package cbrf_client

//go:generate mockery --name=(.+)Mock --case=underscore  --with-expecter

type ClientMock interface {
	IClient
}
