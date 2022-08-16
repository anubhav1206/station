package wallet

import (
	"encoding/json"
	"io/ioutil"
	"sync"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/wallet"
)

func NewGetOne(walletStorage *sync.Map) operations.MgmtWalletGetOneHandler {
	return &walletGetOne{walletStorage: walletStorage}
}

type walletGetOne struct {
	walletStorage *sync.Map
}

// TODO Clean the struct mapping here + correct KeyPairs not returned & Panic(error)
func (c *walletGetOne) Handle(params operations.MgmtWalletGetOneParams) middleware.Responder {

	bytesInput, err := ioutil.ReadFile("wallet_" + params.Nickname + ".json")
	errorCode := ""
	msg := ""

	if err != nil {
		errorCode = errorCodeWalletFileNotFound
		msg = "Error: No wallet " + "wallet_" + params.Nickname + ".json found."
		return operations.NewMgmtWalletGetOneBadRequest().WithPayload(&models.Error{
			Code:    &errorCode,
			Message: &msg,
		})
	}

	wallet := wallet.Wallet{}
	err = json.Unmarshal(bytesInput, &wallet)
	if err != nil {
		errorCode = errorCodeWalletCorruptedFile
		msg = "Error: File wallet_" + params.Nickname + ".json corrupted."
		return operations.NewMgmtWalletGetOneBadRequest().WithPayload(&models.Error{
			Code:    &errorCode,
			Message: &msg,
		})
	}

	modelWallet := &models.Wallet{
		Nickname: &wallet.Nickname,
		Address:  &wallet.Address,
		KeyPairs: []*models.WalletKeyPairsItems0{}}

	return operations.NewMgmtWalletGetOneOK().WithPayload(modelWallet)
}
