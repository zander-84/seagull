package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/zander-84/seagull/contract"
	"github.com/zander-84/seagull/think"
	"github.com/zander-84/seagull/tool"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"time"
)

type myMongo struct {
	db        *mongo.Database
	tableName string
	chunkSize int
}

func NewMongo(db *mongo.Database, tableName string, chunkSize int) contract.Mongo {
	out := new(myMongo)
	out.db = db
	out.tableName = tableName
	out.chunkSize = chunkSize
	return out
}

func (m *myMongo) Create(entity contract.IEntity) (string, error) {
	n := time.Now().UnixMilli()
	entity.UpdateVersion(1)
	entity.UpdateCreatedAt(n)
	entity.UpdateUpdatedAt(n)
	res, err := m.db.Collection(m.tableName).InsertOne(context.Background(), entity)
	if err != nil {
		return "", err
	}

	if id, ok := res.InsertedID.(primitive.ObjectID); ok {
		return id.Hex(), nil
	}

	return "", nil
}

func (m *myMongo) FindByID(id string, entity contract.IEntity) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	err = m.db.Collection(m.tableName).FindOne(context.Background(), bson.D{{"_id", objectID}}).Decode(entity)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return think.RecordNotFound
		}
		return err
	}
	return err
}

func (m *myMongo) Exist(field string, val any) (bool, error) {
	res := m.db.Collection(m.tableName).FindOne(context.Background(), bson.D{{Key: field, Value: val}})
	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, res.Err()
	}

	return true, nil
}
func (m *myMongo) FindOneByField(field string, val any, entity contract.IEntity) error {
	err := m.db.Collection(m.tableName).FindOne(context.Background(), bson.D{{field, val}}).Decode(entity)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return think.RecordNotFound
		}
		return err
	}
	return err
}

func (m *myMongo) FindIn(field string, val any, ptrSliceData interface{}) error {
	if reflect.ValueOf(ptrSliceData).Type().Kind() != reflect.Ptr {
		return errors.New("data  must be ptr type")
	}
	if reflect.ValueOf(ptrSliceData).Elem().Type().Kind() != reflect.Slice {
		return errors.New("data  must be slice ptr")
	}

	reflectValue := reflect.ValueOf(ptrSliceData).Elem()
	query := make([][]any, 0)
	size := m.chunkSize
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
	case []primitive.ObjectID:
		in = tool.SliceUnique(in)
		query = tool.SliceChunkAny(in, size)
	case []any:
		in = tool.SliceUniqueAny(in)
		query = tool.SliceChunkAny2Any(in, size)
	default:
		return errors.New("data  must be slice integer or string")
	}
	for _, v := range query {
		cursor, err := m.db.Collection(m.tableName).Find(
			context.Background(),
			bson.D{{field, bson.D{{"$in", v}}}},
		)
		if err != nil {
			return err
		}

		for cursor.Next(context.Background()) {
			tmp := reflect.New(reflectValue.Type().Elem())
			if err = cursor.Decode(tmp.Interface()); err != nil {
				cursor.Close(context.Background())
				return err
			}
			reflectValue.Set(reflect.Append(reflectValue, tmp.Elem()))
		}
		err = cursor.Close(context.Background())
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *myMongo) ReplaceOne(id string, version int, entity contract.IEntity) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	n := time.Now().UnixMilli()
	entity.UpdateUpdatedAt(n)
	where := bson.D{{"_id", objectID}}

	if version > 0 {
		entity.UpdateVersion(version + 1)
		where = append(where, bson.E{Key: "version", Value: version})
	} else {
		entity.UpdateVersion(entity.GetVersion() + 1)
	}

	res, err := m.db.Collection(m.tableName).ReplaceOne(context.Background(), where, entity)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 || res.ModifiedCount == 0 {
		return think.RecordNotFound
	}
	return nil
}
func (m *myMongo) ReplaceOneByKv(field string, val any, version int, entity contract.IEntity) error {
	entity.UpdateUpdatedAt(time.Now().UnixMilli())
	where := bson.D{{field, val}}
	if version > 0 {
		entity.UpdateVersion(version + 1)
		where = append(where, bson.E{Key: "version", Value: version})
	} else {
		entity.UpdateVersion(entity.GetVersion() + 1)
	}

	res, err := m.db.Collection(m.tableName).ReplaceOne(context.Background(), where, entity)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 || res.ModifiedCount == 0 {
		return think.RecordNotFound
	}
	return nil
}

func (m *myMongo) Update(id string, version int, entity contract.IEntity) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return m.UpdateByKv("_id", objectID, version, entity)
}
func (m *myMongo) UpdateByKv(field string, val any, version int, entity contract.IEntity) error {

	n := time.Now().UnixMilli()
	entity.UpdateUpdatedAt(n)
	where := bson.D{{field, val}}
	if version > 0 {
		where = append(where, bson.E{Key: "version", Value: version})
	}
	updateFields := entity.UpdatedFields()
	updateData := bson.D{}
	for k, v := range updateFields {
		if k == "version" {
			continue
		}
		updateData = append(updateData, bson.E{Key: k, Value: v})
	}

	res, err := m.db.Collection(m.tableName).UpdateOne(context.Background(), where, bson.D{{"$set", updateData}, {"$inc", bson.D{{"version", 1}}}})
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 || res.ModifiedCount == 0 {
		return think.RecordNotFound
	}
	if version > 0 {
		entity.UpdateVersion(version + 1)
	}
	return nil
}
func (m *myMongo) DelOneByKv(key string, val any) error {
	_, err := m.db.Collection(m.tableName).DeleteOne(context.Background(), bson.D{{Key: key, Value: val}})
	return err
}

func (m *myMongo) DelByKv(key string, val any) error {
	_, err := m.db.Collection(m.tableName).DeleteMany(context.Background(), bson.D{{Key: key, Value: val}})
	return err
}

func (m *myMongo) DelByCondition(builder contract.MongoBuilder) error {
	_, err := m.db.Collection(m.tableName).DeleteMany(context.Background(), builder.BuildQuery())
	return err
}

func (m *myMongo) Search(searchMeta contract.SearchMeta, searchParams contract.MongoBuilder, ptrSliceData interface{}, cnt *int64) (err error) {
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
	query := searchParams.BuildQuery()
	if searchMeta.IsCount() {
		total, err := m.db.Collection(m.tableName).CountDocuments(
			context.Background(),
			query,
		)
		if err != nil {
			return err
		}
		*cnt = total
	}

	ops := new(options.FindOptions)
	if len(searchParams.Fields()) > 0 {
		ops.SetProjection(searchParams.Fields())
	}

	if searchMeta.IsPage() {
		ops.SetLimit(int64(searchMeta.PageSize()))
		ops.SetSkip(int64(searchMeta.Offset()))
	}
	if len(searchParams.OrderBy()) > 0 {
		ops.SetSort(searchParams.OrderBy())
		//ops.SetAllowDiskUse(true)
	}

	cursor, err := m.db.Collection(m.tableName).Find(
		context.Background(),
		query,
		ops,
	)
	if err != nil {
		return err
	}
	if !searchMeta.IsCursor() {
		defer cursor.Close(context.Background())
		if err := cursor.All(context.Background(), ptrSliceData); err != nil {
			return err
		}
	} else {
		tmp := make([]bson.Raw, 0)
		for cursor.Next(context.Background()) {
			tmp = append(tmp, cursor.Current)

			//tmp := reflect.New(reflectValue.Type().Elem())
			//if err = cursor.Decode(tmp.Interface()); err != nil {
			//	cursor.Close(context.Background())
			//	return err
			//}
			//reflectValue.Set(reflect.Append(reflectValue, tmp.Elem()))
		}
		err = cursor.Close(context.Background())
		if err != nil {
			return err
		}
		for k, _ := range tmp {
			rowData := reflect.New(reflectValue.Type().Elem())
			if err := bson.UnmarshalWithRegistry(bson.DefaultRegistry, tmp[k], rowData.Interface()); err != nil {
				return err
			}
			reflectValue.Set(reflect.Append(reflectValue, rowData.Elem()))
		}
	}

	return err
}
