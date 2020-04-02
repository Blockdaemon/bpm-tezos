package main

import (
	"fmt"
	"path/filepath"

	"github.com/Blockdaemon/bpm-sdk/pkg/docker"
	"github.com/Blockdaemon/bpm-sdk/pkg/node"
	"github.com/Blockdaemon/bpm-sdk/pkg/plugin"
)

var version string

const (
	tezosContainerName  = "tezos"
	tezosDataVolumeName = "tezos-data"
	tezosCmdFile        = "configs/tezos.dockercmd"

	collectorContainerName = "collector"
	collectorImage         = "docker.io/blockdaemon/tezos-collector:0.5.0"
	collectorEnvFile       = "configs/collector.env"

	testContainerName = "tezos-test"
	testImage         = "docker.io/blockdaemon/tezos-tests:1.0.0"
)

func getContainerImage(network string) string {
	// TODO: Not all containers are versioned yet but they will be as soon as there are new commits
	// https://gitlab.com/tezos/tezos/issues/682
	if network == "mainnet" {
		return "docker.io/tezos/tezos-bare:mainnet"
	} else if network == "carthagenet" {
		return "docker.io/tezos/tezos-bare:carthagenet"
	} else {
		panic(fmt.Sprintf("Unknown network: %q", network))
	}
}

func getContainers(currentNode node.Node) []docker.Container {
	network := currentNode.StrParameters["network"]
	image := getContainerImage(network)

	tezosContainer := docker.Container{
		Name:    tezosContainerName,
		Image:   image,
		CmdFile: tezosCmdFile,
		User:    "root",
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

	collectorContainer := docker.Container{
		Name:        collectorContainerName,
		Image:       collectorImage,
		EnvFilename: collectorEnvFile,
		Mounts: []docker.Mount{
			{
				Type: "bind",
				From: "logs",
				To:   "/data/nodestate",
			},
		},
		CollectLogs: true,
	}

	containers := []docker.Container{
		tezosContainer,
		collectorContainer,
	}

	return containers
}

func main() {
	templates := map[string]string{
		tezosCmdFile:     tezosCmdTpl,
		collectorEnvFile: collectorEnvTpl,
	}

	parameters := []plugin.Parameter{
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
	}

	description := "A tezos package"

	containers := []docker.Container{} // Passing in empty containers because we'll replace the container handlers later anyway

	tezosPlugin := plugin.NewDockerPlugin("tezos", version, description, parameters, templates, containers)
	// Replace handlers with tezos specific ones
	tezosPlugin.ParameterValidator = NewTezosParameterValidator(tezosPlugin.Meta().Parameters)
	tezosPlugin.IdentityCreator = NewTezosIdentityCreator()
	tezosPlugin.LifecycleHandler = NewTezosLifecycleHandler()
	tezosPlugin.Upgrader = NewTezosUpgrader()
	tezosPlugin.Tester = NewTezosTester()

	plugin.Initialize(tezosPlugin)
}
