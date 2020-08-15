package explorerutils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"
)

type ExplorerInput struct {
	NetworkName           string
	ExternalNetworkName   string
	ChannelName           string
	DiscoverAsLocalHost   bool
	CryptoConfigPath      string
	ExplorerPort          int
	TLSEnable             bool
	AdminUserName         string
	AdminPassword         string
	Organization          string
	PeerID                string
	PeerPort              int
	CAName                string
	CAPort                int
	ExplorerOutConfigPath string
}

func (configInput *ExplorerInput) ExplorerDown() {
	cmd := exec.Command("docker-compose", "-f", "docker-compose-explorer.yaml", "down", "-v")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ":compose down error --------- " + string(output))
		//	log.Fatalf(string(output))
	}
	time.Sleep(30 * time.Second)
}

func (configInput *ExplorerInput) ExplorerUp() {
	cmd := exec.Command("docker-compose", "-f", "docker-compose-explorer.yaml", "up", "-d")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ":compose up error --------- " + string(output))
		log.Fatalf(string(output))
	}
	time.Sleep(30 * time.Second)
}

func LaunchExplorer(jsoninput string) {

	var configInput ExplorerInput

	inputFile, err := os.Open(jsoninput)
	if err != nil {
		fmt.Printf("Unable to open the json input file: %s\n", err)
		return
	}
	defer inputFile.Close()

	jsonbytes, _ := ioutil.ReadAll(inputFile)

	err = json.Unmarshal([]byte(jsonbytes), &configInput)
	if err != nil {
		fmt.Printf("Unable to unmarshal json %s\n", err)
		return
	}

	OriginalDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Unable to get CWD : %s\n", err)
	}

	if os.Getenv("EXPLORER_OUT_CONFIG_PATH") == "" {
		os.Setenv("EXPLORER_OUT_CONFIG_PATH", "/tmp/")
	}
	configInput.ExplorerOutConfigPath = os.Getenv("EXPLORER_OUT_CONFIG_PATH")
	explorerPath := fmt.Sprintf("%s/%s/explorer", configInput.ExplorerOutConfigPath, configInput.NetworkName)
	err = os.MkdirAll(explorerPath, 0777)
	if err != nil {
		fmt.Printf("Unable to create directory : %s\n", err)
		return
	}
	err = os.Chdir(explorerPath)
	if err != nil {
		fmt.Printf("Unable to change directory : %s\n", err)
		return
	}
	configInput.GenerateExplorerConfig()
	configInput.GenerateDockerCompose()
	err = os.Mkdir("connection-profile", 0777)
	if err != nil {
		fmt.Printf("Unable to create directory : %s\n", err)
		return
	}
	err = os.Chdir("connection-profile")
	if err != nil {
		fmt.Printf("Unable to change directory : %s\n", err)
		return
	}
	configInput.GenerateConectionProfile()
	err = os.Chdir("..")
	if err != nil {
		fmt.Printf("Unable to change directory : %s\n", err)
		return
	}
	configInput.ExplorerDown()
	configInput.ExplorerUp()

	os.Chdir(OriginalDir)
}
