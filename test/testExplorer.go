package main

import "explorerutils"

func main() {

	explorerInput := explorerutils.ExplorerInput{
		NetworkName:         "testnetwork",
		ChannelName:         "testchannel",
		DiscoverAsLocalHost: false,
		CryptoConfigPath:    "/tmp/testnetwork/crypto-config/",
		ExplorerPort:        8080,
		TLSEnable:           true,
		AdminUserName:       "exploreradmin",
		AdminPassword:       "exploreradminpw",
		Organization:        "supplier",
		PeerID:              "peer0.supplier.com",
		PeerPort:            7051,
		CAPort:              6054,
	}
	explorerInput.LaunchExplorer()
}
