package main

import (
	"context"
	"fmt"
	"time"

	"go.blockdaemon.com/bpm/sdk/pkg/docker"
	"go.blockdaemon.com/bpm/sdk/pkg/node"
)

// TezosTester tests a tezos node
type TezosTester struct{}

// Test tests the current node
func (t TezosTester) Test(currentNode node.Node) (bool, error) {
	endpoint := fmt.Sprintf("http://%s%s:8732", currentNode.NamePrefix(), tezosContainerName)

	// Run container to test the node
	testContainer := docker.Container{
		Name:  testContainerName,
		Image: testImage,
		Cmd: []string{
			endpoint,
			"main",
		},
	}

	// client, err := docker.InitializeClient(currentNode)
	client, err := docker.NewBasicManager(currentNode)
	if err != nil {
		return false, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	output, err := client.RunTransientContainer(ctx, testContainer)
	fmt.Println(output)

	success := err == nil
	return success, err

}

// NewTezosTester creates a new instance of TezosTester
func NewTezosTester() TezosTester {
	return TezosTester{}
}
