all: up

install:
	vagrant plugin install dotenv
	vagrant plugin install vagrant-proxyconf
	vagrant plugin install vagrant-libvirt

provision: install
	vagrant up --provision --provider=libvirt --no-parallel --no-destroy-on-error

up:
	vagrant up

destroy:
	vagrant destroy --force --graceful --no-tty

validate: up
	vagrant validate

package: up
	mkdir -p ../../output
	vagrant package --output ../../output/provisioned-ubuntu2204.box

cleanup:
	-vagrant destroy --force --graceful --no-tty
	-vagrant box remove -f ubuntu-22.04/runner
	rm -f ../../output/provisioned-ubuntu2204.box
	-virsh destroy ubuntu-22.04_provision
	-virsh undefine ubuntu-22.04_provision --snapshots-metadata --managed-save
