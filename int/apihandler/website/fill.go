package websites

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"strings"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/node/base58"
	"github.com/massalabs/thyra/pkg/node/getters"
	"github.com/massalabs/thyra/pkg/node/sendoperation"
	"github.com/massalabs/thyra/pkg/node/sendoperation/callsc"
	"github.com/massalabs/thyra/pkg/wallet"
)

func NewFillWebsitePost() operations.FillWebPostHandler {
	return &fillWeb{todelete: "todelete"}
}

type fillWeb struct {
	todelete string
}

type Bytes struct {
	Data string `json:"data"`
}

func (c *fillWeb) Handle(params operations.FillWebPostParams) middleware.Responder {
	errorCode := ""
	msg := ""

	client := node.NewClient()
	wallet, err := wallet.GetWallet(params.Nickname)

	if err != nil {
		errorCode = errorCodeWebsiteGetWallet
		msg = "Error: get Wallet failed"
		return operations.NewFillWebPostInternalServerError().WithPayload(&models.Error{
			Code:    &errorCode,
			Message: &msg,
		})
	}
	status, err := getters.GetNodeStatus(client)
	if err != nil {
		errorCode = errorCodeWebsiteGetWallet
		msg = "Error: get Node status failed"
		return operations.NewFillWebPostInternalServerError().WithPayload(&models.Error{
			Code:    &errorCode,
			Message: &msg,
		})
	}

	address, _, err := base58.VersionedCheckDecode(params.Website[1:])
	if err != nil {
		errorCode = errorCodeWebsiteGetWallet
		msg = "Error: decode website deployer address failed"
		return operations.NewFillWebPostInternalServerError().WithPayload(&models.Error{
			Code:    &errorCode,
			Message: &msg,
		})
	}

	buf := new(strings.Builder)
	_, err = io.Copy(buf, params.Zipfile)
	if err != nil {
		errorCode = errorCodeWebsiteGetWallet
		msg = "Error: get Node status failed"
		return operations.NewFillWebPostInternalServerError().WithPayload(&models.Error{
			Code:    &errorCode,
			Message: &msg,
		})
	}
	sEnc := base64.StdEncoding.EncodeToString([]byte(buf.String()))

	d, err := json.Marshal(Bytes{
		Data: sEnc,
	})
	if err != nil {
		errorCode = errorCodeWebsiteGetWallet
		msg = "Error: serialize website data"
		return operations.NewFillWebPostInternalServerError().WithPayload(&models.Error{
			Code:    &errorCode,
			Message: &msg,
		})
	}

	callSC := callsc.New(address, "initializeWebsite", d, 0, 700000000, 0, 0)
	id, err := sendoperation.Call(client, uint64(status.NextSlot.Period+2), 0, callSC, wallet.KeyPairs[0].PublicKey, wallet.KeyPairs[0].PrivateKey)

	if err != nil {
		errorCode = errorCodeWebsiteGetWallet
		msg = "Error: insert first chunk in the Smart contract"
		return operations.NewFillWebPostInternalServerError().WithPayload(&models.Error{
			Code:    &errorCode,
			Message: &msg,
		})
	}
	b := false
	n := 0
	for n < 3 && !b {

		time.Sleep(10 * time.Second)
		events, err := getters.GetEvents(client, nil, nil, nil, nil, &id)
		if err != nil {
			errorCode = errorCodeWebsiteGetWallet
			msg = "Error: get events "
			return operations.NewFillWebPostInternalServerError().WithPayload(&models.Error{
				Code:    &errorCode,
				Message: &msg,
			})
		}

		eventsValue := *events
		if len(eventsValue) > 0 {
			b = true
		}
		n++

	}
	website := &models.Websites{
		Name:    "Name",
		Address: params.Website}

	return operations.NewFillWebPostOK().WithPayload(website)
}
