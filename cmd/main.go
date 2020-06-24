package main

import (
	"path/filepath"

	"go.blockdaemon.com/bpm/sdk/pkg/docker"
	"go.blockdaemon.com/bpm/sdk/pkg/plugin"
)

var version string

const (
	tezosContainerName = "tezos"
	tezosImage         = "docker.io/tezos/tezos-bare:v7.2"
	tezosConfigFile    = "configs/config.json"

	collectorContainerName = "collector"
	collectorImage         = "docker.io/blockdaemon/tezos-collector:0.5.0"
	collectorEnvFile       = "configs/collector.env"

	testContainerName = "tezos-test"
	testImage         = "docker.io/blockdaemon/tezos-tests:1.0.0"
)

func getContainers() []docker.Container {
	tezosContainer := docker.Container{
		Name:  tezosContainerName,
		Image: tezosImage,
		Cmd: []string{
			"tezos-node",
			"run",
			"--config",
			"/config/config.json",
		},
		User: "root",
		Mounts: []docker.Mount{
			{
				Type: "bind",
				From: "configs/config.json",
				To:   "/config/config.json",
			},
			{
				Type: "bind",
				From: "{{ index .Node.StrParameters \"data-dir\" }}",
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
		tezosConfigFile:  tezosConfigTpl,
		collectorEnvFile: collectorEnvTpl,
	}

	parameters := []plugin.Parameter{
		{
			Name:        "subtype",
			Type:        plugin.ParameterTypeString,
			Description: "The type of node. Only `full` supported currently",
			Mandatory:   false,
			Default:     "full",
		},
		{
			Name:        "network",
			Type:        plugin.ParameterTypeString,
			Description: "The network. Can be either `mainnet` or `carthagenet`",
			Mandatory:   false,
			Default:     "mainnet",
		},
	}

	description := "Tezos full node package"
	containers := getContainers()

	tezosPlugin := plugin.NewDockerPlugin("tezos", version, description, parameters, templates, containers)
	tezosPlugin.ParameterValidator = NewTezosParameterValidator(tezosPlugin.Meta().Parameters)
	tezosPlugin.IdentityCreator = NewTezosIdentityCreator()
	tezosPlugin.Tester = NewTezosTester()

	plugin.Initialize(tezosPlugin)
}
