package websites

import (
	"encoding/json"
	"os"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/node/getters"
	"github.com/massalabs/thyra/pkg/onchain/dns"
	"github.com/massalabs/thyra/pkg/onchain/website"
	"github.com/massalabs/thyra/pkg/wallet"
)

func NewWebsitePost() operations.UploadWebPostHandler {
	return &newWebsitePost{todelete: "todelete"}
}

type newWebsitePost struct {
	todelete string
}

func (c *newWebsitePost) Handle(params operations.UploadWebPostParams) middleware.Responder {
	client := node.NewClient()
	errorCode := ""
	msg := ""
	// Get status for expire period
	expirePeriod, err := getters.GetExpirePeriod(client)
	if err != nil {
		errorCode = errorCodeWebsiteGetExpirePeriod
		msg = "Error: Get Node status failed"
		return operations.NewUploadWebPostInternalServerError().WithPayload(&models.Error{
			Code:    &errorCode,
			Message: &msg,
		})
	}

	// Get first wallet
	wallet, err := wallet.GetWallet(params.Nickname)
	if err != nil {
		errorCode = errorCodeWebsiteGetWallet
		msg = "Error: Get Wallet failed"
		return operations.NewUploadWebPostInternalServerError().WithPayload(&models.Error{
			Code:    &errorCode,
			Message: &msg,
		})
	}

	// Deploy Smart contract deployer
	smartContract, err := website.DeployWebsiteDeployer(client, *wallet, *expirePeriod)
	if err != nil {
		errorCode = errorCodeWebsiteDeployWebsiteDeployer
		msg = "Error: deploying website deployer Smart Contract"
		return operations.NewUploadWebPostInternalServerError().WithPayload(&models.Error{
			Code:    &errorCode,
			Message: &msg,
		})
	}

	// // Set DNS Approval
	// _, err = dns.SetDnsApproval(c, *wallet, true, *expirePeriod)
	// if err != nil {
	// 	return nil, err
	// }
	// time.Sleep(15 * time.Second)

	// Set DNS Resolver
	_, err = dns.SetDnsResolver(client, *wallet, params.Dnsname, *smartContract, *expirePeriod)
	if err != nil {
		errorCode = errorCodeWebsiteSetDnsResolver
		msg = "Error: Setting the DNS resolver failed"
		return operations.NewUploadWebPostInternalServerError().WithPayload(&models.Error{
			Code:    &errorCode,
			Message: &msg,
		})
	}
	deployers, err := website.GetDeployers()
	if err != nil {
		errorCode = errorCodeWebsiteGetDeployers
		msg = "Error: Get website deployers failed"
		return operations.NewUploadWebPostInternalServerError().WithPayload(&models.Error{
			Code:    &errorCode,
			Message: &msg,
		})
	}
	dep := website.WebsiteDeployer{
		DnsName: &params.Dnsname,
		Address: smartContract,
	}
	deployers = append(deployers, dep)

	bytesOutput, err := json.Marshal(deployers)
	if err != nil {
		errorCode = errorCodeWebsiteDeserializeDeployers
		msg = "Error: Deserialize website deployers filed failed"
		return operations.NewUploadWebPostInternalServerError().WithPayload(&models.Error{
			Code:    &errorCode,
			Message: &msg,
		})
	}

	err = os.WriteFile("deployers.json", bytesOutput, 0o644)
	if err != nil {
		errorCode = errorCodeWebsiteSerializeDeployers
		msg = "Error: Writting of new deployers file failed"
		return operations.NewUploadWebPostInternalServerError().WithPayload(&models.Error{
			Code:    &errorCode,
			Message: &msg,
		})
	}

	newWebsite := &models.Websites{
		Name:    params.Dnsname,
		Address: *smartContract}

	return operations.NewUploadWebPostOK().WithPayload(newWebsite)
}
