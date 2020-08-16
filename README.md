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
  1. git clone https://github.com/saanvijay/hyperledger-explorer-made-easy.git
  2. cp package/* $GOPATH/src/explorerutils/*
  3. Make sure you have running Hyperledger Fabric network
  4. By default all the explorer out files will be generated in "/tmp" dir, but you can set env variable to override it (export EXPLORER_OUT_CONFIG_PATH=/your/explorer/out/path)
  5. cd test
  6. Edit input fields as per your requirements (explorerinput.json)
  7. go run launchExplorer.go (wait for couple of mins)
  8. Open browser and type "localhost:8080" (default port)

## Written by

Vijaya Prakash<br />
https://www.linkedin.com/in/saanvijay/<br />
