// Package gqlclient creates a client to connect with the GQL Server GraphQL API
package gqlclient

import (
	"net/http"

	"github.com/kiwisheets/gql-server/client"
	"github.com/kiwisheets/util"
	"github.com/sethgrid/pester"
)

// NewClient creates a new GQL Server client
func NewClient(cfg *util.ClientConfig) *client.Client {
	return client.NewClient(pester.New(), cfg.BaseURL, func(req *http.Request) {
		if cfg.CfClientID != "" && cfg.CfClientSecret != "" {
			req.Header.Set("CF-Access-Client-Id", cfg.CfClientID)
			req.Header.Set("CF-Access-Client-Secret", cfg.CfClientSecret)
		}
	})
}
