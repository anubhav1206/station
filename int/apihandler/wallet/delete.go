package wallet

import (
	"os"
	"sync"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
)

func NewDelete(walletStorage *sync.Map) operations.MgmtWalletDeleteHandler {
	return &walletDelete{walletStorage: walletStorage}
}

type walletDelete struct {
	walletStorage *sync.Map
}

func (c *walletDelete) Handle(params operations.MgmtWalletDeleteParams) middleware.Responder {
	if len(params.Nickname) == 0 {
		e := errorCodeWalletDeleteNoNickname
		msg := "Error: nickname field is mandatory."

		return operations.NewMgmtWalletDeleteBadRequest().WithPayload(
			&models.Error{
				Code:    &e,
				Message: &msg,
			})
	}

	err := os.Remove("wallet_" + params.Nickname + ".json")
	if err != nil {
		e := errorCodeWalletDeleteNoNickname
		msg := "Error: Can't delete wallet_" + params.Nickname + ".json file"

		return operations.NewMgmtWalletCreateInternalServerError().WithPayload(
			&models.Error{
				Code:    &e,
				Message: &msg,
			})
	}

	return operations.NewMgmtWalletDeleteNoContent()
}
