package websites

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station/api/swagger/server/models"
	"github.com/massalabs/station/api/swagger/server/restapi/operations"
	"github.com/massalabs/station/pkg/config"
	sendOperation "github.com/massalabs/station/pkg/node/sendoperation"
	"github.com/massalabs/station/pkg/onchain/website"
)

func NewWebsiteUploadMissedChunkHandler(config *config.AppConfig) operations.WebsiteUploadMissingChunksHandler {
	return &uploadMissedChunkHandler{config: config}
}

type uploadMissedChunkHandler struct {
	config *config.AppConfig
}

func (h *uploadMissedChunkHandler) Handle(params operations.WebsiteUploadMissingChunksParams) middleware.Responder {
	archive, errorResponse := readAndCheckArchive(params.Zipfile)
	if errorResponse != nil {
		return errorResponse
	}

	_, err := website.UploadMissedChunks(
		*h.config,
		params.Address,
		archive,
		params.Nickname,
		params.MissedChunks,
		sendOperation.OperationBatch{
			NewBatch:      true,
			CorrelationID: "",
		},
	)
	if err != nil {
		return createInternalServerError(errorCodeWebCreatorUpload, err.Error())
	}

	return operations.NewWebsiteUploadMissingChunksOK().
		WithPayload(&models.Websites{
			Name:         "",
			Description:  "",
			Address:      params.Address,
			BrokenChunks: nil,
		})
}
