package contract

type Unique interface {
	ID() string
	Check(id string) error
}
