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

