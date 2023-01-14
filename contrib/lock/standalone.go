package lock

import (
	"context"
	"github.com/zander-84/seagull/contract"
	"sync"
	"sync/atomic"
	"time"
)

type standaloneLocked struct {
	key    string
	id     string
	engine *standalone
}

func (s *standaloneLocked) Release(ctx context.Context) error {
	return s.engine.Release(ctx, s.key, s.id)
}

func (s *standaloneLocked) GetID() string {
	return s.id
}
func newStandaloneLocked(engine *standalone, key string, id string) contract.Locked {
	out := new(standaloneLocked)
	out.key = key
	out.id = id
	out.engine = engine
	return out
}

type standalone struct {
	sMap     sync.Map
	exit     chan struct{}
	unique   contract.Unique
	exitOnce sync.Once

	leaserInterval time.Duration
	leaser         map[string]chan struct{}
	leaserLocker   sync.RWMutex
	processor      contract.Processor
}

type singleLockData struct {
	expiredAt time.Time // 过期时间
	tag       int32
	val       interface{}
	locker    sync.Mutex
}

func NewStandaloneLock(unique contract.Unique, processor contract.Processor, leaserInterval time.Duration) (contract.Locker, func()) {
	s := new(standalone)
	s.unique = unique
	s.exit = make(chan struct{}, 1)
	go func() {
		for {
			select {
			case <-s.exit:
				return
			case <-time.After(time.Second * 60 * 30):
				s.clean()
			}
		}
	}()
	s.leaserInterval = leaserInterval
	s.leaser = make(map[string]chan struct{}, 0)
	s.processor = processor
	return s, s.cancel
}
func (s *standalone) Lock(ctx context.Context, key string, minExpiration time.Duration) (locked contract.Locked, err error) {
	id, ok, err := s.lock(ctx, key, minExpiration)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, contract.LockFailed
	}
	//续租
	if minExpiration > 0 {
		s.lease(ctx, key, id, minExpiration)
	}
	return newStandaloneLocked(s, key, id), nil
}
func (s *standalone) lock(ctx context.Context, key string, expiration time.Duration) (id string, ok bool, err error) {
	var expiredAt time.Time
	if expiration < 1 {
		expiredAt = time.Time{}
	} else {
		expiredAt = time.Now().Add(expiration)
	}
	identify := s.unique.ID()
	iData, ok := s.sMap.LoadOrStore(key, &singleLockData{
		expiredAt: expiredAt,
		tag:       0,
		val:       identify,
	})

	//不存在 第一次
	if !ok {
		return identify, true, nil
	}

	// 非第一次
	user := iData.(*singleLockData)
	//并发下 同一个人只有第一次人才可以进去判断时间，第二次以上直接返回错误
	if atomic.CompareAndSwapInt32(&user.tag, 0, 1) {
		user.locker.Lock()
		defer user.locker.Unlock()
		// 时间过期下才可以返回true
		if !user.expiredAt.IsZero() && time.Now().After(user.expiredAt) {
			user.expiredAt = expiredAt
			user.val = identify
			user.tag = 0
			s.sMap.Store(key, user)
			return identify, true, nil
		} else {
			user.tag = 0
			return identify, false, contract.LockFailed
		}
	} else {
		return identify, false, contract.LockFailed
	}
}

func (s *standalone) Release(ctx context.Context, key string, identify string) error {
	s.exitAndDelLeaser(key, identify)
	valInterface, ok := s.sMap.Load(key)
	if !ok {
		return nil
	}
	lockData, ok := valInterface.(*singleLockData)
	if ok && lockData.val == identify {
		s.sMap.Delete(key)
	}
	return nil
}

func (s *standalone) clean() {
	s.sMap.Range(func(key, value interface{}) bool {
		user, ok := value.(*singleLockData)
		if !ok {
			return true
		}
		// 时间0  永不过期
		if user.expiredAt.IsZero() {
			return true
		}
		if time.Now().After(user.expiredAt) {
			s.sMap.Delete(key)
		}
		return true
	})
}

func (s *standalone) cancel() {
	s.exitOnce.Do(func() {
		close(s.exit)
	})
}

func (s *standalone) key(key string, id string) string {
	return key + ":" + id
}

func (s *standalone) addLeaser(key string, id string) <-chan struct{} {
	s.leaserLocker.Lock()
	defer s.leaserLocker.Unlock()
	realKey := s.key(key, id)
	if _, ok := s.leaser[realKey]; !ok {
		s.leaser[realKey] = make(chan struct{}, 0)
	}

	return s.leaser[realKey]
}
func (s *standalone) exitAndDelLeaser(key string, id string) {
	s.leaserLocker.Lock()
	defer s.leaserLocker.Unlock()

	realKey := s.key(key, id)
	leaser, ok := s.leaser[realKey]
	if ok {
		close(leaser)
		delete(s.leaser, realKey)
	}
}

func (s *standalone) lease(ctx context.Context, key string, id string, expiration time.Duration) {
	leaserChan := s.addLeaser(key, id)
	realKey := s.key(key, id)

	s.processor.Go(func() {
		for {
			select {
			case <-leaserChan:
				return
			case <-time.After(s.leaserInterval):
				// 续租
				var expiredAt time.Time
				if expiration < 1 {
					expiredAt = time.Time{}
				} else {
					expiredAt = time.Now().Add(expiration)
				}
				iData, ok := s.sMap.Load(realKey)
				if !ok {
					return
				}
				user := iData.(*singleLockData)
				user.locker.Lock()
				if user.val != id {
					user.locker.Unlock()
					return
				}
				user.expiredAt = expiredAt
				user.locker.Unlock()
			}
		}
	}, nil)
}
