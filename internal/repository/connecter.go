package repository

type Connecter interface {
	Connect() (interface{}, error)
	Close() error
}
