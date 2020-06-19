package kvstore

type BadgerStore struct {

}

func (s BadgerStore) Save(key string, value interface{}) error {
	return nil
}

func (s BadgerStore) Get(key string, receiver interface{}) error {
	return nil
}

func (s BadgerStore) Delete(key string) error {
	return nil
}

func Close() error {
	return nil
}