package stores

// Repository is the interface which a store must implement
type Repository interface {
	Get(item interface{}) error
	GetSingle(id string, item interface{}) error
	Add(item interface{}) error
	Remove(id string) error
	Update(id string, item interface{}) error
}
