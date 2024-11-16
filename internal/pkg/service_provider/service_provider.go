package service_provider

type ServiceProvider struct {
	clients  clients
	services services
}

var serviceProvider *ServiceProvider

func GetServiceProvider() *ServiceProvider {
	if serviceProvider == nil {
		serviceProvider = &ServiceProvider{}
	}
	return serviceProvider
}
