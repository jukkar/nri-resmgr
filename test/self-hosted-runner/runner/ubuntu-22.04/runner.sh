#!/usr/bin/env bash

# This will create the VM from vagrant box file, start the
# runner and wait for connections. After the run, the VM is
# destroyed and then re-created.

STOPPED=0
trap ctrl_c INT TERM

ctrl_c() {
    STOPPED=1
}

debug() {
    if [ -z "$DEBUG" ]; then
	return
    fi

    printf "$*\n" > /dev/tty
}

wait_for_ssh() {
    while [ "$(vagrant status --machine-readable 2>/dev/null | awk -F, '/,state,/ { print $4 }')" != "running" ]; do
        echo "Waiting for VM SSH server to respond"
        sleep 1

	if [ $STOPPED -eq 1 ]; then
	    echo "Exiting."
	    exit 1
	fi
    done
}

build_jwt_payload() {
        jq -c --arg iat_str "$(date +%s)" --arg app_id "${app_id}" \
        '
        ($iat_str | tonumber) as $iat
        | .iat = $iat
        | .exp = ($iat + 300)
        | .iss = ($app_id | tonumber)
        ' <<< "${payload_template}" | tr -d '\n'
}

b64enc() {
    openssl enc -base64 -A | tr '+/' '-_' | tr -d '='
}

json() {
    jq -c . | LC_CTYPE=C tr -d '\n'
}

rs256_sign() {
    openssl dgst -binary -sha256 -sign <(printf '%s\n' "$1")
}

export BASE_PACKAGE=`pwd`/../../output/provisioned-ubuntu2204.box

if [ ! -f $BASE_PACKAGE ]; then
    echo "Create a box file $BASE_PACKAGE that contains the provisioned environment."
    echo "Follow the instructions in ../../README.md file and provision the base image first."
    exit 1
fi

# All the configuration data is stored in env file.
. ../../env

if [ ! -f "$CONTAINERD_SRC/bin/ctr" ]; then
    echo "########"
    echo "WARNING: containerd sources / binaries not found in $CONTAINERD_SRC"
    echo "########"
fi

get_runner_token() {
    local payload sig access_token

    registration_url="https://api.github.com/repos/${GH_REPO}/actions/runners/registration-token"
    app_id=${GH_APP_ID}
    app_private_key="$(< ${GH_APP_PRIVATE_KEY_FILE})"
    payload_template='{}'
    header='{"alg": "RS256","typ": "JWT"}'

    payload=$(build_jwt_payload) || return
    signed_content="$(json <<<"$header" | b64enc).$(json <<<"$payload" | b64enc)"
    sig=$(printf %s "$signed_content" | rs256_sign "$app_private_key" | b64enc)

    generated_jwt="${signed_content}.${sig}"
    app_installations_url="https://api.github.com/app/installations"
    app_installations_response=$(curl -sX GET -H "Authorization: Bearer  ${generated_jwt}" -H "Accept: application/vnd.github.v3+json" ${app_installations_url})
    access_token_url=$(echo $app_installations_response | jq '.[] | select (.app_id  == '${app_id}') .access_tokens_url' --raw-output)
    access_token_response=$(curl -sX POST -H "Authorization: Bearer  ${generated_jwt}" -H "Accept: application/vnd.github.v3+json" ${access_token_url})
    access_token=$(echo $access_token_response | jq .token --raw-output)
    payload=$(curl -sX POST -H "Authorization: Bearer  ${access_token}" -H "Accept: application/vnd.github.v3+json" ${registration_url})

    debug "GH application id : '${GH_APP_ID}'"
    debug "GH private key    : '${GH_APP_PRIVATE_KEY_FILE}'"
    debug "GH registration URL      : '${registration_url}'"
    debug "GH app installations URL : '${app_installations_url}'"
    debug "GH app installations response : '${app_installations_response}'"
    debug "GH access token URL : '${access_token_url}'"
    debug "GH access token response : '${access_token_response}'"
    debug "GH access token : '${access_token}'"
    debug "GH reponse payload : $payload"

    echo $(echo $payload | jq .token --raw-output)
}

get_latest_runner() {
    LATEST_RUNNER=.github_self_hosted_runner_version

    latest_runner_version=$(curl -I -v -s https://github.com/actions/runner/releases/latest 2>&1 | sed -n 's/^< location: \(.*\)$/\1/p' | awk -F/ '{ print $NF }' | sed 's/v//' | tr -d '\r')

    if [ -e $LATEST_RUNNER ]; then
	if [ -f actions-runner/config.sh ]; then
	    github_self_hosted_runner_version=$(cat $LATEST_RUNNER)
	    if [ "$github_self_hosted_runner_version" == "$latest_runner_version" ]; then
		return
	    fi

	    if [ ! -z "$github_self_hosted_runner_version" ]; then
		echo "Current runner version ($github_self_hosted_runner_version) is not up to date."
	    fi

	    if [ -d actions-runner.old ]; then
		rm -rf actions-runner.old
	    fi

	    mv actions-runner actions-runner.old >/dev/null 2>&1
	fi
    fi

    echo "Installing latest runner version ($latest_runner_version)"

    mkdir -p actions-runner

    curl -sL "https://github.com/actions/runner/releases/download/v${latest_runner_version}/actions-runner-linux-x64-${latest_runner_version}.tar.gz" | tar xzC actions-runner

    echo $latest_runner_version > $LATEST_RUNNER
}

while [ $STOPPED -eq 0 ]; do
    GH_RUNNER_TOKEN=$(get_runner_token)
    if [ "$GH_RUNNER_TOKEN" == "null" ]; then
	echo "Cannot get self-hosted runner registration token!"
	break
    fi

    get_latest_runner

    make up

    wait_for_ssh

    # Generate env file for DNS stuff
    vagrant ssh -c "echo export dns_nameserver=$DNS_NAMESERVER > env; echo export dns_search_domain=$DNS_SEARCH_DOMAIN >> env"

    # Make sure the mount in Vagrantfile is pointing to same directory.
    vagrant ssh -c "echo export containerd_src=/home/vagrant/containerd >> env"

    # Configure the runner
    vagrant ssh -c "cd actions-runner; ./config.sh --replace --name '$GH_RUNNER_NAME' --url '$GH_RUNNER_URL' --token '$GH_RUNNER_TOKEN' --unattended --ephemeral" | \
	tee /dev/tty | egrep -q -e "Http response code: NotFound from " -e "Invalid configuration provided for url"
    if [ $? -eq 0 ]; then
	echo "Action runner configuration failed. Fix things and retry."
	# User stopped the script, we should quit
	STOPPED=1
	break
    fi

    # and then run it
    vagrant ssh -c "cd actions-runner; ./run.sh" | tee /dev/tty | egrep -q -e "Exiting\.\.\." -e "An error occurred: Not configured\. Run config"
    if [ $? -eq 0 ]; then
	# User stopped the script, we should quit
	STOPPED=1
    fi

    # Remove the runner after we have finished working with it
    vagrant ssh -c "cd actions-runner; ./config.sh remove --token '${GH_RUNNER_TOKEN}'"

    make destroy

    # Remove the work directory as it is now useless
    if [ -d actions-runner/_work ]; then
	rm -rf actions-runner/_work
    fi

    # Re-read the values if user has changed them.
    . ../../env
done
