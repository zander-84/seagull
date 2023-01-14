package storage

import (
	"github.com/zander-84/seagull/contract"
	"go.mongodb.org/mongo-driver/bson"
)

type mongoBuilder struct {
	query   bson.D
	fields  bson.D
	orderBy bson.D
	tag     string
}

func NewMongoBuilder() contract.MongoBuilder {
	out := new(mongoBuilder)
	out.query = bson.D{}
	out.fields = bson.D{}
	out.orderBy = bson.D{}
	return out
}

func (m *mongoBuilder) AppendWhere(be bson.E) contract.MongoBuilder {
	m.query = append(m.query, be)
	return m
}

func (m *mongoBuilder) BuildQuery() bson.D {
	return m.query
}

func (m *mongoBuilder) SetFields(fs []string) contract.MongoBuilder {
	for _, v := range fs {
		m.fields = append(m.fields, bson.E{Key: v, Value: 1})
	}
	return m
}

func (m *mongoBuilder) Fields() bson.D {
	return m.fields
}

func (m *mongoBuilder) SetOrderBy(orderBy bson.D) contract.MongoBuilder {
	m.orderBy = orderBy
	return m
}

func (m *mongoBuilder) OrderBy() bson.D {
	return m.orderBy
}

func (m *mongoBuilder) SetTag(tag string) contract.MongoBuilder {
	m.tag = tag
	return m
}

func (m *mongoBuilder) Tag() string {
	return m.tag
}
