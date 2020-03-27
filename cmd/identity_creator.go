package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Blockdaemon/bpm-sdk/pkg/docker"
	"github.com/Blockdaemon/bpm-sdk/pkg/node"
)

const (
	tezosInitContainerName = "tezos-init"
	identityDirectory      = "identity"
)

// TezosIdentityCreator provides functions to create and remove the identity (e.g. private keys) of a node
type TezosIdentityCreator struct{}

// CreateIdentity creates the identity of a node
func (t TezosIdentityCreator) CreateIdentity(currentNode node.Node) error {
	identityPath := filepath.Join(currentNode.NodeDirectory(), identityDirectory)

	// network := currentNode.StrParameters["network"]
	// if network != "mainnet" && network != "carthagenet" && network != "babylonnet" && network != "zeronet" {
	// 	return fmt.Errorf("unknown network: %q", network)

	// }

	// Create identity directory if it doesn't exist yet
	if err := os.MkdirAll(identityPath, os.ModePerm); err != nil {
		return err
	}

	// Check if there already is an identiy file
	_, err := os.Stat(filepath.Join(identityPath, "identity.json"))
	if err == nil {
		fmt.Println("Identity file already exists, skipping creation")
		return nil
	}

	// Run container to create an identity
	tezosInitContainer := docker.Container{
		Name:  tezosInitContainerName,
		Image: getContainerImage(currentNode.StrParameters["network"]),
		User:  "root",
		Cmd: []string{
			"tezos-node",
			"identity",
			"generate",
			"--data-dir=/data",
		},
		Mounts: []docker.Mount{
			{
				Type: "bind",
				From: identityPath,
				To:   "/data",
			},
		},
	}

	client, err := docker.InitializeClient(currentNode)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	if err := client.NetworkExists(ctx, currentNode.StrParameters["docker-network"]); err != nil {
		return err
	}

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
