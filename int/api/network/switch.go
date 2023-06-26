package network

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station/api/swagger/server/models"
	"github.com/massalabs/station/api/swagger/server/restapi/operations"
	"github.com/massalabs/station/pkg/config"
)

type switchNetworkHandler struct {
	networkManager *config.NetworkManager
}

// NewSwitchNetworkHandler creates a new switchNetworkHandler instance.
func NewSwitchNetworkHandler(networkManager *config.NetworkManager) operations.SwitchNetworkHandler {
	return &switchNetworkHandler{networkManager: networkManager}
}

// handles the request for switching the network.
func (h *switchNetworkHandler) Handle(params operations.SwitchNetworkParams) middleware.Responder {
	err := h.networkManager.SwitchNetwork(params.Network)
	if err != nil {
		// If the network is not found, return a 404 response with an error message.
		return operations.NewSwitchNetworkNotFound().WithPayload(
			&models.Error{
				Code:    "404",
				Message: "Network not found",
			},
		)
	}

	// Build the response with the current network information.
	response := &models.NetworkManagerItem{
		CurrentNetwork:     &h.networkManager.Network().Network,
		AvailableNetworks: *h.networkManager.Networks(),
	}

	return operations.NewSwitchNetworkOK().WithPayload(response)
}
