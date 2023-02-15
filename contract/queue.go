package contract

import (
	"context"
	"errors"
)

/**
队列3种模型：
1. 通用模型：失败存持久化在queue中，但是不好观察数据
2. 优雅模型：错误时候持久化，调整完毕在继续消费
3. 备份模型：存储原始数据，队列数据一份到队列 一份到mongo

*/

type QCode int32
type QDataComeFrom int32

const (
	QCodeOrigin      QCode = 0 //原始数据
	QCodeWaiting     QCode = 1 //等待消费
	QCodeConsumeFail QCode = 2 //数据消费失败
	QCodeDel         QCode = 3 //软删除
	QCodeInQFail     QCode = 4 //入q失败

	QDataFromFirstTime QDataComeFrom = 0
	QDataFromBackup    QDataComeFrom = 1
)

type QProducer interface {
	// SendOrigin 发送原始数据
	SendOrigin(message *QMessage) error
	// Send 发送数据
	Send(message *QMessage) error

	// SendFromBackup 从备份数据发送
	SendFromBackup(message *QMessage) error

	Close(ctx context.Context) error
}

type QConsumer interface {
	// Consume 消费
	Consume(workers map[string]func(data string) error) error

	Close(ctx context.Context) error
}

type QManager interface {

	// Messages 消息列表
	Messages(searchMeta SearchMeta, searchParams MongoBuilder) (message []QMessage, cnt *int64, err error)

	// AdjustMessage 调整消息
	AdjustMessage(message *QMessage) error
	// SendFromBackup 从备份数据发送
	SendFromBackup(Id string) error

	// ReleaseMessage 释放数据
	ReleaseMessage(id string) error

	Close(ctx context.Context) error
}

//--------------------------------------

type QMess struct {
	Kind string
	UID  string
	Body []byte
}

// QMessage 采用json是因为数据可能要入mongo
// 只要写入 Kind  UID Data 即可，其余都自动生成
type QMessage struct {
	Topic      string        `json:"topic,omitempty" bson:"topic"`             // 主题
	Kind       string        `json:"kind,omitempty" bson:"kind"`               // 分类
	UID        string        `json:"uid,omitempty" bson:"uid"`                 // 消息ID 唯一
	ForeignKey string        `json:"foreign_key,omitempty" bson:"foreign_key"` // 业务ID 用于搜索
	Code       QCode         `json:"code,omitempty" bson:"code"`               // 消息状态码
	ComeFrom   QDataComeFrom `json:"come_from,omitempty" bson:"come_from"`     // 来源
	Reason     string        `json:"reason,omitempty" bson:"reason"`           // 备注：失败原因
	Data       string        `json:"data,omitempty" bson:"data"`               // 数据
	Version    int           `json:"version,omitempty" bson:"version"`         // 数据版本号
	CreatedAt  int64         `json:"created_at,omitempty" bson:"created_at"`   // 毫秒
	UpdatedAt  int64         `json:"updated_at,omitempty" bson:"updated_at"`   // 毫秒

	updateFields UpdateFields `json:"-" bson:"-"`
}

func (m *QMessage) NewQMessage() *QMessage {
	out := new(QMessage)
	return out
}

func (m *QMessage) Valid() error {
	if m.UID == "" {
		return errors.New("UID can not empty")
	}
	if m.Data == "" {
		return errors.New("data can not empty")
	}
	return nil
}

func (m *QMessage) UpdatedFields() map[string]any {
	return m.updateFields.Get()
}

func (m *QMessage) update(key string, val any) {
	m.updateFields.Update(key, val)
}

func (m *QMessage) SetTopic(topic string) {
	m.update("topic", topic)
	m.Topic = topic
}

func (m *QMessage) SetKind(kind string) {
	m.update("kind", kind)
	m.Kind = kind
}

func (m *QMessage) SetUID(UID string) {
	m.update("uid", UID)
	m.UID = UID
}

func (m *QMessage) SetForeignKey(ForeignKey string) {
	m.update("foreign_key", ForeignKey)
	m.ForeignKey = ForeignKey
}

func (m *QMessage) SetCode(c QCode) {
	m.update("code", c)
	m.Code = c
}

func (m *QMessage) SetData(data string) {
	m.update("data", data)
	m.Data = data
}
func (m *QMessage) SetComeFrom(data QDataComeFrom) {
	m.update("come_from", data)
	m.ComeFrom = data
}
func (m *QMessage) SetReason(reason string) {
	m.update("reason", reason)
	m.Reason = reason
}

func (m *QMessage) UpdateVersion(version int) {
	m.Version = version
}
func (m *QMessage) GetVersion() int {
	return m.Version
}

func (m *QMessage) UpdateUpdatedAt(updatedAt int64) {
	m.update("updated_at", updatedAt)
	m.UpdatedAt = updatedAt
}

func (m *QMessage) UpdateCreatedAt(createdAt int64) {
	m.update("created_at", createdAt)
	m.CreatedAt = createdAt
}
