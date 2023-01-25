/*
Copyright 2019 Intel Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package agent

import (
	"fmt"
	"sync"

	"github.com/intel/nri-resmgr/pkg/log"
	policyapi "github.com/intel/nri-resmgr/pkg/policy"
	k8sclient "k8s.io/client-go/kubernetes"

	resmgrcs "github.com/intel/nri-resmgr/pkg/apis/resmgr/generated/clientset/versioned/typed/resmgr/v1alpha1"
	resmgr "github.com/intel/nri-resmgr/pkg/apis/resmgr/v1alpha1"
	nrtapi "github.com/k8stopologyawareschedwg/noderesourcetopology-api/pkg/generated/clientset/versioned/typed/topology/v1alpha1"
)

// Get cri-resmgr config
type getConfigFn func() resmgrConfig

// resmgrConfig represents cri-resmgr configuration
type resmgrConfig map[string]string

// resmgrAdjustment represents external adjustments for the resource-manager
type resmgrAdjustment map[string]*resmgr.Adjustment

// resmgrStatus represents the status of an external adjustment update
type resmgrStatus struct {
	request error
	errors  map[string]string
}

// ResourceManagerAgent is the interface exposed for the CRI Resource Manager Congig Agent
type ResourceManagerAgent interface {
	Run() error
	UpdateNrtCR(policy string, zones []*policyapi.TopologyZone) error
}

// agent implements ResourceManagerAgent
type agent struct {
	log.Logger                      // Our logging interface
	cli        *k8sclient.Clientset // K8s client
	extCli     *resmgrcs.CriresmgrV1alpha1Client
	nrtCli     *nrtapi.TopologyV1alpha1Client
	server     agentServer   // gRPC server listening for requests from cri-resource-manager
	watcher    k8sWatcher    // Watcher monitoring events in K8s cluster
	updater    configUpdater // Client sending config updates to cri-resource-manager
	nrtLock    sync.Mutex    // serialize async CR updates
}

// NewResourceManagerAgent creates a new instance of ResourceManagerAgent
func NewResourceManagerAgent() (ResourceManagerAgent, error) {
	var err error

	a := &agent{
		Logger: log.NewLogger("resource-manager-agent"),
	}

	if a.cli, a.extCli, a.nrtCli, err = a.getK8sClient(opts.kubeconfig); err != nil {
		return nil, agentError("failed to get k8s client: %v", err)
	}

	if a.watcher, err = newK8sWatcher(a.cli, a.extCli); err != nil {
		return nil, agentError("failed to initialize watcher instance: %v", err)
	}

	if a.updater, err = newConfigUpdater(); err != nil {
		return nil, agentError("failed to initialize config updater instance: %v", err)
	}

	return a, nil
}

// Start starts the resource manager.
func (a *agent) Run() error {
	a.Info("starting CRI Resource Manager Agent")

	if err := a.watcher.Start(); err != nil {
		return agentError("failed to start watcher: %v", err)
	}

	if err := a.updater.Start(); err != nil {
		return agentError("failed to start config updater: %v", err)
	}

	for {
		select {
		case config, ok := <-a.watcher.ConfigChan():
			if ok {
				a.updater.UpdateConfig(&config)
			}
		case adjust, ok := <-a.watcher.AdjustmentChan():
			if ok {
				a.updater.UpdateAdjustment(&adjust)
			}
		case status, ok := <-a.updater.StatusChan():
			if ok {
				a.Info("got status %v", status)
				if err := a.watcher.UpdateStatus(status); err != nil {
					a.Error("failed to update adjustment node status: %v", err)
				}
			}
		}
	}
}

func agentError(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}
