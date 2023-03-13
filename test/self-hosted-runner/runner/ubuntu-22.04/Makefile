all: up

install:
	vagrant plugin install dotenv
	vagrant plugin install vagrant-proxyconf
	vagrant plugin install vagrant-libvirt

up: prepare
	vagrant up

destroy:
	vagrant destroy --force --graceful --no-tty

validate: up
	vagrant validate

BASE_PACKAGE=`pwd`/../../output/provisioned-ubuntu2204.box

prepare:
	./create-vagrantfile.sh $(BASE_PACKAGE)

cleanup:
	-vagrant destroy --force --graceful --no-tty
	-virsh destroy ubuntu-22.04_runner
	-virsh vol-delete --pool default runner_vagrant_box_image_0_box.img
	rm -r .vagrant
	rm -f Vagrantfile
