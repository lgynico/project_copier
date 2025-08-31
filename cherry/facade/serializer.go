package cherryFacade

type ISerializer interface {
	Marshal(any) ([]byte, error)
	Unmarshal([]byte, any) error
	Name() string
}
