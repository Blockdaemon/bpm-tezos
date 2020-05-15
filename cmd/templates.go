package main

const (
	// Originally generated using the `tezos-node config init` command
	tezosConfigTpl = `{
	"data-dir": "/data",
	"rpc": {
		"listen-addrs": [
			"0.0.0.0:8732"
		],
		"cors-origin": [
			"*"
		]
	},
	"p2p": {
		"listen-addr": "0.0.0.0:9732"
	},
	"shell": {
		"history_mode": "full"
	},
	"network": "{{ .Node.StrParameters.network }}"
	}	
`

	collectorEnvTpl = `SERVICE_PORT=8732
SERVICE_HOST={{ .Node.NamePrefix }}` + tezosContainerName
)
