package stores

// Repository is the interface which a store must implement
type Repository interface {
	Get(item interface{}, query interface{}, selector interface{}) error
	GetSingle(id string, item interface{}) error
	Add(item interface{}) (string, error)
	Remove(id string) error
	Update(id string, item interface{}) error
	IsValidID(id string) bool
}
