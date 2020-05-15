package main

import (
	"fmt"

	"go.blockdaemon.com/bpm/sdk/pkg/node"
	"go.blockdaemon.com/bpm/sdk/pkg/plugin"
)

// TezosParameterValidator validates Tezos parameters
type TezosParameterValidator struct {
	plugin.SimpleParameterValidator
}

// ValidateParameters uses SimpleParameterValidator but also check is the network is correct
func (t TezosParameterValidator) ValidateParameters(currentNode node.Node) error {
	if err := t.SimpleParameterValidator.ValidateParameters(currentNode); err != nil {
		return err
	}

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
