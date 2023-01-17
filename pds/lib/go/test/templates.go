package test

import (
	"bytes"
	"fmt"
	"github.com/onflow/flow-go-sdk"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"text/template"
)

// Handle relative paths by making these regular expressions
const (
	nftAddressPlaceholder      = "\"[^\"]*NonFungibleToken.cdc\""
	IPackNFTAddressPlaceholder = "\"[^\"]*IPackNFT.cdc\""

	PackNftPath  = "../../../contracts/PackNFT.cdc"
	PDSPath      = "../../../contracts/PDS.cdc"
	IPackNftPath = "../../../contracts/IPackNFT.cdc"

	DeployPackNftPath = "../../../transactions/deploy/deploy-packNFT-with-auth.cdc"
	DeployPDSPath     = "../../../transactions/deploy/deploy-pds-with-auth.cdc"
)

func LoadPackNFT(nftAddress flow.Address, ipackAddress flow.Address) []byte {
	code := readFile(PackNftPath)

	nftRe := regexp.MustCompile(nftAddressPlaceholder)
	code = nftRe.ReplaceAll(code, []byte("0x"+nftAddress.String()))

	ipackNftRe := regexp.MustCompile(IPackNFTAddressPlaceholder)
	code = ipackNftRe.ReplaceAll(code, []byte("0x"+ipackAddress.String()))

	return code
}

func LoadPDS(nftAddress flow.Address, ipackAddress flow.Address) []byte {
	code := readFile(PDSPath)

	nftRe := regexp.MustCompile(nftAddressPlaceholder)
	code = nftRe.ReplaceAll(code, []byte("0x"+nftAddress.String()))

	ipackNftRe := regexp.MustCompile(IPackNFTAddressPlaceholder)
	code = ipackNftRe.ReplaceAll(code, []byte("0x"+ipackAddress.String()))

	return code
}

func LoadIPackNFT(nftAddress flow.Address) []byte {
	code := readFile(IPackNftPath)

	nftRe := regexp.MustCompile(nftAddressPlaceholder)
	code = nftRe.ReplaceAll(code, []byte("0x"+nftAddress.String()))
	return code
}

func ParseCadenceTemplate(templatePath string, data interface{}) ([]byte, error) {
	fmt.Println(filepath.Abs(templatePath))
	fb, err := ioutil.ReadFile(templatePath) /* #nosec */
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New("Template").Parse(string(fb))
	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, data)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func LoadPackNFTDeployScript() ([]byte, error) {
	return ParseCadenceTemplate(DeployPackNftPath, struct {
	}{})
}

func LoadPDSDeployScript() ([]byte, error) {
	return ParseCadenceTemplate(DeployPDSPath, struct {
	}{})
}
