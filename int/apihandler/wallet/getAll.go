package wallet

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/wallet"
)

func NewGet(walletStorage *sync.Map) operations.MgmtWalletGetHandler {
	return &walletGet{walletStorage: walletStorage}
}

type walletGet struct {
	walletStorage *sync.Map
}

// TODO Clean the struct mapping here + correct KeyPairs not returned & Panic(error)
func (c *walletGet) Handle(params operations.MgmtWalletGetParams) middleware.Responder {

	errorCode := ""
	msg := ""
	wd, err := os.Getwd()

	if err != nil {
		return operations.NewMgmtWalletGetInternalServerError()
	}

	wallets := []wallet.Wallet{}
	files, err := ioutil.ReadDir(wd)
	if err != nil {
		return operations.NewMgmtWalletGetInternalServerError()
	}

	for _, f := range files {
		fileName := f.Name()
		if strings.HasPrefix(fileName, "wallet_") && strings.HasSuffix(fileName, ".json") {
			bytesInput, err := ioutil.ReadFile(fileName)
			if err != nil {
				errorCode = errorCodeWalletFileNotFound
				msg = "Error: No wallet " + fileName + " found."
				return operations.NewMgmtWalletGetInternalServerError().WithPayload(&models.Error{
					Code:    &errorCode,
					Message: &msg,
				})
			}
			wallet := wallet.Wallet{}
			err = json.Unmarshal(bytesInput, &wallet)
			if err != nil {
				errorCode = errorCodeWalletCorruptedFile
				msg = "Error: File " + fileName + " corrupted."
				return operations.NewMgmtWalletGetInternalServerError().WithPayload(&models.Error{
					Code:    &errorCode,
					Message: &msg,
				})
			}
			wallets = append(wallets, wallet)
		}
	}

	var walll []*models.Wallet

	for i := 0; i < len(wallets); i++ {
		walletss := &models.Wallet{
			Nickname: &wallets[i].Nickname,
			Address:  &wallets[i].Address,
			KeyPairs: []*models.WalletKeyPairsItems0{}}
		walll = append(walll, walletss)
	}

	return operations.NewMgmtWalletGetOK().WithPayload(walll)
}
