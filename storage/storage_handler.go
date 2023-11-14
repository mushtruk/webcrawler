package storage

type StorageHandler interface {
	Handle(data interface{}) error
}
