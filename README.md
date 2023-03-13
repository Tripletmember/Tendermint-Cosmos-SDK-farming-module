[![codecov](https://codecov.io/gh/tendermint/farming/branch/main/graph/badge.svg)](https://codecov.io/gh/tendermint/farming?branch=main)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/tendermint/farming)](https://pkg.go.dev/github.com/tendermint/farming)

# Farming Module

The farming module is a Cosmos SDK module that implements farming functionality, which provides farming rewards to participants called farmers. A primary use case is to use this module to provide incentives for liquidity pool investors for their pool participation. 

- see the [main](https://github.com/tendermint/farming/tree/main) branch for the latest 
- see [releases](https://github.com/tendermint/farming/releases) for the latest release

## Dependencies

If you haven't already, install Golang by following the [official docs](https://golang.org/doc/install). Make sure that your `GOPATH` and `GOBIN` environment variables are properly set up.

Requirement | Notes
----------- | -----------------
Go version  | Go1.16 or higher
Cosmos SDK  | v0.44.5 or higher

## Installation

```bash
# Use git to clone farming module source code and install `farmingd`
git clone https://github.com/tendermint/farming.git
cd farming
make install
```

## Getting Started

To get started to the project, visit the [TECHNICAL-SETUP.md](./TECHNICAL-SETUP.md) docs.

