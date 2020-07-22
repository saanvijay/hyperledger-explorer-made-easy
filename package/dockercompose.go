package explorerutils

import (
	"bytes"
	"fmt"
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type ExplorerServices struct {
	Version  string                     `yaml:"version"`
	Volumes  map[string]interface{}     `yaml:"volumes"`
	Networks map[string]Networks        `yaml:"networks"`
	Services map[string]ExplorerService `yaml:"services"`
}
type External struct {
	Name string `yaml:"name"`
}
type Networks struct {
	External External `yaml:"external"`
}
type Healthcheck struct {
	Test     string `yaml:"test"`
	Interval string `yaml:"interval"`
	Timeout  string `yaml:"timeout"`
	Retries  int    `yaml:"retries"`
}

type DependsOn struct {
	Condition string `yaml:"condition"`
}
type ExplorerService struct {
	Image         string               `yaml:"image,omitempty"`
	ContainerName string               `yaml:"container_name,omitempty"`
	Hostname      string               `yaml:"hostname,omitempty"`
	Environment   []string             `yaml:"environment,omitempty"`
	Healthcheck   Healthcheck          `yaml:"healthcheck,omitempty"`
	Volumes       []string             `yaml:"volumes,omitempty"`
	Command       string               `yaml:"command,omitempty"`
	Ports         []string             `yaml:"ports,omitempty"`
	DependsOn     map[string]DependsOn `yaml:"depends_on,omitempty"`
	Networks      []string             `yaml:"networks,omitempty"`
}

func (configInput *ExplorerInput) GenerateDockerCompose() {
	fid, err := os.Create("docker-compose-explorer.yaml")
	if err != nil {
		log.Fatalf("Unable to create file docker-compose-explorer.yaml")
		return
	}
	defer fid.Close()
	// Version
	var explorerServices ExplorerServices
	explorerServices.Version = "2.1"

	//Volumes
	var volumeMap map[string]interface{}
	volumeMap = make(map[string]interface{})
	volumeMap["pgdata"] = nil
	volumeMap["walletstore"] = nil
	explorerServices.Volumes = volumeMap

	// External networks
	var networks Networks
	var networkMap map[string]Networks
	networkMap = make(map[string]Networks)
	networks.External.Name = fmt.Sprintf("organizations_%s", configInput.NetworkName)
	networkMap[fmt.Sprintf("%s.com", configInput.NetworkName)] = networks
	explorerServices.Networks = networkMap

	// Services
	var explorerServiceMap map[string]ExplorerService
	explorerServiceMap = make(map[string]ExplorerService)

	// DB Service
	var explorerService1 ExplorerService
	explorerService1.Image = "hyperledger/explorer-db:latest"
	explorerService1.ContainerName = fmt.Sprintf("explorerdb.%s.com", configInput.NetworkName)
	explorerService1.Hostname = fmt.Sprintf("explorerdb.%s.com", configInput.NetworkName)
	explorerService1.Environment = []string{
		"DATABASE_DATABASE=fabricexplorer",
		"DATABASE_USERNAME=hppoc",
		"DATABASE_PASSWORD=password",
	}
	explorerService1.Healthcheck.Test = "pg_isready -h localhost -p 5432 -q -U postgres"
	explorerService1.Healthcheck.Interval = "30s"
	explorerService1.Healthcheck.Timeout = "10s"
	explorerService1.Healthcheck.Retries = 5

	explorerService1.Volumes = []string{
		"pgdata:/var/lib/postgresql/data",
	}
	explorerService1.Networks = []string{
		fmt.Sprintf("%s.com", configInput.NetworkName),
	}
	// Explorer Service
	var explorerService2 ExplorerService
	explorerService2.Image = "hyperledger/explorer:latest"
	explorerService2.ContainerName = fmt.Sprintf("explorer.%s.com", configInput.NetworkName)
	explorerService2.Hostname = fmt.Sprintf("explorer.%s.com", configInput.NetworkName)
	explorerService2.Environment = []string{
		fmt.Sprintf("DATABASE_HOST=explorerdb.%s.com", configInput.NetworkName),
		"DATABASE_DATABASE=fabricexplorer",
		"DATABASE_USERNAME=hppoc",
		"DATABASE_PASSWD=password",
		"LOG_LEVEL_APP=debug",
		"LOG_LEVEL_DB=debug",
		"LOG_LEVEL_CONSOLE=info",
		"LOG_CONSOLE_STDOUT=true",
		fmt.Sprintf("DISCOVERY_AS_LOCALHOST=%t", configInput.DiscoverAsLocalHost),
	}
	explorerService2.Volumes = []string{
		"./config.json:/opt/explorer/app/platform/fabric/config.json",
		"./connection-profile:/opt/explorer/app/platform/fabric/connection-profile",
		fmt.Sprintf("%s:/tmp/crypto", configInput.CryptoConfigPath),
		"walletstore:/opt/wallet",
	}
	explorerService2.Command = "sh -c \"node /opt/explorer/main.js && tail -f /dev/null\""
	explorerService2.Ports = []string{
		fmt.Sprintf("%d:8080", configInput.ExplorerPort),
	}

	var dependsOnMap map[string]DependsOn
	dependsOnMap = make(map[string]DependsOn)
	var dependsOn DependsOn
	dependsOn.Condition = "service_healthy"
	dependsOnMap[fmt.Sprintf("explorerdb.%s.com", configInput.NetworkName)] = dependsOn
	explorerService2.DependsOn = dependsOnMap
	explorerService2.Networks = []string{
		fmt.Sprintf("%s.com", configInput.NetworkName),
	}

	explorerServiceMap[fmt.Sprintf("explorerdb.%s.com", configInput.NetworkName)] = explorerService1
	explorerServiceMap[fmt.Sprintf("explorer.%s.com", configInput.NetworkName)] = explorerService2
	explorerServices.Services = explorerServiceMap

	var buffer bytes.Buffer
	err = yaml.NewEncoder(&buffer).Encode(explorerServices)
	if err != nil {
		fmt.Println(err)
	}
	fid.Write([]byte(buffer.String()))
}
