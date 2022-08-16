package website

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"time"

	"github.com/massalabs/thyra/pkg/contracts"
	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/node/getters"
	sendOperation "github.com/massalabs/thyra/pkg/node/sendoperation"

	"github.com/massalabs/thyra/pkg/node/sendoperation/executesc"
	"github.com/massalabs/thyra/pkg/wallet"
)

func DeployWebsiteDeployer(c *node.Client, wallet wallet.Wallet, expire uint64) (*string, error) {
	exeSC := executesc.New([]byte(contracts.WebstiteDeployerContract), 700000, 0, 0)
	id, err := sendOperation.Call(c, expire, 0, exeSC, wallet.KeyPairs[0].PublicKey, wallet.KeyPairs[0].PrivateKey)
	if err != nil {
		return nil, err
	}

	// Get SC Contract
	smartContract := ""
	n := 0
	for n < 3 && smartContract == "" {

		time.Sleep(15 * time.Second)
		events, err := getters.GetEvents(c, nil, nil, nil, nil, &id)
		if err != nil {
			return nil, err
		}

		eventsValue := *events
		if len(eventsValue) > 0 {
			smartContract = strings.Split(eventsValue[0].Data, ":")[1]
		}
		n++

	}
	return &smartContract, nil

}

type WebsiteDeployer struct {
	DnsName *string `json:"dnsName"`
	Address *string `json:"address"`
}

func GetDeployers() ([]WebsiteDeployer, error) {
	deployers := []WebsiteDeployer{}
	bytesInput, err := ioutil.ReadFile("deployers.json")
	if err != nil {
		return deployers, nil
	}

	err = json.Unmarshal(bytesInput, &deployers)
	if err != nil {
		return nil, err
	}
	return deployers, nil
}
