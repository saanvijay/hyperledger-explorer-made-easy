package explorerutils

import (
	"bytes"
	"fmt"
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type ExplorerServices struct {
	Version  string   `yaml:"version"`
	Volumes  Volumes  `yaml:"volumes"`
	Networks Networks `yaml:"networks"`
	Services Services `yaml:"services"`
}
type Volumes struct {
	Pgdata      interface{} `yaml:"pgdata"`
	Walletstore interface{} `yaml:"walletstore"`
}
type External struct {
	Name string `yaml:"name"`
}
type TestnetworkCom struct {
	External External `yaml:"external"`
}
type Networks struct {
	TestnetworkCom TestnetworkCom `yaml:"testnetwork.com"`
}
type Healthcheck struct {
	Test     string `yaml:"test"`
	Interval string `yaml:"interval"`
	Timeout  string `yaml:"timeout"`
	Retries  int    `yaml:"retries"`
}
type ExplorerdbTestnetworkCom struct {
	Image         string      `yaml:"image"`
	ContainerName string      `yaml:"container_name"`
	Hostname      string      `yaml:"hostname"`
	Environment   []string    `yaml:"environment"`
	Healthcheck   Healthcheck `yaml:"healthcheck"`
	Volumes       []string    `yaml:"volumes"`
	Networks      []string    `yaml:"networks"`
}
type ExplorerdbTestnetworkCom struct {
	Condition string `yaml:"condition"`
}
type DependsOn struct {
	ExplorerdbTestnetworkCom ExplorerdbTestnetworkCom `yaml:"explorerdb.testnetwork.com"`
}
type ExplorerTestnetworkCom struct {
	Image         string    `yaml:"image"`
	ContainerName string    `yaml:"container_name"`
	Hostname      string    `yaml:"hostname"`
	Environment   []string  `yaml:"environment"`
	Volumes       []string  `yaml:"volumes"`
	Command       string    `yaml:"command"`
	Ports         []string  `yaml:"ports"`
	DependsOn     DependsOn `yaml:"depends_on"`
	Networks      []string  `yaml:"networks"`
}
type Services struct {
	ExplorerdbTestnetworkCom ExplorerdbTestnetworkCom `yaml:"explorerdb.testnetwork.com"`
	ExplorerTestnetworkCom   ExplorerTestnetworkCom   `yaml:"explorer.testnetwork.com"`
}

func (configInput *ExplorerInput) GenerateDockerCompose() {
	fid, err := os.Create("docker-compose.yaml")
	if err != nil {
		log.Fatalf("Unable to create file docker-compose.yaml")
		return
	}
	defer fid.Close()
	var explorerServices ExplorerServices

	var buffer bytes.Buffer
	err = yaml.NewEncoder(&buffer).Encode(explorerServices)
	if err != nil {
		fmt.Println(err)
	}
	fid.Write([]byte(buffer.String()))
}
