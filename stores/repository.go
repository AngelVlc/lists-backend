package stores

// Repository is the interface which a store must implement
type Repository interface {
	Get(item interface{}, query interface{}, selector interface{}) error
	GetOne(item interface{}, query interface{}, selector interface{}) error
	Add(item interface{}) (string, error)
	Remove(query interface{}) error
	Update(query interface{}, item interface{}) error
	IsValidID(id string) bool
}
