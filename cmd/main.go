package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Blockdaemon/bpm-sdk/pkg/docker"
	"github.com/Blockdaemon/bpm-sdk/pkg/node"
	"github.com/Blockdaemon/bpm-sdk/pkg/plugin"
)

var version string

const (
	tezosContainerName     = "tezos"
	tezosInitContainerName = "tezos-init"
	tezosDataVolumeName    = "tezos-data"
	tezosCmdFile           = "tezos.dockercmd"

	networkName     = "tezos"
	initNetworkName = "tezos-init"
)

// TezosLifecycleHandler makes use of DockerLifecycleHandler to manage containers but calls it with different docker images.
//
// This is necessary because Tezos uses different images per network. It's a bit of a hack so hopefully this can go away soon
// once Tezos provides one universal image/binary for all networks
type TezosLifecycleHandler struct{}

func (t TezosLifecycleHandler) Start(currentNode node.Node) error {
	handler := plugin.NewDockerLifecycleHandler(getContainers(currentNode))
	return handler.Start(currentNode)
}

func (t TezosLifecycleHandler) Stop(currentNode node.Node) error {
	handler := plugin.NewDockerLifecycleHandler(getContainers(currentNode))
	return handler.Stop(currentNode)
}

func (t TezosLifecycleHandler) Status(currentNode node.Node) (string, error) {
	handler := plugin.NewDockerLifecycleHandler(getContainers(currentNode))
	return handler.Status(currentNode)
}

func (t TezosLifecycleHandler) RemoveData(currentNode node.Node) error {
	handler := plugin.NewDockerLifecycleHandler(getContainers(currentNode))
	return handler.RemoveData(currentNode)
}

func (t TezosLifecycleHandler) RemoveRuntime(currentNode node.Node) error {
	handler := plugin.NewDockerLifecycleHandler(getContainers(currentNode))
	return handler.RemoveRuntime(currentNode)
}

func NewTezosLifecycleHandler() TezosLifecycleHandler {
	return TezosLifecycleHandler{}
}

// TezosUpgrader makes use of DockerUpgradeer but calls it with different docker images.
//
// See TezosLifecycleHandler for the reasoning
type TezosUpgrader struct{}

func (t TezosUpgrader) Upgrade(currentNode node.Node) error {
	upgrader := plugin.NewDockerUpgrader(getContainers(currentNode))
	return upgrader.Upgrade(currentNode)
}

func NewTezosUpgrader() TezosUpgrader {
	return TezosUpgrader{}
}

// TezosConfigurator makes use of FileConfigurator for rendering config files and adds functionality to create the tezos node identity
type TezosConfigurator struct {
	plugin.FileConfigurator
}

func (t TezosConfigurator) Configure(currentNode node.Node) error {
	if err := t.FileConfigurator.ValidateParameters(currentNode); err != nil {
		return err
	}
	network := currentNode.StrParameters["network"]
	if network != "mainnet" && network != "carthagenet" && network != "babylonnet" && network != "zeronet" {
		return fmt.Errorf("unknown network: %q", network)

	}

	if err := t.FileConfigurator.Configure(currentNode); err != nil {
		return err
	}

	identityDir := filepath.Join(currentNode.ConfigsDirectory(), "identity")

	_, err := os.Stat(filepath.Join(identityDir, "identity.json"))
	if err == nil {
		fmt.Println("Identity file already exists, skipping creation")
		return nil
	}

	// Tezos expects and empty directory to create an identity in
	if err := os.MkdirAll(identityDir, os.ModePerm); err != nil {
		return err
	}

	tezosInitContainer := docker.Container{
		Name:      tezosInitContainerName,
		Image:     getContainerImage(network),
		NetworkID: initNetworkName,
		Cmd: []string{
			"tezos-node",
			"identity",
			"generate",
			"--data-dir=/data",
		},
		Mounts: []docker.Mount{
			{
				Type: "bind",
				From: identityDir,
				To:   "/data",
			},
		},
	}

	client, err := docker.NewBasicManager(currentNode.NamePrefix(), currentNode.ConfigsDirectory())
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	if err := client.NetworkExists(ctx, tezosInitContainer.NetworkID); err != nil {
		return err
	}

	output, err := client.RunTransientContainer(ctx, tezosInitContainer)

	// err = os.Chown("test.txt", os.Getuid(), os.Getgid())
	fmt.Println(output)

	return err
}

func (t TezosConfigurator) RemoveConfig(currentNode node.Node) error {
	if err := t.FileConfigurator.RemoveConfig(currentNode); err != nil {
		return err
	}

	identityDir := filepath.Join(currentNode.ConfigsDirectory(), "identity")
	fmt.Printf("Removing directory %q\n", identityDir)
	return os.RemoveAll(identityDir)
}

func NewTezosConfigurator(configFilesAndTemplates map[string]string, pluginParameters []plugin.Parameter) TezosConfigurator {
	fileConfigurator := plugin.NewFileConfigurator(map[string]string{
		tezosCmdFile: tezosCmdTpl,
	}, pluginParameters)

	return TezosConfigurator{
		FileConfigurator: fileConfigurator,
	}
}

func getContainerImage(network string) string {
	// TODO: Not all containers are versioned yet but they will be as soon as there are new commits
	// https://gitlab.com/tezos/tezos/issues/682
	if network == "mainnet" {
		return "docker.io/tezos/tezos-bare:master_6d2aa96e_20200212132052"
	} else if network == "carthagenet" {
		return "docker.io/tezos/tezos-bare:carthagenet"
	} else if network == "babylonnet" {
		return "docker.io/tezos/tezos-bare:babylonnet"
	} else if network == "zeronet" {
		return "docker.io/tezos/tezos-bare:zeronet"
	} else {
		panic("Unknown network")
	}
}

func getContainers(currentNode node.Node) []docker.Container {
	network := currentNode.StrParameters["network"]
	image := getContainerImage(network)

	tezosContainer := docker.Container{
		Name:      tezosContainerName,
		Image:     image,
		CmdFile:   tezosCmdFile,
		NetworkID: networkName,
		User:      "root",
		Mounts: []docker.Mount{
			{
				Type: "volume",
				From: tezosDataVolumeName,
				To:   "/data",
			},
			{
				Type: "bind",
				From: filepath.Join("identity", "identity.json"),
				To:   "/data/identity.json", // Tezos needs the identity file in /data
			},
		},
		Ports: []docker.Port{
			{
				HostIP:        "0.0.0.0",
				HostPort:      "9732",
				ContainerPort: "9732",
				Protocol:      "tcp",
			},
			{
				HostIP:        "127.0.0.1",
				HostPort:      "8732",
				ContainerPort: "8732",
				Protocol:      "tcp",
			},
		},
		CollectLogs: true,
	}

	// TODO: Add collector

	containers := []docker.Container{
		tezosContainer,
	}

	return containers
}

func main() {
	tezosTemplates := map[string]string{
		tezosCmdFile: tezosCmdTpl,
	}

	meta := plugin.MetaInfo{
		Version:         version,
		Description:     "A tezos package",
		ProtocolVersion: "1.1.0",
		Parameters: []plugin.Parameter{
			{
				Name:        "subtype",
				Type:        plugin.ParameterTypeString,
				Description: "The type of node. Only `watcher` supported currently",
				Mandatory:   false,
				Default:     "watcher",
			},
			{
				Name:        "network",
				Type:        plugin.ParameterTypeString,
				Description: "The network. Can be one of `mainnet`, `carthagenet`, `babylonnet`, `zeronet`",
				Mandatory:   false,
				Default:     "mainnet",
			},
		},
		Supported: []string{
			plugin.SupportsTest,
		},
	}

	lifeCycleHandler := NewTezosLifecycleHandler()

	configurator := NewTezosConfigurator(tezosTemplates, meta.Parameters)

	upgrader := NewTezosUpgrader()

	tester := plugin.NewDummyTester()

	tezosPlugin := plugin.NewPlugin(
		"tezos",
		meta,
		configurator,
		lifeCycleHandler,
		upgrader,
		tester,
	)

	plugin.Initialize(tezosPlugin)
}
