package contract

import (
	"go.mongodb.org/mongo-driver/bson"
	"sync"
)

type SearchMeta interface {
	UseCursor(cursor bool) SearchMeta
	IsCursor() bool

	UsePage(page bool) SearchMeta
	IsPage() bool

	SetPage(page int) SearchMeta
	SetPageSize(pageSize int) SearchMeta
	SetMaxPage(page int) SearchMeta
	SetMaxPageSize(pageSize int) SearchMeta

	Page() int
	PageSize() int
	Offset() int

	IsCount() bool
	UseCount(cnt bool) SearchMeta
}

type MysqlBuilder interface {
	AppendWhere(query string, args ...any) MysqlBuilder
	BuildQuery() string
	Args() []any
	SetFields(fs string) MysqlBuilder
	Fields() string
	SetOrderBy(orderBy string) MysqlBuilder
	OrderBy() string
	SetTag(tag string) MysqlBuilder
	Tag() string
}

type Mysql interface {
	Create(entity IEntity) error
	CreateTx(tx any, entity IEntity) error

	FindByID(id int64, entity IEntity) error
	FindOneByField(field string, val any, entity IEntity) error
	FindIn(field string, val any, ptrSliceData interface{}) error
	Exist(field string, val any) (bool, error)

	UpdateMap(id int64, version int, data map[string]interface{}) error
	UpdateMapTx(tx any, id int64, version int, data map[string]interface{}) error
	Update(id int64, version int, entity IEntity) error
	UpdateTx(tx any, id int64, version int, entity IEntity) error

	Search(searchMeta SearchMeta, searchParams MysqlBuilder, ptrSliceData interface{}, cnt *int64) (err error)
}

type MongoBuilder interface {
	AppendWhere(be bson.E) MongoBuilder
	BuildQuery() bson.D
	SetFields(fs []string) MongoBuilder
	Fields() bson.D
	SetOrderBy(orderBy bson.D) MongoBuilder
	OrderBy() bson.D
	SetTag(tag string) MongoBuilder
	Tag() string
}

// Mongo 一维文档操作
type Mongo interface {
	Create(entity IEntity) (string, error)
	FindByID(id string, entity IEntity) error
	FindOneByField(field string, val any, entity IEntity) error
	FindIn(field string, val any, ptrSliceData interface{}) error
	Exist(field string, val any) (bool, error)

	// ReplaceOne version没实现原子递增
	ReplaceOne(id string, version int, entity IEntity) error
	ReplaceOneByKv(field string, val any, version int, entity IEntity) error

	Update(id string, version int, entity IEntity) error
	UpdateByKv(field string, val any, version int, entity IEntity) error

	DelOneByKv(key string, val any) error
	DelByKv(key string, val any) error
	DelByCondition(builder MongoBuilder) error
	Search(searchMeta SearchMeta, searchParams MongoBuilder, ptrSliceData interface{}, cnt *int64) (err error)
}

type IEntity interface {
	UpdatedFields() map[string]any
	UpdateUpdatedAt(updatedAt int64)
	UpdateCreatedAt(createdAt int64)
	UpdateVersion(version int)
	GetVersion() int
}

type UpdateFields struct {
	data sync.Map `json:"-" bson:"-"`
}

func (u *UpdateFields) Update(key string, val any) {
	u.data.Store(key, val)
}

func (u *UpdateFields) Reset() {
	u.data = sync.Map{}
}

func (u *UpdateFields) Get() map[string]any {
	out := map[string]any{}
	u.data.Range(func(key, value any) bool {
		keyStr, ok := key.(string)
		if ok {
			out[keyStr] = value
		}
		return true
	})
	return out
}
