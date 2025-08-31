package cherryFacade

type INetParser interface {
	Load(app IApplication)
	AddConnector(connector IConnector)
	Connectors() []IConnector
}
