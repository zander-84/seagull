package contract

// Balancer yields endpoints according to some heuristic.
type Balancer interface {
	Update()
	Next() (any, error)
	All() ([]any, error)
	Used() map[any]int64
	Get(uid any) (any, error) //用于hash一致性
}
