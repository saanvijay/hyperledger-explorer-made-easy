package explorerutils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type ExplorerConfig struct {
	NetworkConfigs map[string]NetworkConfigs `json:"network-configs"`
	License        string                    `json:"license"`
}
type NetworkConfigs struct {
	Name    string `json:"name"`
	Profile string `json:"profile"`
}

func (configInput *ExplorerInput) GenerateExplorerConfig() {
	fid, err := os.Create("config.json")
	if err != nil {
		log.Fatalf("Unable to create file config.json")
		return
	}
	defer fid.Close()

	var explorerConfig ExplorerConfig
	var configMap map[string]NetworkConfigs
	var networkConfigs NetworkConfigs

	configMap = make(map[string]NetworkConfigs)
	networkConfigs.Name = configInput.NetworkName
	networkConfigs.Profile = fmt.Sprintf("./connection-profile/%s.json", configInput.NetworkName)
	configMap[configInput.NetworkName] = networkConfigs
	explorerConfig.NetworkConfigs = configMap
	explorerConfig.License = fmt.Sprintf("Apache-2.0")

	var buffer bytes.Buffer
	err = json.NewEncoder(&buffer).Encode(explorerConfig)
	if err != nil {
		fmt.Println(err)
	}
	fid.Write([]byte(buffer.String()))
}
