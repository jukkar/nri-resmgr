# Enabling CRI-RM

## 1. Install the package

For Fedora, CentOS, SUSE:
 - CRIRM_VERSION=`curl -s "https://api.github.com/repos/intel/cri-resource-manager/releases/latest" | \
               jq .tag_name | tr -d '"v'`
 - source /etc/os-release
 - [ "$ID" = "sles" -o "$ID" = "opensuse-leap" ] && export ID=suse
 - sudo rpm -Uvh https://github.com/intel/cri-resource-manager/releases/download/v${CRIRM_VERSION}/cri-resource-manager-${CRIRM_VERSION}-0.${ID}-${VERSION_ID}.x86_64.rpm

For Ubuntu, Debian:
 - CRIRM_VERSION=`curl -s "https://api.github.com/repos/intel/cri-resource-manager/releases/latest" | \
               jq .tag_name | tr -d '"v'`
 - source /etc/os-release
 - pkg=cri-resource-manager_${CRIRM_VERSION}_${ID}-${VERSION_ID}_amd64.deb; curl -LO https://github.com/intel/cri-resource-manager/releases/download/v${CRIRM_VERSION}/${pkg}; sudo dpkg -i ${pkg}; rm ${pkg}

## 2. Setup & verify:
 - sudo cp /etc/cri-resmgr/fallback.cfg.sample /etc/cri-resmgr/fallback.cfg
 - sudo systemctl enable cri-resource-manager && sudo systemctl start cri-resource-manager
 - systemctl status cri-resource-manager

## 3. Run: "systemctl cat kubelet". Look for: EnvironmentFile=/.../kubeadm-flags.env

## 4. Run: "sudo vim /.../kubeadm-flags.env"

## 5. Change the --container-runtime-endpoint from containerd socket to cri-resmgr.scok:
 - KUBELET_KUBEADM_ARGS="--container-runtime-endpoint=unix:///var/run/cri-resmgr/cri-resmgr.sock --pod-infra-container-image=registry.k8s.io/pause:3.9"

 Alternatively run:
 - kubelet <other-kubelet-options> --container-runtime=remote \
     --container-runtime-endpoint=unix:///var/run/cri-resmgr/cri-resmgr.sock

## 6. Restart kubelet: "systemctl restart kubelet"

## 7. Run CRI-RM with your desired policy with the following command:
 - cri-resmgr --force-config <config-file> --runtime-socket unix:///var/run/containerd/containerd.sock

## 8. Deploy your pod

## 9. See what resources the container got assigned with:
 - kubectl exec -c $container $pod -- cat /proc/self/status | grep _allowed_list

Output should be something along the lines of:

Cpus_allowed_list:	10-15 </br>
Mems_allowed_list:	1

# Enabling NRI

## 1. With Containerd as the runtime:

Replace containerd in the system with 1.7 or newer version (NRI server not supported in older versions).

Change the Kubelet --container-runtime-endpoint to containerd socket:
 - sudo vim /.../kubeadm-flags.env
 - KUBELET_KUBEADM_ARGS="--container-runtime-endpoint=unix:///var/run/containerd/containerd.sock --pod-infra-container-image=registry.k8s.io/pause:3.9"
 - sudo vim /etc/sysconfig/kubelet 
 - KUBELET_EXTRA_ARGS= --container-runtime-endpoint=/var/run/containerd/containerd.sock <- Remember this aswell
 - systemctl restart containerd
 - systemctl restart kubelet

Edit the containerd config file and look for the section "io.containerd.nri.v1.nri" and replace "disable = true" with "disable = false":
 - "vim /etc/containerd/config.toml"

## 1. With CRI-O as the runtime:

Replace crio in the system with v1.26.2 or newer version.

Change the Kubelet --container-runtime-endpoint to crio socket:
 - sudo vim /.../kubeadm-flags.env
 - KUBELET_KUBEADM_ARGS="--container-runtime-endpoint=unix:///var/run/crio/crio.sock --pod-infra-container-image=registry.k8s.io/pause:3.9"
 - sudo vim /etc/sysconfig/kubelet
 - KUBELET_EXTRA_ARGS= --container-runtime-endpoint=/var/run/crio/crio.sock <- Remember this aswell
 - systemctl restart crio
 - systemctl restart kubelet

Edit the CRI-O config file and look for "[crio.nri]" section, switch "enable_nri = false" with "enable_nri = true" and uncomment the NRI config lines.
 - "sudo vim /etc/crio/crio.conf"

Create NRI configuration (not needed with Containerd v1.7:
 - sudo sh -c "mkdir -p /etc/nri; touch /etc/nri/nri.conf; systemctl restart crio"

## 2. Build the policies:
 - git clone https://github.com/jukkar/nri-resmgr.git
 - cd nri-resmgr
 - make
 - make images

## 3. Apply required crds:
 - kubectl apply -f deployment/base/crds/noderesourcetopology_crd.yaml

## 4. Start the NRI plugin you want to run:
 For Containerd:
 - ctr -n k8s.io images import build/images/nri-resmgr-topology-aware-image-9797e8de7107.tar
 For CRI-O????????????????

 Deploy the plugin:
 - kubectl apply -f build/images/nri-resmgr-topology-aware-deployment.yaml

## 5. Deploy your pod.

## 6. See the resources the pod got assigned with:
 - kubectl exec $pod -c $container  -- grep allowed_list: /proc/self/status
 - Output should look similar to the output of CRI-RM

## Common error:
NRI socket errors, ex: "failed to start resource manager: failed to start NRI plugin: failed to connect to NRI service: dial unix /var/run/nri.sock: connect: connection refused"
 - sudo vim /etc/containerd/config.toml (or the crio.conf if using crio)
 - Go to the "io.containerd.nri.v1.nri" and look for "socket_path"
 - Make sure this matches the socket path of your NRI socket.
 - systemctl restart containerd
