package kvstore

type Store interface {
	Save(key string, value interface{}) error
	Get(key string, receiver interface{}) error
	Delete(key string) error
	Close() error
}