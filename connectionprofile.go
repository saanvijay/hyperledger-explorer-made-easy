package explorerutils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type ConnectionProfile struct {
	Name                   string                 `json:"name"`
	Version                string                 `json:"version"`
	License                string                 `json:"license"`
	Client                 Client                 `json:"client"`
	Channels               Channels               `json:"channels"`
	Organizations          Organizations          `json:"organizations"`
	Peers                  Peers                  `json:"peers"`
	CertificateAuthorities CertificateAuthorities `json:"certificateAuthorities"`
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
	Endorser string `json:"endorser"`
}
type Timeout struct {
	Peer    Peer   `json:"peer"`
	Orderer string `json:"orderer"`
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
type Peer0SupplierCom struct {
}
type Peers struct {
	Peer0SupplierCom Peer0SupplierCom `json:"peer0.supplier.com"`
}
type Testchannel struct {
	Peers      Peers      `json:"peers"`
	Connection Connection `json:"connection"`
}
type Channels struct {
	Testchannel Testchannel `json:"testchannel"`
}
type AdminPrivateKey struct {
	Path string `json:"path"`
}
type SignedCert struct {
	Path string `json:"path"`
}
type SupplierMSP struct {
	Mspid           string          `json:"mspid"`
	AdminPrivateKey AdminPrivateKey `json:"adminPrivateKey"`
	Peers           []string        `json:"peers"`
	SignedCert      SignedCert      `json:"signedCert"`
}
type Organizations struct {
	SupplierMSP SupplierMSP `json:"supplierMSP"`
}
type HTTPOptions struct {
	Verify bool `json:"verify"`
}
type TLSCACerts struct {
	Path string `json:"path"`
}
type SupplierCaCom struct {
	URL         string      `json:"url"`
	HTTPOptions HTTPOptions `json:"httpOptions"`
	TLSCACerts  TLSCACerts  `json:"tlsCACerts"`
	CaName      string      `json:"caName"`
}
type CertificateAuthorities struct {
	SupplierCaCom SupplierCaCom `json:"supplier.ca.com"`
}

func (configInput *ExplorerInput) GenerateConectionProfile() {
	fid, err := os.Create(configInput.NetworkName)
	if err != nil {
		log.Fatalf("Unable to create file %s\n", configInput.NetworkName)
		return
	}
	defer fid.Close()
	var connectionProfile ConnectionProfile

	var buffer bytes.Buffer
	err = json.NewEncoder(&buffer).Encode(connectionProfile)
	if err != nil {
		fmt.Println(err)
	}
	fid.Write([]byte(buffer.String()))
}
