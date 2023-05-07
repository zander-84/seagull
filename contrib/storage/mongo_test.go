package storage

import (
	"fmt"
	"github.com/zander-84/seagull/drive/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

type student2 struct {
	Id        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Status    int                `bson:"status"`
	Version   int                `bson:"version"`
	CreatedAt int64              `bson:"created_at"`
	UpdatedAt int64              `bson:"updated_at"`

	updateFields map[string]any `bson:"-"`
}

func newStudent2() *student2 {
	out := new(student2)
	out.updateFields = map[string]any{}
	return out
}
func (s *student2) UpdatedFields() map[string]any {
	var out = map[string]interface{}{}
	for k, v := range s.updateFields {
		out[k] = v
	}
	return out
}

func (s *student2) UpdateName(name string) {
	s.updateFields["name"] = name
	s.Name = name
}
func (s *student2) UpdateStatus(status int) {
	s.updateFields["status"] = status
	s.Status = status
}

func (s *student2) UpdateUpdatedAt(updatedAt int64) {
	s.updateFields["updated_at"] = updatedAt
	s.UpdatedAt = updatedAt
}
func (s *student2) UpdateCreatedAt(createdAt int64) {
	s.updateFields["created_at"] = createdAt
	s.CreatedAt = createdAt
}

func (s *student2) UpdateVersion(version int) {
	s.Version = version
}
func (s *student2) GetVersion() int {
	return s.Version
}
func TestNewMongo(t *testing.T) {
	mdb := mongo.NewMongo(mongo.Conf{
		Host:            "172.16.86.160",
		Port:            "27017",
		MaxPoolSize:     100,
		MinPoolSize:     10,
		MaxConnIdleTime: 5,
		Database:        "test",
		User:            "zander",
		Pwd:             "zander",
	})
	if err := mdb.Start(); err != nil {
		t.Fatal(err.Error())
	}

	m := NewMongo(mdb.DB(), "student3", 100)

	//for i := 0; i < 10; i++ {
	//	s := newStudent2()
	//	s.Name = fmt.Sprintf("name-%d", i)
	//	s.Status = i % 128
	//	_, err := m.Create(s)
	//	if err != nil {
	//		t.Fatal(err.Error())
	//	}
	//	//fmt.Println(s, id)
	//}

	searchMeta2 := NewSearchMeta()
	//searchMeta2.SetPage(2)
	//searchMeta2.SetPageSize(3)
	searchMeta2.UseCount(false)
	searchMeta2.UseCursor(true)
	searchMeta2.UsePage(false)

	mysqlBuilder2 := NewMongoBuilder()
	mysqlBuilder2.AppendWhere(bson.E{Key: "version", Value: 1})

	res := make([]student2, 0)
	var cnt int64 = 0
	if err := m.Search(searchMeta2, mysqlBuilder2, &res, &cnt); err != nil {
		t.Fatal(err.Error())
	}
	//
	fmt.Println(cnt, len(res), res)

	//s := newStudent2()
	//if err := m.FindByID("63c1460a8abbf40a24e6c81d", s); err != nil {
	//	t.Fatal(err.Error())
	//}
	//fmt.Println(s)
	//
	//s.Status = 99
	//if err := m.ReplaceOne("63c1460a8abbf40a24e6c81d", s.Version, s); err != nil {
	//	t.Fatal(err.Error())
	//}
	//fmt.Println(s)

	//s1 := newStudent2()
	//if err := m.FindOneByField("name", "name-3", s1); err != nil {
	//	t.Fatal(err.Error())
	//}
	//fmt.Println(s1)
	////
	//s1.SetStatus(88)
	//if err := m.UpdatePart(s1.Id.Hex(), s1.Version, s1); err != nil {
	//	t.Fatal(err.Error())
	//}
	//fmt.Println(s1)
	////
	//s1.Name = "marvin"
	//s1.SetStatus(66)
	//if err := m.UpdatePart(s1.Id.Hex(), s1.Version, s1); err != nil {
	//	t.Fatal(err.Error())
	//}
	//fmt.Println(s1)

	//q := make([]string, 0)
	//for i := 0; i < 1000; i++ {
	//	q = append(q, fmt.Sprintf("name-%d", i))
	//	q = append(q, fmt.Sprintf("name-%d", i))
	//}
	//res := make([]student2, 0)
	//if err := m.FindIn("name", q, &res); err != nil {
	//	t.Fatal(err.Error())
	//}
	//fmt.Println(len(res))
	//fmt.Println(res)

	//q2 := make([]primitive.ObjectID, 0)
	//for _, v := range []string{"63c148f6d545515d9c6322ee", "63c148f6d545515d9c6322ef", "63c148f6d545515d9c6322f0"} {
	//	id, _ := primitive.ObjectIDFromHex(v)
	//	q2 = append(q2, id)
	//}
	//
	//res2 := make([]student2, 0)
	//if err := m.FindIn("_id", q2, &res2); err != nil {
	//	t.Fatal(err.Error())
	//}
	//fmt.Println(len(res2))
	//fmt.Println(res2)

	fmt.Println(m.Exist("name", fmt.Sprintf("name-%d", 1)))
}
