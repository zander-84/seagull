package storage

import (
	"github.com/zander-84/seagull/contract"
	"strings"
)

type mysqlBuilder struct {
	query   []string
	fields  string
	orderBy string
	args    []any
	tag     string
}

func NewMysqlBuilder() contract.MysqlBuilder {
	out := new(mysqlBuilder)
	out.query = make([]string, 0)
	out.args = make([]any, 0)
	return out
}

func (m *mysqlBuilder) AppendWhere(query string, args ...any) contract.MysqlBuilder {
	m.query = append(m.query, query)
	m.args = append(m.args, args...)
	return nil
}

func (m *mysqlBuilder) BuildQuery() string {
	if len(m.query) > 0 {
		return strings.Join(m.query, " And ")
	} else {
		return ""
	}
}

func (m *mysqlBuilder) Args() []any {
	return m.args
}

func (m *mysqlBuilder) SetFields(fields string) contract.MysqlBuilder {
	m.fields = fields
	return m
}

func (m *mysqlBuilder) Fields() string {
	if m.fields == "" {
		return "*"
	}
	return m.fields
}

func (m *mysqlBuilder) SetOrderBy(orderBy string) contract.MysqlBuilder {
	m.orderBy = orderBy
	return m
}

func (m *mysqlBuilder) OrderBy() string {
	return m.orderBy
}

func (m *mysqlBuilder) SetTag(tag string) contract.MysqlBuilder {
	m.tag = tag
	return m
}

func (m *mysqlBuilder) Tag() string {
	return m.tag
}
