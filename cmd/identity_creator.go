package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"go.blockdaemon.com/bpm/sdk/pkg/docker"
	"go.blockdaemon.com/bpm/sdk/pkg/node"
)

const (
	tezosInitContainerName  = "tezos-init"
	identityDirectory       = "identity"
	identityConfigDirectory = "identity-config"
)

// TezosIdentityCreator provides functions to create and remove the identity (e.g. private keys) of a node
type TezosIdentityCreator struct{}

// CreateIdentity creates the identity of a node
func (t TezosIdentityCreator) CreateIdentity(currentNode node.Node) error {
	// Create identity directory if it doesn't exist yet
	identityPath := filepath.Join(currentNode.NodeDirectory(), identityDirectory)
	if err := os.MkdirAll(identityPath, os.FileMode(0750)); err != nil {
		return err
	}

	// Check if there already is an identiy file
	_, err := os.Stat(filepath.Join(identityPath, "identity.json"))
	if err == nil {
		fmt.Println("Identity file already exists, skipping creation")
		return nil
	}

	// Write a minimal configuration file to configure the network for the identity creation
	// We cannot use the real configuration file because it gets created after the identity
	// { "data-dir": "//home/tezos/.tezos-node/", "p2p": {}, "network": "carthagenet" }
	identityConfigPath := filepath.Join(currentNode.NodeDirectory(), identityConfigDirectory)
	if err := os.MkdirAll(identityConfigPath, os.FileMode(0750)); err != nil {
		return err
	}
	minimalConfig := filepath.Join(identityConfigPath, "minimal-config.json")

	content := []byte(fmt.Sprintf(`{"p2p": {}, "network": "%s"}`, currentNode.StrParameters["network"]))
	fmt.Println(minimalConfig)
	if err := ioutil.WriteFile(minimalConfig, content, os.FileMode(0600)); err != nil {
		return err
	}
	time.Sleep(time.Duration(300) * time.Second)

	// Run container to create an identity
	tezosInitContainer := docker.Container{
		Name:  tezosInitContainerName,
		Image: tezosImage,
		User:  "root",
		Cmd: []string{
			"tezos-node",
			"identity",
			"generate",
			"--config",
			"/config/minimal-config.json",
			"--data-dir",
			"/identity",
		},
		Mounts: []docker.Mount{
			{
				Type: "bind",
				From: identityPath,
				To:   "/identity",
			},
			{
				Type: "bind",
				From: identityConfigPath,
				To:   "/config",
			},
		},
	}

	// client, err := docker.InitializeClient(currentNode)
	client, err := docker.InitializeClient(currentNode)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	output, err := client.RunTransientContainer(ctx, tezosInitContainer)
	fmt.Println(output)

	return err
}

// RemoveIdentity removes identity related to the node
func (t TezosIdentityCreator) RemoveIdentity(currentNode node.Node) error {
	identityPath := filepath.Join(currentNode.NodeDirectory(), identityDirectory)
	fmt.Printf("Removing directory %q\n", identityPath)
	return os.RemoveAll(identityPath)
}

// NewTezosIdentityCreator creates a new instance of TezosIdentityCreator
func NewTezosIdentityCreator() TezosIdentityCreator {
	return TezosIdentityCreator{}
}
