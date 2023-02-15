package lock

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/zander-84/seagull/contract"
	"sync"
	"time"
)

type redisV2Lock struct {
	engine        *redis.Client
	unique        contract.Unique
	listenChannel string
	waitTime      time.Duration
	waiter        map[string]chan struct{}
	waiterCnt     map[string]int
	waiterLocker  sync.RWMutex
	pubSub        *redis.PubSub
	redisLeaser   *redisLeaser
}
type redisV2Locked struct {
	key    string
	id     string
	engine *redisV2Lock
}

func (r *redisV2Locked) Release(ctx context.Context) error {
	return r.engine.Release(ctx, r.key, r.id)
}

func (r *redisV2Locked) GetID() string {
	return r.id
}

func newRedisV2Locked(engine *redisV2Lock, key string, id string) contract.Locked {
	out := new(redisV2Locked)
	out.key = key
	out.id = id
	out.engine = engine
	return out
}

// NewRedisWaitLocker  一把基于redis消息订阅的等待锁
func NewRedisWaitLocker(engine *redis.Client, unique contract.Unique, processor contract.Processor, listenChannel string, waitTime time.Duration) (locker contract.Locker, cancel func(), err error) {

	out := &redisV2Lock{
		engine:        engine,
		unique:        unique,
		listenChannel: listenChannel,
		waitTime:      waitTime,

		waiter:    make(map[string]chan struct{}, 0),
		waiterCnt: make(map[string]int, 0),
	}
	out.pubSub = out.engine.Subscribe(context.Background(), out.listenChannel)

	go out.subscribe()
	out.redisLeaser, err = newRedisLeaser(engine, processor)
	//go func() {
	//	for {
	//		time.Sleep(time.Second / 2)
	//		fmt.Println("test:", out.getLeasers())
	//	}
	//}()
	return out, out.cancel, err
}

func (r *redisV2Lock) cancel() {
	r.pubSub.Close()
}
func (r *redisV2Lock) subscribe() {
	for {
		message, err := r.pubSub.ReceiveMessage(context.Background())
		if err != nil {
			if err == redis.ErrClosed {
				return
			}
			time.Sleep(time.Second)
		} else {
			waiter, ok := r.getWaiter(message.Payload)
			if !ok {
				continue
			}
			waiter <- struct{}{}
		}
	}
}
func (r *redisV2Lock) Lock(ctx context.Context, key string, minExpiration time.Duration) (contract.Locked, error) {
	id, ok, err := r.lock(ctx, key, minExpiration)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, contract.LockFailed
	}
	if err == nil && ok {
		// 续租  Release 时候释放
		if minExpiration > 0 {
			r.redisLeaser.lease(ctx, key, id, minExpiration)
		}
	}

	return newRedisV2Locked(r, key, id), err
}
func (r *redisV2Lock) lock(ctx context.Context, key string, expiration time.Duration) (string, bool, error) {
	identify := r.unique.ID()
	ok, err := r.engine.SetNX(ctx, key, identify, expiration).Result()
	if err != nil {
		return "", false, err
	}
	// 成功
	if ok {
		return identify, ok, err
	}

	waiter := r.addWaiter(key)
	defer func() {
		r.delWaiter(key)
	}()
	ctxDone := ctx.Done()
	waitTime := time.NewTimer(r.waitTime)

	if ctxDone != nil {
		for {
			select {
			case <-waiter:
				ok, err := r.engine.SetNX(ctx, key, identify, expiration).Result()
				if err != nil {
					return "", false, err
				}
				if ok {
					return identify, ok, nil
				} else {
					continue
				}
			case <-time.After(time.Second * 3):
				ok, err := r.engine.SetNX(ctx, key, identify, expiration).Result()
				if err != nil {
					return "", false, err
				}
				if ok {
					return identify, ok, nil
				} else {
					continue
				}
			case <-waitTime.C:
				return "", false, nil
			case <-ctxDone:
				return "", false, ctx.Err()
			}
		}
	} else {
		for {
			select {
			case <-waiter:
				ok, err := r.engine.SetNX(ctx, key, identify, expiration).Result()
				if err != nil {
					return "", false, err
				}
				if ok {
					return identify, ok, nil
				} else {
					continue
				}
			case <-time.After(time.Second * 3):
				ok, err := r.engine.SetNX(ctx, key, identify, expiration).Result()
				if err != nil {
					return "", false, err
				}
				if ok {
					return identify, ok, nil
				} else {
					continue
				}
			case <-waitTime.C:
				return "", false, nil
			}
		}
	}
}

func (r *redisV2Lock) Release(ctx context.Context, key string, identify string) error {
	defer r.engine.Publish(ctx, r.listenChannel, key)
	_ = r.redisLeaser.release(ctx, key, identify)
	data, err := r.engine.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		return err
	}

	if identify == "" || data == identify {
		if err := r.engine.Del(ctx, key).Err(); err != nil {
			if err == redis.Nil {
				return nil
			}
			err = r.engine.Del(ctx, key).Err()
			if err == redis.Nil {
				return nil
			}
			return err
		}
	}
	return err
}

func (r *redisV2Lock) addWaiter(key string) <-chan struct{} {
	r.waiterLocker.Lock()
	defer r.waiterLocker.Unlock()
	if _, ok := r.waiter[key]; !ok {
		r.waiter[key] = make(chan struct{}, 0)
	}
	if _, ok := r.waiterCnt[key]; !ok {
		r.waiterCnt[key] = 1
	} else {
		r.waiterCnt[key] += 1
	}
	return r.waiter[key]
}

func (r *redisV2Lock) delWaiter(key string) {
	r.waiterLocker.Lock()
	defer r.waiterLocker.Unlock()
	if _, ok := r.waiterCnt[key]; ok {
		r.waiterCnt[key] -= 1
		if r.waiterCnt[key] < 1 {
			delete(r.waiterCnt, key)
			delete(r.waiter, key)
		}
	}
}

func (r *redisV2Lock) getWaiter(key string) (chan struct{}, bool) {
	r.waiterLocker.RLock()
	defer r.waiterLocker.RUnlock()
	out, ok := r.waiter[key]
	return out, ok
}
