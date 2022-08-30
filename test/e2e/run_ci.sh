#!/bin/bash

# This is executed by the Github Workflow Actions and it will run e2e tests
# with suitable parameters.

set -x

if [ -z "$1" ]; then
    echo "Usage: $0 <test-directory-to-use>"
    exit 1
fi

cd $GITHUB_WORKSPACE/test/e2e

# Set any site specific environment variables to the env file so that they
# will be available to the run_tests.sh script.
if [ -f $GITHUB_WORKSPACE/../env ]; then
    . $GITHUB_WORKSPACE/../env
fi

# Set the govm cpu to be like this VM cpu.
VM_QEMU_EXTRA="-cpu host,topoext=on" ./run_tests.sh $1
