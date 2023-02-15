package tool

import "sync"

type ConcurrentMap struct {
	lock sync.RWMutex
	val  map[any]any
}

func NewConcurrentMap() *ConcurrentMap {
	data := new(ConcurrentMap)
	data.val = make(map[any]any, 0)
	return data
}

func (cm *ConcurrentMap) Set(key any, val any) {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	cm.val[key] = val
}
func (cm *ConcurrentMap) Get(key any) (any, bool) {
	cm.lock.RLock()
	defer cm.lock.RUnlock()
	data, ok := cm.val[key]
	return data, ok
}

func (cm *ConcurrentMap) Exist(key any) bool {
	cm.lock.RLock()
	defer cm.lock.RUnlock()
	_, ok := cm.val[key]
	return ok
}

func (cm *ConcurrentMap) GetMap() map[any]any {
	cm.lock.RLock()
	defer cm.lock.RUnlock()
	res := make(map[any]interface{}, 0)
	for k, v := range cm.val {
		res[k] = v
	}
	return res
}

func (cm *ConcurrentMap) ShouldGetString(key any) string {
	data, ok := cm.Get(key)
	if ok {
		res, _ := data.(string)
		return res
	}
	return ""
}

func (cm *ConcurrentMap) ShouldGetInt32(key any) int32 {
	data, ok := cm.Get(key)
	if ok {
		res, _ := data.(int32)
		return res
	}
	return 0
}

func (cm *ConcurrentMap) ShouldGetInt64(key any) int64 {
	data, ok := cm.Get(key)
	if ok {
		res, _ := data.(int64)
		return res
	}
	return 0
}

func (cm *ConcurrentMap) ShouldGetInt(key any) int {
	data, ok := cm.Get(key)
	if ok {
		res, _ := data.(int)
		return res
	}
	return 0
}

func (cm *ConcurrentMap) ShouldGetFloat64(key any) float64 {
	data, ok := cm.Get(key)
	if ok {
		res, _ := data.(float64)
		return res
	}
	return 0
}

func (cm *ConcurrentMap) ShouldGetFloat32(key any) float32 {
	data, ok := cm.Get(key)
	if ok {
		res, _ := data.(float32)
		return res
	}
	return 0
}

func (cm *ConcurrentMap) ShouldGetBool(key any) bool {
	data, ok := cm.Get(key)
	if ok {
		res, _ := data.(bool)
		return res
	}
	return false
}
