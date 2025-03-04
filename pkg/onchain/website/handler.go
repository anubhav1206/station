package website

import (
	"net/http"
	"strings"

	"github.com/massalabs/station/api/interceptor"
	"github.com/massalabs/station/pkg/config"
	"github.com/massalabs/station/pkg/node"
	"github.com/massalabs/station/pkg/onchain/dns"
)

func handleMassaDomainRequest(writer http.ResponseWriter, reader *http.Request, index int, config config.AppConfig) {
	name := reader.Host[:index]

	rpcClient := node.NewClient(config.NodeURL)

	addr, err := dns.Resolve(config, rpcClient, name)
	if err != nil {
		panic(err)
	}

	var target string
	if reader.URL.Path == "/" {
		target = "index.html"
	} else {
		target = reader.URL.Path[1:]
	}

	Request(writer, reader, rpcClient, addr, target)
}

// MassaTLDInterceptor intercepts request for web on-chain.
func MassaTLDInterceptor(req *interceptor.Interceptor, appConfig config.AppConfig) *interceptor.Interceptor {
	if req == nil {
		return nil
	}

	massaIndex := strings.Index(req.Request.Host, ".massa")

	if massaIndex > 0 && !strings.HasPrefix(req.Request.Host, config.MassaStationURL) {
		handleMassaDomainRequest(req.Writer, req.Request, massaIndex, appConfig)

		return nil
	}

	return req
}
