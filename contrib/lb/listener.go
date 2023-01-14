package lb

// Listener 动态数据监听
type Listener interface {
	SetErr(err error)
	Err() error
	Get() (map[any]int, []any, uint64)
	Version() uint64
	Remove(data any) error
	AddWithWeight(data any, weight int) error
	Add(data any) error
	Set(data map[any]int) error
	Exist(data any) bool
	Close()
}
