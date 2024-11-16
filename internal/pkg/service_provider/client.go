package service_provider

import (
	cbrf_client "github.com/ivangurin/cbrf-go/internal/pkg/client/cbrf"
)

type clients struct {
	cbrfClient cbrf_client.IClient
}

func (sp *ServiceProvider) GetCbrfClient() cbrf_client.IClient {
	if sp.clients.cbrfClient == nil {
		sp.clients.cbrfClient = cbrf_client.NewClient()
	}
	return sp.clients.cbrfClient
}
