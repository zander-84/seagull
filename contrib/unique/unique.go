package unique

import (
	rand2 "crypto/rand"
	"errors"
	"fmt"
	"github.com/zander-84/seagull/contract"
	"math/big"
	"math/rand"
	"strings"
	"sync"
	"time"
)

type unique struct {
	incrementID     uint64 // 自增ID
	incrementTimeID uint64 // 自增ID
	joinSymbol      string // 连接符
	lock            sync.Mutex
	lastTime        int64  // 上次更新时间
	serverID        string // 上次更新时间
	salt            string
	checkCode       contract.CheckCode
}

func (u *unique) ID() string {
	return u.id()
}

func (u *unique) Check(id string) error {
	id = strings.TrimPrefix(id, u.serverID)

	if len(id) < 3 {
		return errors.New("err id")
	}
	checkCode := string([]rune(id)[0:3])
	id = strings.TrimPrefix(id, checkCode)

	return u.checkCode.Check(id+u.salt, checkCode)
}

func (u *unique) realRandomInt64(min int64, max int64) int64 {
	if result, err := rand2.Int(rand2.Reader, big.NewInt(max-min)); err == nil {
		data := result.Int64()
		return data + min
	} else {
		res := rand.Intn(int(max - min))
		return int64(res) + min
	}
}

func (u *unique) now() time.Time {
	return time.Now()
}

func (u *unique) id() string {
	u.lock.Lock()
	currentTime := u.now()
	currentTimeUnix := currentTime.Unix()
	if u.lastTime == 0 {
		u.lastTime = currentTimeUnix
	}

	u.incrementID += 1
	if u.incrementID > 999999 {
		u.incrementID = 1

		if u.lastTime >= currentTimeUnix { // 同秒内超过99w次 重置自增位，时间加1，否则自增加1
			time.Sleep(1001 * time.Millisecond) //防止一秒破百万造成全局不唯一
			currentTime = u.now()
			currentTimeUnix = currentTime.Unix()
			u.lastTime = currentTimeUnix
		}
	}

	if u.lastTime > currentTimeUnix { // 时间回滚下 修改时间标志位 重置时间
		u.incrementTimeID = u.incrementTimeID + 1
		if u.incrementTimeID > 99 {
			u.incrementTimeID = 10
		}
		currentTime = u.now()
		currentTimeUnix = currentTime.Unix()
		u.lastTime = currentTimeUnix
	}
	u.lastTime = currentTimeUnix
	d := u.incrementID
	dt := u.incrementTimeID
	u.lock.Unlock()

	code := fmt.Sprintf("%d%d", dt, rand.Intn(10)) + fmt.Sprintf("%06d", d) + u.joinSymbol + currentTime.Format("060102150405")
	checkCode := u.checkCode.Sign(code + u.salt)

	return u.serverID + checkCode + code
}

func New(serverID string, joinSymbol string, salt string, checkCode contract.CheckCode) contract.Unique {
	rand.Seed(time.Now().UnixNano())

	u := new(unique)
	u.serverID = serverID
	u.joinSymbol = joinSymbol
	u.incrementTimeID = uint64(u.realRandomInt64(10, 100))
	u.checkCode = checkCode
	u.salt = salt
	return u
}
