# Tezos BPM package

## Tutorial: Starting a Tezos carthagenet node

First, let's list all commands available for managing nodes:

```bash
bpm nodes --help
```

You should see the following output:

```
Manage blockchain nodes

Usage:
  bpm nodes [command]

Available Commands:
  configure   Configure a new blockchain node
  remove      Remove blockchain node data and configuration
  show        Print a resource to stdout
  start       Start a blockchain node
  status      Display statuses of configured nodes
  stop        Stops a running blockchain node
  test        Tests a running blockchain node

Flags:
  -h, --help   help for nodes

Global Flags:
      --base-dir string           The directory plugins and configuration are stored (default "~/.bpm/")
      --debug                     Enable debug output
      --package-registry string   The package registry provides packages to install (default "https://dev.registry.blockdaemon.com")
  -y, --yes                       Automatic yes to prompts; assume "yes" as answer to all prompts and run non-interactively

Use "bpm nodes [command] --help" for more information about a command.
```

### Installing the tezos package

Before diving into the actual node management we need to install the tezos package:

```bash
bpm packages install tezos
```

### Configuring the node

Run the configure command with `--help` to see the available parameters:

```
bpm nodes configure tezos --help
```

There are two parameters that are specific to the tezos package:

```
      --network mainnet      The network. Can be one of mainnet, `carthagenet`, `babylonnet`, `zeronet` (default "mainnet")
      --subtype watcher      The type of node. Only watcher supported currently (default "watcher")
```

For now we will leave the defaults which creates a configuration for a non-validating tezos node. Run the configure command to create it now:

```bash
bpm nodes configure tezos --network carthagenet
```

You should see an output similar to this (but with different IDs):

```
No credentials found in "/Users/enode/.bpm/beats", skipping configuration of Blockdaemon monitoring. Please configure your own monitoring in the node configuration files.

Writing file '/Users/enode/.bpm/nodes/bpel6tstp9lnoi401vcg/configs/tezos.dockercmd'
Creating network 'bpm-bpel6tstp9lnoi401vcg-tezos-init'
Creating container 'bpm-bpel6tstp9lnoi401vcg-tezos-init'
Starting container 'bpm-bpel6tstp9lnoi401vcg-tezos-init'
Container 'bpm-bpel6tstp9lnoi401vcg-tezos-init' is not running, skipping stop
Removing container 'bpm-bpel6tstp9lnoi401vcg-tezos-init'
�Generating a new identity... (level: 26.00)
UStored the new identity (idt9NiApmo1gg9cFb23ZCeC3ReGin5) into '/data/identity.json'.


Node with id "bpel6tstp9lnoi401vcg" has been initialized.

To change the configuration, modify the files here:
    /Users/enode/.bpm/nodes/bpel6tstp9lnoi401vcg/configs
To start the node, run:
    bpm nodes start bpel6tstp9lnoi401vcg
To see the status of configured nodes, run:
    bpm nodes status
```

### Node status

You can verify the status of all nodes by running the `status` command:

```bash
bpm nodes status
```

### Starting and stopping the node

To start the node, run the `node start` command. Replace `<node-id>` with the ID outputed by the `configure` or `status` command:

```bash
bpm nodes start <node-id>
```

You should get an output similar to:

```
Creating network 'bpm-bpel6tstp9lnoi401vcg-tezos'
Creating container 'bpm-bpel6tstp9lnoi401vcg-tezos'
Starting container 'bpm-bpel6tstp9lnoi401vcg-tezos'
The node "bpel6tstp9lnoi401vcg" has been started.
```

This shows BPM creating a dedicated Docker network for this node as well as starting the Docker containers.

Verify the node status by running:

```bash
bpm nodes status
```

The node should now show up as `running`, similar to below:

```
        NODE ID        | PACKAGE | STATUS
+----------------------+---------+--------+
  bpel6tstp9lnoi401vcg | running | tezos
```

To stop it temporarily, run:

```bash
bpm nodes stop <node-id>
```

### Removing a node

When removing a node you need to consider the following. A node consists of:

1. Node configuration and secrets (e.g. accounts, passwords, private keys)
2. Runtime (e.g. Docker networks and containers)
3. Data (typically the parts of the Blockchain that have already been synced)

Depending on the use-case it may be desirable to remove all or only parts of the node. For example:

* In order to re-configure a node one might only want to remove the configuration but leave the data intact to avoid having to re-sync the Blockchain
* If the node crashed due to an unexpected error it can make sense to remove the runtime and start it again but keep the configuration and data
* If something went wrong during the initial sync it can help to remove the data and then start the node again to start syncing from scratch

To support the above use-cases plus others we have allowed parameters/flags to be used with our `remove` command.

You can view all the available parameters/flags by running the following BPM command:

```bash
bpm nodes remove --help
```

Which should return:

```
Remove blockchain node data and configuration

Usage:
  bpm nodes remove <id> [flags]

Flags:
      --all       Remove all data, configuration files and node information
      --config    Remove all configuration files but keep data and node information
      --data      Remove all data but keep configuration files and node information
  -h, --help      help for remove
      --runtime   Remove all runtimes but keep configuration files and node information

Global Flags:
      --base-dir string           The directory plugins and configuration are stored (default "~/.bpm/")
      --debug                     Enable debug output
      --package-registry string   The package registry provides packages to install (default "https://dev.registry.blockdaemon.com")
  -y, --yes                       Automatic yes to prompts; assume "yes" as answer to all prompts and run non-interactively
```

For now, let's remove the whole node:

```bash
nodes remove --all <node-id>
```

!!! warning
    By removing the whole node you will also remove the secrets. It's always advisable to backup the `configs` directory to a safe place before doing anything with the node.
