#!/usr/bin/env bash

# fail on unset variables and command errors
#set -eu -o pipefail # -x: is for debugging

export GOPATH="/home/vagrant/go"
mkdir -p $GOPATH/bin
mkdir -p $GOPATH/src

# golangci-lint
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "$(go env GOPATH)/bin" v1.46.2

grep GOPATH ~/.bashrc > /dev/null
if [ $? -ne 0 ]; then
  echo "export GOPATH=$GOPATH" >> ~/.bashrc
  echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.bashrc
  mkdir -p $GOPATH/bin
fi

error() {
    (echo ""; echo "error: $1" ) >&2
    exit 1
}

# Vagrant needs plugins to be installed
vagrant plugin install dotenv
vagrant plugin install vagrant-proxyconf
vagrant plugin install vagrant-qemu

# Make sure the base images needed by e2e tests are pre-loaded
vagrant box add --provider libvirt generic/ubuntu2204
vagrant box add --provider libvirt generic/fedora37

HOST_VM_DIR="/home/vagrant/files"
mkdir -p "$HOST_VM_DIR" ||
        error "cannot create directory for VM images: $HOST_VM_DIR"

fetch-vm-file() {
    local url="$1"
    local file=$(basename $url)
    local image decompress
    case $file in
        *.xz)
            image=${file%.xz}
            decompress="xz -d"
            ;;
        *.bz2)
            image=${file%.bz2}
            decompress="bzip -d"
            ;;
        *.gz)
            image=${file%.gz}
            decompress="gzip -d"
            ;;
        *)
            image="$file"
            decompress=""
            ;;
    esac
    [ -s "$HOST_VM_DIR/$image" ] || {
        echo "VM image $HOST_VM_DIR/$image not found..."
        [ -s "$HOST_VM_DIR/$file" ] || {
            echo "downloading VM file $image..."
            wget --progress=dot:giga -O "$HOST_VM_DIR/$file" "$url" ||
            echo "failed to download VM file ($url), skipping it"
        }
        if [ -s "$HOST_VM_DIR/$file" ]; then
            if [ -n "$decompress" ]; then
		echo "decompressing VM file $file..."
		( cd "$HOST_VM_DIR" && $decompress $file ) ||
                    error "failed to decompress $file to $image using $decompress"
            fi
            if [ ! -s "$HOST_VM_DIR/$image" ]; then
		error "internal error, fetching+decompressing $url did not produce $HOST_VM_DIR/$image"
            fi
	fi
    }
}

# Pre-install any large files that are needed when the e2e test scripts.
#fetch-vm-file "<add-here>"
