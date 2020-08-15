# hyperledger-explorer-made-easy

PreRequisites:
1. docker
2. docker-compose
3. golang
4. running hyperledger fabric network (supports latest version 2.2 as of now)

How to test?
  1. Copy the package/* files in $GOPATH/src/explorerutils/*
  2. Make sure you have running Hyperledger Fabric network
  3. Go to test dir and edit input fields as per your requirements and run "go run launchExplorer.go" (wait for couple of mins)
  4. Open browser and type "localhost:8080" (default port)
