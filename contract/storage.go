package contract

import "go.mongodb.org/mongo-driver/bson"

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

	Update(id int64, version int64, data map[string]interface{}) error
	UpdateTx(tx any, id int64, version int64, data map[string]interface{}) error
	UpdatePart(id int64, version int64, entity IEntity) error
	UpdateTxPart(tx any, id int64, version int64, entity IEntity) error

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

	// Update 全量更新 必须带版本号
	Update(id string, version int64, entity IEntity) error

	UpdatePart(id string, version int64, entity IEntity) error

	Search(searchMeta SearchMeta, searchParams MongoBuilder, ptrSliceData interface{}, cnt *int64) (err error)
}

type IEntity interface {
	UpdatedFields() map[string]any
	SetUpdatedAt(updatedAt int64)
	SetCreatedAt(createdAt int64)
	SetVersion(version int64)
}
