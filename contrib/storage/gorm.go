package storage

import (
	"errors"
	"fmt"
	"github.com/zander-84/seagull/contract"
	"github.com/zander-84/seagull/think"
	"github.com/zander-84/seagull/tool"
	"gorm.io/gorm"
	"reflect"
	"time"
)

type gormMysql struct {
	db        *gorm.DB
	tableName string
}

func NewGormMysql(db *gorm.DB, tableName string) contract.Mysql {
	out := new(gormMysql)
	out.db = db
	out.tableName = tableName
	return out
}

func (g *gormMysql) Create(entity contract.IEntity) error {
	entity.SetVersion(1)
	return g.CreateTx(g.db, entity)
}
func (g *gormMysql) CreateTx(tx any, entity contract.IEntity) error {
	tx2, ok := tx.(*gorm.DB)
	if !ok {
		return think.SystemError
	}
	n := time.Now().UnixMilli()
	entity.SetCreatedAt(n)
	entity.SetUpdatedAt(n)
	return tx2.Session(&gorm.Session{SkipHooks: true}).Table(g.tableName).Create(entity).Error
}

func (g *gormMysql) FindByID(id int64, entity contract.IEntity) error {
	err := g.db.Table(g.tableName).Where("id=?", id).First(entity).Error
	if err != nil && gorm.ErrRecordNotFound == err {
		err = think.RecordNotFound
	}
	return err
}

func (g *gormMysql) FindOneByField(field string, val any, entity contract.IEntity) error {
	err := g.db.Table(g.tableName).Where(fmt.Sprintf("`%s`=?", field), val).First(entity).Error
	if err != nil && gorm.ErrRecordNotFound == err {
		err = think.RecordNotFound
	}
	return err
}

func (g *gormMysql) FindIn(field string, val any, ptrSliceData interface{}) error {
	if reflect.ValueOf(ptrSliceData).Type().Kind() != reflect.Ptr {
		return errors.New("data  must be ptr type")
	}
	if reflect.ValueOf(ptrSliceData).Elem().Type().Kind() != reflect.Slice {
		return errors.New("data  must be slice ptr")
	}
	reflectValue := reflect.ValueOf(ptrSliceData).Elem()

	query := make([][]any, 0)
	size := 100
	switch in := val.(type) {
	case []int:
		in = tool.SliceUnique(in)
		query = tool.SliceChunkAny(in, size)
	case []int32:
		in = tool.SliceUnique(in)
		query = tool.SliceChunkAny(in, size)
	case []int64:
		in = tool.SliceUnique(in)
		query = tool.SliceChunkAny(in, size)
	case []string:
		in = tool.SliceUnique(in)
		query = tool.SliceChunkAny(in, size)
	default:
		return errors.New("data  must be slice integer or string")
	}
	for _, v := range query {
		tx := g.db.Session(&gorm.Session{SkipHooks: true})
		rows, err := tx.Table(g.tableName).Where(fmt.Sprintf("`%s` IN ?", field), v).Rows()
		if err != nil {
			return err
		}
		for rows.Next() {
			tmp := reflect.New(reflectValue.Type().Elem())
			if err := tx.ScanRows(rows, tmp.Interface()); err != nil {
				return err
			}

			reflectValue.Set(reflect.Append(reflectValue, tmp.Elem()))
		}
		if err := rows.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (g *gormMysql) Update(id int64, version int64, data map[string]interface{}) error {
	return g.UpdateTx(g.db, id, version, data)
}

func (g *gormMysql) UpdateTx(tx any, id int64, version int64, data map[string]interface{}) error {
	tx2, ok := tx.(*gorm.DB)
	if !ok {
		return think.SystemError
	}

	tx2 = tx2.Session(&gorm.Session{SkipHooks: true})
	data["version"] = gorm.Expr("version+?", 1)
	data["updated_at"] = time.Now().UnixMilli()
	tx2 = tx2.Table(g.tableName).Where("id=?", id)
	if version > 0 {
		tx2 = tx2.Where("version = ?", version)
	}
	res := tx2.Updates(data)
	if res.RowsAffected < 1 {
		return think.RecordNotFound
	}

	return res.Error
}

func (g *gormMysql) UpdatePart(id int64, version int64, entity contract.IEntity) error {
	return g.UpdateTxPart(g.db, id, version, entity)
}

func (g *gormMysql) UpdateTxPart(tx any, id int64, version int64, entity contract.IEntity) error {
	tx2, ok := tx.(*gorm.DB)
	if !ok {
		return think.SystemError
	}

	updateFields := entity.UpdatedFields()
	err := g.UpdateTx(tx2, id, version, updateFields)
	if err != nil {
		return err
	}

	//更新时间+version更新
	if updateAt, ok := updateFields["updated_at"]; ok {
		if updateAt2, ok := updateAt.(int64); ok {
			entity.SetUpdatedAt(updateAt2)
		}
	}

	if version > 0 {
		entity.SetVersion(version + 1)
	}
	return nil
}

func (g *gormMysql) Search(searchMeta contract.SearchMeta, searchParams contract.MysqlBuilder, ptrSliceData interface{}, cnt *int64) (err error) {
	defer func() {
		if recoverError := recover(); recoverError != nil {
			err = fmt.Errorf("recover: %v", recoverError)
		}
	}()
	if reflect.ValueOf(ptrSliceData).Type().Kind() != reflect.Ptr {
		return errors.New("data  must be ptr type")
	}
	if reflect.ValueOf(ptrSliceData).Elem().Type().Kind() != reflect.Slice {
		return errors.New("data  must be slice ptr")
	}

	reflectValue := reflect.ValueOf(ptrSliceData).Elem()
	db := g.db.Session(&gorm.Session{SkipHooks: true}).Table(g.tableName).Select(searchParams.Fields())

	query := searchParams.BuildQuery()
	if query != "" {
		db = db.Where(query, searchParams.Args()...)
	}
	if searchMeta.IsCount() {
		db = db.Count(cnt)
		err = db.Error
		if err != nil {
			return err
		}
	}

	if searchParams.OrderBy() != "" {
		db = db.Order(searchParams.OrderBy())
	}

	if searchMeta.IsPage() {
		db = db.Limit(searchMeta.PageSize()).Offset(searchMeta.Offset())
	}

	if !searchMeta.IsCursor() {
		resDb := db.Find(ptrSliceData)
		err = resDb.Error
		return err
	} else {
		rows, err := db.Rows()
		if err != nil {
			return err
		}

		for rows.Next() {
			tmp := reflect.New(reflectValue.Type().Elem())
			if err := db.ScanRows(rows, tmp.Interface()); err != nil {
				_ = rows.Close()
				return err
			}
			reflectValue.Set(reflect.Append(reflectValue, tmp.Elem()))
		}

		if err := rows.Close(); err != nil {
			return err
		}

		return err
	}
}
