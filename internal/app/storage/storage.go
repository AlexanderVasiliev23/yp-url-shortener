package storage

type Storage interface {
	Add(token, url string) error
	Get(token string) (string, error)
}
