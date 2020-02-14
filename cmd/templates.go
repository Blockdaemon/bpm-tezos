package main

const (
	tezosCmdTpl = `tezos-node
run
--bootstrap-threshold=10
--net-addr=0.0.0.0:9732
--history-mode=full
--cors-origin=*
--rpc-addr=0.0.0.0:8732
--data-dir
/data
`
)
