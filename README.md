# hyperledger-explorer-made-easy
![Platforms](https://img.shields.io/badge/platform-osx%20%7C%20linux-lightgray.svg)

Hyperledger explorer setup is now pretty easy. The explorerdb and explorer both functionality is shipped with docker images, however to bring-up the actual explorer again we have to do few additional steps like creating config files and yaml files. The
hyperledger-explorer-made-easy breaks all those file generation barriers for you. Just download all pre-requisites and modify input json file and then launch explorer, that's all.

# PreRequisites:
1. docker
2. docker-compose
3. golang
4. running hyperledger fabric network (supports latest version 2.2 as of now)

# How to test?
  1. Copy the package/* files in $GOPATH/src/explorerutils/*
  2. Make sure you have running Hyperledger Fabric network
  3. Go to test dir and edit input fields as per your requirements and run "go run launchExplorer.go" (wait for couple of mins)
  4. Open browser and type "localhost:8080" (default port)
