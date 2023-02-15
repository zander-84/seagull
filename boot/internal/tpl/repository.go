package tpl

var RepositoryMysqlCache = `package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/groupcache/singleflight"
	"github.com/zander-84/seagull/contract"
	"github.com/zander-84/seagull/contract/def"
	"github.com/zander-84/seagull/contrib/cache/wrapcache"
	"github.com/zander-84/seagull/contrib/storage"
	"github.com/zander-84/seagull/drive/gormmysql"
	"${project}/apps/${server}/internal/entity"
	"time"
)

type ${entityName}Mysql struct {
	db                    *gormmysql.Gdb
	dbHelper              contract.Mysql
	singleFlight          singleflight.Group
	cache                 contract.Cache
	_cachePrefix          string
	cacheGetOrSetDuration time.Duration
	cacheSetDuration      time.Duration
	version               int
}

func new${EntityName}Mysql(db *gormmysql.Gdb, cache contract.Cache, cachePrefix string) *${entityName}Mysql {
	out := &${entityName}Mysql{
		db:                    db,
		singleFlight:          singleflight.Group{},
		cache:                 cache,
		_cachePrefix:          cachePrefix,
		cacheGetOrSetDuration: ${cacheGetOrSetDuration},
		cacheSetDuration:      ${cacheSetDuration},
		version:              ${version},
	}
	out.dbHelper = storage.NewGormMysql(db.Engine(), out.Name(), 50)
	return out
}

func (${shortEntityName} *${entityName}Mysql) Name() string {
	return "${tableName}"
}

func (${shortEntityName} *${entityName}Mysql) Version() int {
	return ${shortEntityName}.version
}

func (${shortEntityName} *${entityName}Mysql) CachePrefix() string {
	return ${shortEntityName}._cachePrefix
}
func (${shortEntityName} *${entityName}Mysql) assert(in interface{}) (*entity.${EntityName}, error) {
	if _${entityName}, ok := in.(*entity.${EntityName}); !ok || _${entityName} == nil {
		return nil, errors.New("err type")
	} else {
		return _${entityName}, nil
	}
}

func (${shortEntityName} *${entityName}Mysql) Create(${entityName} *entity.${EntityName}) error {
	return ${shortEntityName}.CreateTx(${shortEntityName}.db.Engine(), ${entityName})
}

func (${shortEntityName} *${entityName}Mysql) CreateTx(db interface{}, ${entityName} *entity.${EntityName}) error {
	return ${shortEntityName}.dbHelper.CreateTx(db, ${entityName})
}

func (${shortEntityName} *${entityName}Mysql) GetRaw(id int64) (*entity.${EntityName}, error) {
	out, err := ${shortEntityName}.singleFlight.Do(fmt.Sprintf("GetRaw:%d", id), func() (interface{}, error) {
		out := new(entity.${EntityName})
		err := ${shortEntityName}.dbHelper.FindByID(id, out)
		return out, err
	})
	if err != nil {
		return nil, err
	}
	return ${shortEntityName}.assert(out)
}

func (${shortEntityName} *${entityName}Mysql) Get(id int64) (*entity.${EntityName}, error) {
	out := new(entity.${EntityName})
	err := ${shortEntityName}.cache.GetOrSet(context.Background(), wrapcache.RepoKey(${shortEntityName} , id), out, ${shortEntityName}.cacheGetOrSetDuration, func(key def.K) (value any, err error) {
		return ${shortEntityName}.GetRaw(id)
	})
	return out, err
}
func (${shortEntityName} *${entityName}Mysql) BatchGetRaw(ids []int64) ([]entity.${EntityName}, error) {
	out := make([]entity.${EntityName}, 0, len(ids))
	err := ${shortEntityName}.dbHelper.FindIn("id", ids, &out)
	return out, err
}

func (${shortEntityName} *${entityName}Mysql) BatchGet(ids []int64) ([]entity.${EntityName}, error) {
	out := make([]entity.${EntityName}, 0, len(ids))
	err := ${shortEntityName}.cache.BatchGetOrSet(context.Background(), wrapcache.RepoKeys(${shortEntityName} , ids), &out, ${shortEntityName}.cacheSetDuration, func(missIds []def.K) (value map[string]any, err error) {
		missDbIds, err := wrapcache.RepoDBKeys(missIds)
		if err != nil {
			return nil, err
		}
		data := make([]entity.${EntityName}, 0, len(missDbIds))
		if err := ${shortEntityName}.dbHelper.FindIn("id", missDbIds, &data); err != nil {
			return nil, err
		}

		value = make(map[string]any)
		for _, v := range data {
			for _, vv := range missIds {
				if v.Id == vv.Alias[0] {
					value[vv.Key] = v
				}
			}
		}
		return value, nil
	})
	return out, err
}

func (${shortEntityName} *${entityName}Mysql) Exist(key string, val any) (bool, error) {
	return ${shortEntityName}.dbHelper.Exist(key, val)
}

func (${shortEntityName} *${entityName}Mysql) UpdateMap(id int64, version int, data map[string]interface{}) error {
	err := ${shortEntityName}.dbHelper.UpdateMap(id, version, data)
	if err != nil {
		return err
	}

	rawData, err := ${shortEntityName}.GetRaw(id)
	if err != nil {
		return err
	}
	err = ${shortEntityName}.cache.Set(context.Background(), wrapcache.RepoKey(${shortEntityName} , id), rawData, ${shortEntityName}.cacheSetDuration)
	if err != nil {
		return err
	}
	return nil
}

func (${shortEntityName} *${entityName}Mysql) UpdateMapTx(tx any, id int64, version int, data map[string]interface{}) (func() error, error) {
	err := ${shortEntityName}.dbHelper.UpdateMapTx(tx, id, version, data)
	if err != nil {
		return nil, err
	}
	return func() error {
		rawData, err := ${shortEntityName}.GetRaw(id)
		if err != nil {
			return err
		}
		return ${shortEntityName}.cache.Set(context.Background(), wrapcache.RepoKey(${shortEntityName} , id), rawData, ${shortEntityName}.cacheSetDuration)
	}, nil
}

func (${shortEntityName} *${entityName}Mysql) Update(id int64, version int, ${entityName} *entity.${EntityName}) error {
	err := ${shortEntityName}.dbHelper.Update(id, version, ${entityName})
	if err != nil {
		return err
	}
	err = ${shortEntityName}.cache.Set(context.Background(), wrapcache.RepoKey(${shortEntityName} , id), ${entityName}, ${shortEntityName}.cacheSetDuration)
	if err != nil {
		return err
	}
	return nil
}

func (${shortEntityName} *${entityName}Mysql) UpdateTx(tx any, id int64, version int, ${entityName} *entity.${EntityName}) (func() error, error) {
	err := ${shortEntityName}.dbHelper.UpdateTx(tx, id, version, ${entityName})
	if err != nil {
		return nil, err
	}
	return func() error {
		return ${shortEntityName}.cache.Set(context.Background(), wrapcache.RepoKey(${shortEntityName} , id), ${entityName}, ${shortEntityName}.cacheSetDuration)
	}, nil
}

func (${shortEntityName} *${entityName}Mysql) Search(searchMeta contract.SearchMeta, searchParams contract.MysqlBuilder) (data []entity.${EntityName}, cnt int64, err error) {
	err = ${shortEntityName}.dbHelper.Search(searchMeta, searchParams, &data, &cnt)
	return data, cnt, err
}
`

var RepositoryMysql = `package repository

import (
	"errors"
	"fmt"
	"github.com/golang/groupcache/singleflight"
	"github.com/zander-84/seagull/contract"
	"github.com/zander-84/seagull/contrib/storage"
	"github.com/zander-84/seagull/drive/gormmysql"
	"${project}/apps/${server}/internal/entity"
)

type ${entityName}Mysql struct {
	db                    *gormmysql.Gdb
	dbHelper              contract.Mysql
	singleFlight          singleflight.Group
	version               int
}

func new${EntityName}Mysql(db *gormmysql.Gdb) *${entityName}Mysql {
	out := &${entityName}Mysql{
		db:                    db,
		singleFlight:          singleflight.Group{},
		version:               1,
	}
	out.dbHelper = storage.NewGormMysql(db.Engine(), out.Name(), 50)
	return out
}

func (${shortEntityName} *${entityName}Mysql) Name() string {
	return "${tableName}"
}
func (${shortEntityName} *${entityName}Mysql) Version() int {
	return ${shortEntityName}.version
}
func (${shortEntityName} *${entityName}Mysql) assert(in interface{}) (*entity.${EntityName}, error) {
	if _${entityName}, ok := in.(*entity.${EntityName}); !ok || _${entityName} == nil {
		return nil, errors.New("err type")
	} else {
		return _${entityName}, nil
	}
}

func (${shortEntityName} *${entityName}Mysql) Create(${entityName} *entity.${EntityName}) error {
	return ${shortEntityName}.CreateTx(${shortEntityName}.db.Engine(), ${entityName})
}

func (${shortEntityName} *${entityName}Mysql) CreateTx(db interface{}, ${entityName} *entity.${EntityName}) error {
	return ${shortEntityName}.dbHelper.CreateTx(db, ${entityName})
}

func (${shortEntityName} *${entityName}Mysql) GetRaw(id int64) (*entity.${EntityName}, error) {
	out, err := ${shortEntityName}.singleFlight.Do(fmt.Sprintf("GetRaw:%d", id), func() (interface{}, error) {
		out := new(entity.${EntityName})
		err := ${shortEntityName}.dbHelper.FindByID(id, out)
		return out, err
	})
	if err != nil {
		return nil, err
	}
	return ${shortEntityName}.assert(out)
}

func (${shortEntityName} *${entityName}Mysql) Get(id int64) (*entity.${EntityName}, error) {
	return ${shortEntityName}.GetRaw(id)
}

func (${shortEntityName} *${entityName}Mysql) BatchGetRaw(ids []int64) ([]entity.${EntityName}, error) {
	out := make([]entity.${EntityName}, 0, len(ids))
	err := ${shortEntityName}.dbHelper.FindIn("id", ids, &out)
	return out, err
}

func (${shortEntityName} *${entityName}Mysql) BatchGet(ids []int64) ([]entity.${EntityName}, error) {
	return ${shortEntityName}.BatchGetRaw(ids)
}

func (${shortEntityName} *${entityName}Mysql) Exist(key string, val any) (bool, error) {
	return ${shortEntityName}.dbHelper.Exist(key, val)
}

func (${shortEntityName} *${entityName}Mysql) UpdateMap(id int64, version int, data map[string]interface{}) error {
	err := ${shortEntityName}.dbHelper.UpdateMap(id, version, data)
	if err != nil {
		return err
	}

	return nil
}

func (${shortEntityName} *${entityName}Mysql) UpdateMapTx(tx any, id int64, version int, data map[string]interface{}) (func() error, error) {
	err := ${shortEntityName}.dbHelper.UpdateMapTx(tx, id, version, data)
	if err != nil {
		return nil, err
	}
	return func() error {
		return nil
	}, nil
}

func (${shortEntityName} *${entityName}Mysql) Update(id int64, version int, ${entityName} *entity.${EntityName}) error {
	err := ${shortEntityName}.dbHelper.Update(id, version, ${entityName})
	if err != nil {
		return err
	}
	return nil
}

func (${shortEntityName} *${entityName}Mysql) UpdateTx(tx any, id int64, version int, ${entityName} *entity.${EntityName}) (func() error, error) {
	err := ${shortEntityName}.dbHelper.UpdateTx(tx, id, version, ${entityName})
	if err != nil {
		return nil, err
	}
	return func() error {
		return nil
	}, nil
}

func (${shortEntityName} *${entityName}Mysql) Search(searchMeta contract.SearchMeta, searchParams contract.MysqlBuilder) (data []entity.${EntityName}, cnt int64, err error) {
	err = ${shortEntityName}.dbHelper.Search(searchMeta, searchParams, &data, &cnt)
	return data, cnt, err
}
`

var RepositoryMongoCache = `package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/groupcache/singleflight"
	"github.com/zander-84/seagull/contract"
	"github.com/zander-84/seagull/contract/def"
	"github.com/zander-84/seagull/contrib/cache/wrapcache"
	"github.com/zander-84/seagull/contrib/storage"
	"github.com/zander-84/seagull/drive/mongo"
	"${project}/apps/${server}/internal/entity"
	"time"
)

type ${entityName}Mongo struct {
	db                    *mongo.Mongo
	dbHelper              contract.Mongo
	singleFlight          singleflight.Group
	cache                 contract.Cache
	_cachePrefix          string
	cacheGetOrSetDuration time.Duration
	cacheSetDuration      time.Duration
	version               int
}

func new${EntityName}Mongo(db *mongo.Mongo, cache contract.Cache, cachePrefix string) *${entityName}Mongo {
	out := &${entityName}Mongo{
		db:                    db,
		singleFlight:          singleflight.Group{},
		cache:                 cache,
		_cachePrefix:          cachePrefix,
		cacheGetOrSetDuration:  ${cacheGetOrSetDuration},
		cacheSetDuration:      ${cacheSetDuration},
		version:               ${version},
	}
	out.dbHelper = storage.NewMongo(db.DB(), out.Name(), 50)
	return out
}

func (${shortEntityName}  *${entityName}Mongo) Name() string {
	return "${tableName}"
}

func (${shortEntityName} *${entityName}Mongo) Version() int {
	return ${shortEntityName}.version
}

func (${shortEntityName} *${entityName}Mongo) CachePrefix() string {
	return ${shortEntityName}._cachePrefix
}
func (${shortEntityName} *${entityName}Mongo) assert(in interface{}) (*entity.${EntityName}, error) {
	if _${shortEntityName}, ok := in.(*entity.${EntityName}); !ok || _${shortEntityName} == nil {
		return nil, errors.New("err type")
	} else {
		return _${shortEntityName}, nil
	}
}

func (${shortEntityName} *${entityName}Mongo) Create(${entityName} *entity.${EntityName}) error {
	_, err := ${shortEntityName}.dbHelper.Create(${entityName})
	if err != nil {
		return err
	}
	return nil
}

func (${shortEntityName} *${entityName}Mongo) GetRaw(id string) (*entity.${EntityName}, error) {
	out, err := ${shortEntityName}.singleFlight.Do(fmt.Sprintf("GetRaw:%d", id), func() (interface{}, error) {
		out := new(entity.${EntityName})
		err := ${shortEntityName}.dbHelper.FindByID(id, out)
		return out, err
	})
	if err != nil {
		return nil, err
	}
	return ${shortEntityName}.assert(out)
}

func (${shortEntityName} *${entityName}Mongo) Get(id string) (*entity.${EntityName}, error) {
	out := new(entity.${EntityName})
	err := ${shortEntityName}.cache.GetOrSet(context.Background(), wrapcache.RepoKey(${shortEntityName} , id), out, ${shortEntityName}.cacheGetOrSetDuration, func(key def.K) (value any, err error) {
		return ${shortEntityName}.GetRaw(id)
	})
	return out, err
}
func (${shortEntityName} *${entityName}Mongo) BatchGetRaw(ids []string) ([]entity.${EntityName}, error) {
	out := make([]entity.${EntityName}, 0, len(ids))
	err := ${shortEntityName}.dbHelper.FindIn("_id", ids, &out)
	return out, err
}

func (${shortEntityName} *${entityName}Mongo) BatchGet(ids []string) ([]entity.${EntityName}, error) {
	out := make([]entity.${EntityName}, 0, len(ids))
	err := ${shortEntityName}.cache.BatchGetOrSet(context.Background(), wrapcache.RepoKeys(${shortEntityName} , ids), &out, ${shortEntityName}.cacheSetDuration, func(missIds []def.K) (value map[string]any, err error) {
		missDbIds, err := wrapcache.RepoDBKeys(missIds)
		if err != nil {
			return nil, err
		}
		data := make([]entity.${EntityName}, 0, len(missDbIds))
		objectIDs, err := storage.MongoPkAny2ObjectIDs(missDbIds)
		if err != nil {
			return nil, err
		}

		if err := ${shortEntityName}.dbHelper.FindIn("_id", objectIDs, &data); err != nil {
			return nil, err
		}

		value = make(map[string]any)
		for _, v := range data {
			for _, vv := range missIds {
				if v.PK.Hex() == vv.Alias[0] {
					value[vv.Key] = v
				}
			}
		}

		return value, nil
	})
	return out, err
}

func (${shortEntityName} *${entityName}Mongo) Exist(key string, val any) (bool, error) {
	return ${shortEntityName}.dbHelper.Exist(key, val)
}

func (${shortEntityName} *${entityName}Mongo) Update(id string, version int, ${entityName} *entity.${EntityName}) error {
	err := ${shortEntityName}.dbHelper.Update(id, version, ${entityName})
	if err != nil {
		return err
	}
	err = ${shortEntityName}.cache.Set(context.Background(), wrapcache.RepoKey(${shortEntityName} , id), ${entityName}, ${shortEntityName}.cacheSetDuration)
	if err != nil {
		return err
	}
	return nil
}

func (${shortEntityName} *${entityName}Mongo) Search(searchMeta contract.SearchMeta, searchParams contract.MongoBuilder) (data []entity.${EntityName}, cnt int64, err error) {
	err = ${shortEntityName}.dbHelper.Search(searchMeta, searchParams, &data, &cnt)
	return data, cnt, err
}
`
var RepositoryMongo = `package repository

import (
	"errors"
	"fmt"
	"github.com/golang/groupcache/singleflight"
	"github.com/zander-84/seagull/contract"
	"github.com/zander-84/seagull/contrib/storage"
	"github.com/zander-84/seagull/drive/mongo"
	"${project}/apps/${server}/internal/entity"
)

type ${entityName}Mongo struct {
	db                    *mongo.Mongo
	dbHelper              contract.Mongo
	singleFlight          singleflight.Group
	version               int
}

func new${EntityName}Mongo(db *mongo.Mongo) *${entityName}Mongo {
	out := &${entityName}Mongo{
		db:                    db,
		singleFlight:          singleflight.Group{},
		version:               ${version},
	}
	out.dbHelper = storage.NewMongo(db.DB(), out.Name(), 50)
	return out
}

func (${shortEntityName}  *${entityName}Mongo) Name() string {
	return "${tableName}"
}

func (${shortEntityName} *${entityName}Mongo) Version() int {
	return ${shortEntityName}.version
}


func (${shortEntityName} *${entityName}Mongo) assert(in interface{}) (*entity.${EntityName}, error) {
	if _${shortEntityName}, ok := in.(*entity.${EntityName}); !ok || _${shortEntityName} == nil {
		return nil, errors.New("err type")
	} else {
		return _${shortEntityName}, nil
	}
}

func (${shortEntityName} *${entityName}Mongo) Create(${entityName} *entity.${EntityName}) error {
	_, err := ${shortEntityName}.dbHelper.Create(${entityName})
	if err != nil {
		return err
	}
	return nil
}

func (${shortEntityName} *${entityName}Mongo) GetRaw(id string) (*entity.${EntityName}, error) {
	out, err := ${shortEntityName}.singleFlight.Do(fmt.Sprintf("GetRaw:%d", id), func() (interface{}, error) {
		out := new(entity.${EntityName})
		err := ${shortEntityName}.dbHelper.FindByID(id, out)
		return out, err
	})
	if err != nil {
		return nil, err
	}
	return ${shortEntityName}.assert(out)
}

func (${shortEntityName} *${entityName}Mongo) Get(id string) (*entity.${EntityName}, error) {
	return ${shortEntityName}.GetRaw(id)
}
func (${shortEntityName} *${entityName}Mongo) BatchGetRaw(ids []string) ([]entity.${EntityName}, error) {
	out := make([]entity.${EntityName}, 0, len(ids))
	err := ${shortEntityName}.dbHelper.FindIn("_id", ids, &out)
	return out, err
}

func (${shortEntityName} *${entityName}Mongo) BatchGet(ids []string) ([]entity.${EntityName}, error) {
	return ${shortEntityName}.BatchGetRaw(ids)
}

func (${shortEntityName} *${entityName}Mongo) Exist(key string, val any) (bool, error) {
	return ${shortEntityName}.dbHelper.Exist(key, val)
}

func (${shortEntityName} *${entityName}Mongo) Update(id string, version int, ${entityName} *entity.${EntityName}) error {
	err := ${shortEntityName}.dbHelper.Update(id, version, ${entityName})
	if err != nil {
		return err
	}
	return nil
}

func (${shortEntityName} *${entityName}Mongo) Search(searchMeta contract.SearchMeta, searchParams contract.MongoBuilder) (data []entity.${EntityName}, cnt int64, err error) {
	err = ${shortEntityName}.dbHelper.Search(searchMeta, searchParams, &data, &cnt)
	return data, cnt, err
}
`

var RepoMysqlCache = `package repository

import (
	"github.com/zander-84/seagull/contract"
	"github.com/zander-84/seagull/drive/gormmysql"
	"${project}/apps/${server}/internal/entity"
)

type ${entityName}Repository struct {
	${entityName}Mysql *${entityName}Mysql
}

func New${EntityName}(db *gormmysql.Gdb, cache contract.Cache, cachePrefix string) *${entityName}Repository {
	out := new(${entityName}Repository)
	out.${entityName}Mysql = new${EntityName}Mysql(db, cache, cachePrefix)
	return out
}

func (${shortEntityName}  *${entityName}Repository) Create${EntityName}(${entityName} *entity.${EntityName}) error {
	return ${shortEntityName} .${entityName}Mysql.Create(${entityName})
}
func (${shortEntityName}   *${entityName}Repository)Exist(field string, val any) (bool, error) {
	return ${shortEntityName}.${entityName}Mysql.Exist(field,val)
}
func (${shortEntityName}  *${entityName}Repository) Get${EntityName}(id int64) (*entity.${EntityName}, error) {
	return ${shortEntityName}.${entityName}Mysql.Get(id)
}
func (${shortEntityName} *${entityName}Repository) BatchGet${EntityName}(ids []int64) ([]entity.${EntityName}, error) {
	return ${shortEntityName}.${entityName}Mysql.BatchGet(ids)
}
func (${shortEntityName} *${entityName}Repository) Update${EntityName}(id int64, version int, entity *entity.${EntityName}) error {
	return ${shortEntityName}.${entityName}Mysql.Update(id, version, entity)
}

func (${shortEntityName} *${entityName}Repository) Update${EntityName}Map(id int64, version int, entity map[string]any) error {
	return ${shortEntityName}.${entityName}Mysql.UpdateMap(id, version, entity)
}

func (${shortEntityName} *${entityName}Repository) Search${EntityName}(searchMeta contract.SearchMeta, searchParams contract.MysqlBuilder) (data []entity.${EntityName}, cnt int64, err error) {
	return ${shortEntityName}.${entityName}Mysql.Search(searchMeta, searchParams)
}
`
var RepoMysql = `package repository

import (
	"github.com/zander-84/seagull/contract"
	"github.com/zander-84/seagull/drive/gormmysql"
	"${project}/apps/${server}/internal/entity"
)

type ${entityName}Repository struct {
	${entityName}Mysql *${entityName}Mysql
}

func New${EntityName}(db *gormmysql.Gdb) *${entityName}Repository {
	out := new(${entityName}Repository)
	out.${entityName}Mysql = new${EntityName}Mysql(db)
	return out
}
func (${shortEntityName}  *${entityName}Repository)Exist${EntityName}(field string, val any) (bool, error) {
	return ${shortEntityName}.${entityName}Mysql.Exist(field,val)
}
func (${shortEntityName} *${entityName}Repository) Create${EntityName}(${entityName} *entity.${EntityName}) error {
	return ${shortEntityName}.${entityName}Mysql.Create(${entityName})
}

func (${shortEntityName} *${entityName}Repository) Get${EntityName}(id int64) (*entity.${EntityName}, error) {
	return ${shortEntityName}.${entityName}Mysql.Get(id)
}
func (${shortEntityName} *${entityName}Repository) BatchGet${EntityName}(ids []int64) ([]entity.${EntityName}, error) {
	return ${shortEntityName}.${entityName}Mysql.BatchGet(ids)
}
func (${shortEntityName} *${entityName}Repository) Update${EntityName}(id int64, version int, entity *entity.${EntityName}) error {
	return ${shortEntityName}.${entityName}Mysql.Update(id, version, entity)
}

func (${shortEntityName} *${entityName}Repository) Update${EntityName}Map(id int64, version int, entity map[string]any) error {
	return ${shortEntityName}.${entityName}Mysql.UpdateMap(id, version, entity)
}

func (${shortEntityName} *${entityName}Repository) Search${EntityName}(searchMeta contract.SearchMeta, searchParams contract.MysqlBuilder) (data []entity.${EntityName}, cnt int64, err error) {
	return ${shortEntityName}.${entityName}Mysql.Search(searchMeta, searchParams)
}
`

var RepoMongo = `package repository

import (
	"github.com/zander-84/seagull/contract"
	"github.com/zander-84/seagull/drive/mongo"
	"${project}/apps/${server}/internal/entity"
)

type ${entityName}Repository struct {
	${entityName}Mongo *${entityName}Mongo
}

func New${EntityName}(db *mongo.Mongo) *${entityName}Repository {
	out := new(${entityName}Repository)
	out.${entityName}Mongo  = new${EntityName}Mongo(db)
	return out
}

func (${shortEntityName} *${entityName}Repository) Create${EntityName}(${entityName} *entity.${EntityName}) error {
	return ${shortEntityName}.${entityName}Mongo.Create(${entityName})
}
func (${shortEntityName}  *${entityName}Repository)Exist${EntityName}(field string, val any) (bool, error) {
	return ${shortEntityName}.${entityName}Mongo.Exist(field,val)
}
func (${shortEntityName} *${entityName}Repository) Get${EntityName}(id string) (*entity.${EntityName}, error) {
	return ${shortEntityName}.${entityName}Mongo.Get(id)
}
func (${shortEntityName} *${entityName}Repository) BatchGet${EntityName}(ids []string) ([]entity.${EntityName}, error) {
	return ${shortEntityName}.${entityName}Mongo.BatchGet(ids)
}
func (${shortEntityName} *${entityName}Repository) Update${EntityName}(id string, version int, entity *entity.${EntityName}) error {
	return ${shortEntityName}.${entityName}Mongo.Update(id, version, entity)
}
func (${shortEntityName} *${entityName}Repository) Search${EntityName}(searchMeta contract.SearchMeta, searchParams contract.MongoBuilder) (data []entity.${EntityName}, cnt int64, err error) {
	return ${shortEntityName}.${entityName}Mongo.Search(searchMeta, searchParams)
}
`
var RepoMongoCache = `package repository

import (
	"github.com/zander-84/seagull/contract"
	"github.com/zander-84/seagull/drive/mongo"
	"${project}/apps/${server}/internal/entity"
)

type ${entityName}Repository struct {
	${entityName}Mongo *${entityName}Mongo
}

func New${EntityName}(db *mongo.Mongo, cache contract.Cache, cachePrefix string) *${entityName}Repository {
	out := new(${entityName}Repository)
	out.${entityName}Mongo  = new${EntityName}Mongo(db, cache, cachePrefix)
	return out
}

func (${shortEntityName}  *${entityName}Repository)Exist(field string, val any) (bool, error) {
	return ${shortEntityName}.${entityName}Mongo.Exist(field,val)
}
func (${shortEntityName} *${entityName}Repository) Create${EntityName}(${entityName} *entity.${EntityName}) error {
	return ${shortEntityName}.${entityName}Mongo.Create(${entityName})
}

func (${shortEntityName} *${entityName}Repository) Get${EntityName}(id string) (*entity.${EntityName}, error) {
	return ${shortEntityName}.${entityName}Mongo.Get(id)
}
func (${shortEntityName} *${entityName}Repository) BatchGet${EntityName}(ids []string) ([]entity.${EntityName}, error) {
	return ${shortEntityName}.${entityName}Mongo.BatchGet(ids)
}
func (${shortEntityName} *${entityName}Repository) Update${EntityName}(id string, version int, entity *entity.${EntityName}) error {
	return ${shortEntityName}.${entityName}Mongo.Update(id, version, entity)
}
func (${shortEntityName} *${entityName}Repository) Search${EntityName}(searchMeta contract.SearchMeta, searchParams contract.MongoBuilder) (data []entity.${EntityName}, cnt int64, err error) {
	return ${shortEntityName}.${entityName}Mongo.Search(searchMeta, searchParams)
}
`
