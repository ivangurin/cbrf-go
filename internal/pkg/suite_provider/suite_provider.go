package suite_provider

import (
	cbrf_client "github.com/ivangurin/cbrf-go/internal/pkg/client/cbrf"
	cbrf_client_mocks "github.com/ivangurin/cbrf-go/internal/pkg/client/cbrf/mocks"
	cbrf_service "github.com/ivangurin/cbrf-go/internal/service/cbrf"
)

type suiteProvider struct {
	cbrfClient     cbrf_client.IClient
	cbrfClientMock *cbrf_client_mocks.ClientMock
	cbrfService    cbrf_service.IService
}

func NewSuiteProvider() *suiteProvider {
	return &suiteProvider{}
}

func (sp *suiteProvider) GetCbrfClientMock() *cbrf_client_mocks.ClientMock {
	if sp.cbrfClientMock == nil {
		sp.cbrfClientMock = &cbrf_client_mocks.ClientMock{}
	}
	return sp.cbrfClientMock
}

func (sp *suiteProvider) GetCbrfClient() cbrf_client.IClient {
	if sp.cbrfClient == nil {
		sp.cbrfClient = sp.GetCbrfClientMock()
	}
	return sp.cbrfClient
}

func (sp *suiteProvider) GetCbrfService() cbrf_service.IService {
	if sp.cbrfService == nil {
		sp.cbrfService = cbrf_service.NewService(
			sp.GetCbrfClient(),
		)
	}
	return sp.cbrfService
}
