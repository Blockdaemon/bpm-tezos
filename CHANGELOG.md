# 1.3.3


### Bug Fixes

* use the official term `full` instead of `watcher` ([1b74877](https://gitlab.com/Blockdaemon/bpm-tezos/commit/1b748774ad3b4c140b30368294c318e25a8d0256))



# 1.3.2


### Bug Fixes

* creating identity sometimes timed out ([b47afd7](https://gitlab.com/Blockdaemon/bpm-tezos/commit/b47afd72dad7df6a39a1e1a48387279668d02df1))



# 1.3.1


### Bug Fixes

* now supports `--data-dir` - it was hardcoded to `data` ([4c396e9](https://gitlab.com/Blockdaemon/bpm-tezos/commit/4c396e954132f199fc2bb15b998dc5100c625d6c))



# 1.3.0


### Features

* upgrade to tezos 7.1 ([6be130c](https://gitlab.com/Blockdaemon/bpm-tezos/commit/6be130c430bcc0e7b20ce7df0ccf9e6fb48c0f39))



# 1.2.0

### Features

* support monitoring packs ([8e44ecd](https://gitlab.com/Blockdaemon/bpm-tezos/commit/8e44ecdbf66fb1f40bb7e1e3d45afa7e5912ba5d))
* upgrade to tezos 7.0 ([92be16c](https://gitlab.com/Blockdaemon/bpm-tezos/commit/92be16ccdf6d0f31efa30945bf965f71bf97c735))

# 1.1.3

Bug fixes:

* Removed unsupported networks from help text (the fix in `1.1.1` only removed them from the validation, not from the help description)

# 1.1.2

Changes:

* `--data-dir` parameter allows to specify where to store the node data

Bug fixes:

* Now correctly validates the `--network` parameter
* Fix checkpoint test, it kept failing on mainnet

# 1.1.1

Bug fixes:

* Removed unsupported networks from help text

# 1.1.0

Changes:

* Dropped support for babylonnet because it is not active anymore
* Dropped support for zeronet because it is a network for core developers, not end-users

Bug fixes:

* Successfull test used to return an error, this is fixed now

