package plugin

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/massalabs/station/api/interceptor"
	"github.com/massalabs/station/pkg/config"
)

func NewAPIHandler(manager *Manager) *APIHandler {
	return &APIHandler{manager: manager}
}

type APIHandler struct {
	manager *Manager
}

func (h *APIHandler) Handle(writer http.ResponseWriter, reader *http.Request, pluginAuthor string, pluginName string) {
	alias := Alias(pluginAuthor, pluginName)

	plugin, err := h.manager.PluginByAlias(alias)
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprint(writer, err)

		return
	}

	plugin.ReverseProxy().ServeHTTP(writer, reader)
}

//nolint:gochecknoglobals
var Handler APIHandler

const EndpointPattern = "/plugin/"

type endpointContent struct {
	pluginAuthor string
	pluginName   string
	subURI       string
}

func splitEndpoint(uri string) *endpointContent {
	// ["", "plugin", "{author-name}", "{plugin-name}", ...]
	exploded := strings.Split(uri, "/")

	return &endpointContent{
		pluginAuthor: FormatTextForURL(exploded[2]),
		pluginName:   FormatTextForURL(exploded[3]),
		subURI:       "/" + strings.Join(exploded[4:], "/"),
	}
}

// Interceptor intercepts requests for plugins.
// The endpoint is expected to have the following structure:
// /plugin/{author-name}/{plugin-name}/{plugin-endpoint}...
func Interceptor(req *interceptor.Interceptor) *interceptor.Interceptor {
	if req == nil {
		return nil
	}

	isMyMassa := strings.HasPrefix(req.Request.Host, config.MassaStationURL)
	indexPluginEndpoint := strings.Index(req.Request.RequestURI, EndpointPattern)

	if isMyMassa && indexPluginEndpoint > -1 {
		endpoint := splitEndpoint(req.Request.RequestURI)

		authorName, err := url.QueryUnescape(endpoint.pluginAuthor)
		if err != nil {
			config.Logger.Error(err.Error())

			return nil
		}

		pluginName, err := url.QueryUnescape(endpoint.pluginName)
		if err != nil {
			config.Logger.Error(err.Error())

			return nil
		}

		Handler.Handle(
			req.Writer, req.Request,
			authorName, pluginName,
		)

		return nil
	}

	return req
}

// modifyRequest rewrite the incoming request URL to match what the plugin is expecting to receive.
// All the `/plugin/{author-name}/{plugin-name}` template is removed.
func modifyRequest(req *http.Request) {
	urlExternal := req.URL.String()

	// the url has the following format:
	// 		http://127.0.0.1:1234/plugin/massalabs/hello-world/web/index.html?name=Massalabs
	// The idea is to rewrite url to remove: /plugin/massalabs/hello-world

	index := strings.Index(urlExternal, EndpointPattern)

	endpoint := splitEndpoint(urlExternal[index:])

	urlRewritten := urlExternal[:index] + endpoint.subURI

	urlNew, err := url.Parse(urlRewritten)
	if err != nil {
		panic(err)
	}

	req.URL = urlNew
}
