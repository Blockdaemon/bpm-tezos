package main

import (
	"fmt"

	"github.com/Blockdaemon/bpm-sdk/pkg/node"
	"github.com/Blockdaemon/bpm-sdk/pkg/plugin"
)

// TezosParameterValidator validates Tezos parameters
type TezosParameterValidator struct {
	plugin.SimpleParameterValidator
}

// Validate uses SimpleParameterValidator but also check is the network is correct
func (t TezosParameterValidator) Validate(currentNode node.Node) error {
	network := currentNode.StrParameters["network"]
	if network != "mainnet" && network != "carthagenet" {
		return fmt.Errorf("unknown network: %q", network)
	}

	return nil
}

// NewTezosParameterValidator creates a new instance of TezosParameterValidator
func NewTezosParameterValidator(pluginParameters []plugin.Parameter) TezosParameterValidator {
	return TezosParameterValidator{
		SimpleParameterValidator: plugin.NewSimpleParameterValidator(pluginParameters),
	}
}
