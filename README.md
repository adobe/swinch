# Swinch

<img src="./img/winch.jpeg" width="400"  />

Manage your Spinnaker pipelines with Kubernetes manifest and objects.
Swinch is a CLI tool that aims at functionality similar to kubectl and helm, but for Spinnaker.

Our goal is to make using Spinnaker friendly for users already familiar with the Kubernetes way of deploying by provide the same language and format used for managing Kubernetes asset, but for Spinnaker.  
Create, delete, edit and manage your Spinnaker CD pipelines reliably, programmatically and sourced controlled, from Helm like templated charts, with built in support for dry-run, diff and validation.

## Install

### Dev setup
Install go for your [platform](https://golang.org/doc/install)  
Set up your GO env, example:
```bash
SWINCH_REPO=git/swinch
export GOPATH=$HOME/go:$HOME/SWINCH_REPO
export GOBIN=$HOME/go/bin
export PATH=${PATH}:$GOBIN
```

Install swinch
```bash
SWINCH_REPO=$HOME/git/swinch
go install
```

### Shell completion
To get shell completion instructions for bash and zsh run:
```bash
swinch completion -h
```

## Basic usage
Generate manifest from the Chart and apply the resulting manifests

```bash
swinch template -c samples/charts/demo -o ./generated 
swinch apply -f ./generated
```
