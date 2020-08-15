package explorerutils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type ConnectionProfile struct {
	Name                   string                            `json:"name"`
	Version                string                            `json:"version"`
	License                string                            `json:"license"`
	Client                 Client                            `json:"client"`
	Channels               map[string]Channel                `json:"channels"`
	Organizations          map[string]Organizations          `json:"organizations"`
	Peers                  map[string]Peers                  `json:"peers"`
	CertificateAuthorities map[string]CertificateAuthorities `json:"certificateAuthorities"`
}
type CaCredential struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}
type AdminCredential struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}
type Peer struct {
	Endorser string `json:"endorser,omitempty"`
	EventHub string `json:"eventHub,omitempty"`
	EventReg string `json:"eventReg,omitempty"`
}
type Timeout struct {
	Peer    Peer   `json:"peer,omitempty"`
	Orderer string `json:"orderer,omitempty"`
}
type Connection struct {
	Timeout Timeout `json:"timeout"`
}
type Client struct {
	TLSEnable            bool            `json:"tlsEnable"`
	CaCredential         CaCredential    `json:"caCredential"`
	AdminCredential      AdminCredential `json:"adminCredential"`
	EnableAuthentication bool            `json:"enableAuthentication"`
	Organization         string          `json:"organization"`
	Connection           Connection      `json:"connection"`
}
type ChannelPeers struct {
	LedgerQuery bool `json:"ledgerQuery,omitempty"`
}
type Channel struct {
	ChannelPeers map[string]ChannelPeers `json:"peers"`
	Connection   Connection              `json:"connection"`
}
type AdminPrivateKey struct {
	Path string `json:"path"`
}
type SignedCert struct {
	Path string `json:"path"`
}
type Organizations struct {
	Mspid           string          `json:"mspid"`
	AdminPrivateKey AdminPrivateKey `json:"adminPrivateKey"`
	Peers           []string        `json:"peers"`
	SignedCert      SignedCert      `json:"signedCert"`
}
type HTTPOptions struct {
	Verify bool `json:"verify"`
}
type TLSCACerts struct {
	Path string `json:"path"`
}
type GRPCOptions struct {
	SSLtargetnameoverride string `json:"ssl-target-name-override"`
}
type Peers struct {
	TLSCACerts  TLSCACerts  `json:"tlsCACerts"`
	URL         string      `json:"url"`
	GRPCOptions GRPCOptions `json:"grpcOptions"`
}
type CertificateAuthorities struct {
	URL         string      `json:"url"`
	HTTPOptions HTTPOptions `json:"httpOptions"`
	TLSCACerts  TLSCACerts  `json:"tlsCACerts"`
	CaName      string      `json:"caName"`
}

func (configInput *ExplorerInput) GenerateConectionProfile() {
	fid, err := os.Create(fmt.Sprintf("%s.json", configInput.NetworkName))
	if err != nil {
		log.Fatalf("Unable to create file %s\n", configInput.NetworkName)
		return
	}
	defer fid.Close()

	var connectionProfile ConnectionProfile
	connectionProfile.Name = configInput.NetworkName
	connectionProfile.Version = "1.0.0"
	connectionProfile.License = "Apache-2.0"

	// Client Configuration
	var client Client
	client.TLSEnable = configInput.TLSEnable
	client.CaCredential.ID = "admin"
	client.CaCredential.Password = "adminpw"
	client.AdminCredential.ID = configInput.AdminUserName
	client.AdminCredential.Password = configInput.AdminPassword
	client.EnableAuthentication = true
	client.Organization = configInput.Organization
	client.Connection.Timeout.Orderer = "300"
	client.Connection.Timeout.Peer.Endorser = "300"
	connectionProfile.Client = client

	// Channel Configuration
	var channel Channel
	var channelsMap map[string]Channel
	var channelPeerMap map[string]ChannelPeers
	channelsMap = make(map[string]Channel)
	channelPeerMap = make(map[string]ChannelPeers)
	channelPeerMap[configInput.PeerID] = ChannelPeers{}
	channel.ChannelPeers = channelPeerMap
	channel.Connection.Timeout.Peer.Endorser = "6000"
	channel.Connection.Timeout.Peer.EventHub = "6000"
	channel.Connection.Timeout.Peer.EventReg = "6000"
	channelsMap[configInput.ChannelName] = channel
	connectionProfile.Channels = channelsMap

	// Organization Configurations
	var organization Organizations
	organization.Mspid = fmt.Sprintf("%sMSP", configInput.Organization)
	organization.AdminPrivateKey.Path = fmt.Sprintf("/tmp/crypto/peerOrganizations/%s.com/users/Admin@%s.com/msp/keystore/priv_sk", configInput.Organization, configInput.Organization)
	organization.Peers = []string{configInput.PeerID}
	organization.SignedCert.Path = fmt.Sprintf("/tmp/crypto/peerOrganizations/%s.com/users/Admin@%s.com/msp/signcerts/Admin@%s.com-cert.pem", configInput.Organization, configInput.Organization, configInput.Organization)
	var organizationMap map[string]Organizations
	organizationMap = make(map[string]Organizations)
	organizationMap[fmt.Sprintf(configInput.Organization)] = organization
	connectionProfile.Organizations = organizationMap

	// Peers Configuration
	var peers Peers
	var peerURL string
	if configInput.TLSEnable {
		peerURL = fmt.Sprintf("grpcs://%s:%d", configInput.PeerID, configInput.PeerPort)
	} else {
		peerURL = fmt.Sprintf("grpc://%s:%d", configInput.PeerID, configInput.PeerPort)
	}
	peers.URL = peerURL
	peers.GRPCOptions.SSLtargetnameoverride = configInput.PeerID
	peers.TLSCACerts.Path = fmt.Sprintf("/tmp/crypto/peerOrganizations/%s.com/peers/%s/tls/ca.crt", configInput.Organization, configInput.PeerID)
	var peersMap map[string]Peers
	peersMap = make(map[string]Peers)
	peersMap[configInput.PeerID] = peers
	connectionProfile.Peers = peersMap

	// CertificateAuthorities Configurations
	var certificateAuthorities CertificateAuthorities
	var caName string
	if configInput.CAName != "" {
		caName = configInput.CAName
	} else {
		caName = fmt.Sprintf("%s-ca-server", configInput.Organization)
	}
	certificateAuthorities.CaName = caName
	var caURL string
	if configInput.TLSEnable {
		caURL = fmt.Sprintf("https://%s.ca.com:%d", configInput.Organization, configInput.CAPort)
	} else {
		caURL = fmt.Sprintf("http://%s.ca.com:%d", configInput.Organization, configInput.CAPort)
	}
	certificateAuthorities.URL = caURL
	certificateAuthorities.HTTPOptions.Verify = false
	certificateAuthorities.TLSCACerts.Path = fmt.Sprintf("/tmp/crypto/peerOrganizations/%s.com/ca/ca.%s.com-cert.pem", configInput.Organization, configInput.Organization)
	var caMap map[string]CertificateAuthorities
	caMap = make(map[string]CertificateAuthorities)
	caMap[fmt.Sprintf("%s.ca.com", configInput.Organization)] = certificateAuthorities
	connectionProfile.CertificateAuthorities = caMap

	var buffer bytes.Buffer
	err = json.NewEncoder(&buffer).Encode(connectionProfile)
	if err != nil {
		fmt.Println(err)
	}
	fid.Write([]byte(buffer.String()))
}
