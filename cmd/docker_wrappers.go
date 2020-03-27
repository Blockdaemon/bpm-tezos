package main

import (
	"github.com/Blockdaemon/bpm-sdk/pkg/node"
	"github.com/Blockdaemon/bpm-sdk/pkg/plugin"
)

// TezosLifecycleHandler makes use of DockerLifecycleHandler to manage containers but calls it with different docker images.
//
// This is necessary because Tezos uses different images per network. It's a bit of a hack so hopefully this can go away soon
// once Tezos provides one universal image/binary for all networks
type TezosLifecycleHandler struct{}

// Start starts a new tezos node
func (t TezosLifecycleHandler) Start(currentNode node.Node) error {
	handler := plugin.NewDockerLifecycleHandler(getContainers(currentNode))
	return handler.Start(currentNode)
}

// Stop stops the containers of a tezos node
func (t TezosLifecycleHandler) Stop(currentNode node.Node) error {
	handler := plugin.NewDockerLifecycleHandler(getContainers(currentNode))
	return handler.Stop(currentNode)
}

// Status returns the status of a tezos node
func (t TezosLifecycleHandler) Status(currentNode node.Node) (string, error) {
	handler := plugin.NewDockerLifecycleHandler(getContainers(currentNode))
	return handler.Status(currentNode)
}

// RemoveData removes the data of a tezos node
func (t TezosLifecycleHandler) RemoveData(currentNode node.Node) error {
	handler := plugin.NewDockerLifecycleHandler(getContainers(currentNode))
	return handler.RemoveData(currentNode)
}

// RemoveRuntime removes the containers of a tezos node
func (t TezosLifecycleHandler) RemoveRuntime(currentNode node.Node) error {
	handler := plugin.NewDockerLifecycleHandler(getContainers(currentNode))
	return handler.RemoveRuntime(currentNode)
}

// NewTezosLifecycleHandler initializes a new instance of TezosLifecycleHandler
func NewTezosLifecycleHandler() TezosLifecycleHandler {
	return TezosLifecycleHandler{}
}

// TezosUpgrader makes use of DockerUpgradeer but calls it with different docker images.
//
// See TezosLifecycleHandler for the reasoning
type TezosUpgrader struct{}

// Upgrade upgrades to a new version of the tezos images
func (t TezosUpgrader) Upgrade(currentNode node.Node) error {
	upgrader := plugin.NewDockerUpgrader(getContainers(currentNode))
	return upgrader.Upgrade(currentNode)
}

// NewTezosUpgrader initializes a new instance of TezosUpgrader
func NewTezosUpgrader() TezosUpgrader {
	return TezosUpgrader{}
}
