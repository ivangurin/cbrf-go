package service_provider

import (
	cbrf_service "github.com/ivangurin/cbrf-go/internal/service/cbrf"
)

type services struct {
	cbrfService cbrf_service.IService
}

func (sp *ServiceProvider) GetCbrfService() cbrf_service.IService {
	if sp.services.cbrfService == nil {
		sp.services.cbrfService = cbrf_service.NewService(
			sp.GetCbrfClient(),
		)
	}
	return sp.services.cbrfService
}
