package explorerutils

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

type ExplorerInput struct {
	NetworkName           string
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
	CAPort                int
	ExplorerOutConfigPath string
}

func (configInput *ExplorerInput) ExplorerDown() {
	cmd := exec.Command("docker-compose", "-f", "docker-compose-explorer.yaml", "down", "-v")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ":compose down error --------- " + string(output))
		log.Fatalf(string(output))
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

func (configInput *ExplorerInput) LaunchExplorer() {

	OriginalDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Unable to get CWD : %s\n", err)
	}

	if os.Getenv("EXPLORER_OUT_CONFIG_PATH") == "" {
		os.Setenv("EXPLORER_OUT_CONFIG_PATH", os.TempDir())
	}
	configInput.ExplorerOutConfigPath = os.Getenv("EXPLORER_OUT_CONFIG_PATH")
	explorerPath := fmt.Sprintf("%s/%s/explorer", configInput.ExplorerOutConfigPath, configInput.NetworkName)
	os.MkdirAll(explorerPath, 0664)
	os.Chdir(explorerPath)
	configInput.GenerateExplorerConfig()
	configInput.GenerateDockerCompose()
	os.Mkdir("connection-profile", 0664)
	os.Chdir("connection-profile")
	configInput.GenerateConectionProfile()
	os.Chdir("..")

	configInput.ExplorerDown()
	configInput.ExplorerUp()

	os.Chdir(OriginalDir)
}