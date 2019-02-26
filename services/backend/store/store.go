package store

type Storer interface {
	Get(id string) ([]byte, error)
	Put(id string, data []byte) error
}
