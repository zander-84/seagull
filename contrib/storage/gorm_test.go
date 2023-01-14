package storage

import (
	"fmt"
	"github.com/zander-84/seagull/drive/gormmysql"
	"testing"
)

/*
CREATE TABLE `student` (

	`id` bigint unsigned NOT NULL AUTO_INCREMENT,
	`name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户名',
	`status` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '1:正常, 2:删除 ',
	`version` int unsigned NOT NULL DEFAULT '0' COMMENT '版本',
	`created_at` bigint NOT NULL DEFAULT '0',
	`updated_at` bigint NOT NULL DEFAULT '0',
	PRIMARY KEY (`id`)

) ENGINE=InnoDB;
*/

type student struct {
	Id        int64
	Name      string
	Status    int
	Version   int64
	CreatedAt int64
	UpdatedAt int64

	updateFields map[string]any
}

func newStudent() *student {
	out := new(student)
	out.updateFields = map[string]any{}
	return out
}
func (s *student) UpdatedFields() map[string]any {
	var out = map[string]interface{}{}
	for k, v := range s.updateFields {
		out[k] = v
	}
	return out
}

func (s *student) SetName(name string) {
	s.updateFields["name"] = name
	s.Name = name
}
func (s *student) SetStatus(status int) {
	s.updateFields["status"] = status
	s.Status = status
}

func (s *student) SetUpdatedAt(updatedAt int64) {
	s.updateFields["updated_at"] = updatedAt
	s.UpdatedAt = updatedAt
}
func (s *student) SetCreatedAt(createdAt int64) {
	s.updateFields["created_at"] = createdAt
	s.CreatedAt = createdAt
}

func (s *student) SetVersion(version int64) {
	s.Version = version
}

func TestNewGormMysql(t *testing.T) {
	gdb := gormmysql.NewGdb(gormmysql.Conf{
		Host:            "172.16.86.160",
		Port:            "3306",
		User:            "zander",
		Pwd:             "zander",
		Database:        "test2",
		Charset:         "utf8mb4",
		MaxIdleconns:    100,
		MaxOpenconns:    1000,
		ConnMaxLifetime: 300,
		Debug:           true,
		TimeZone:        "",
	})
	if err := gdb.Start(); err != nil {
		t.Fatal(err.Error())
	}
	defer gdb.Stop()

	m := NewGormMysql(gdb.Engine(), "student")

	//for i := 0; i < 10000; i++ {
	//	s := newStudent()
	//	s.Name = fmt.Sprintf("name-%d", i)
	//	s.Status = i % 128
	//	if err := m.Create(s); err != nil {
	//		t.Fatal(err.Error())
	//	}
	//	fmt.Println(s)
	//}
	//s := newStudent()
	//if err := m.FindByID(2, s); err != nil {
	//	t.Fatal(err.Error())
	//}
	//fmt.Println(s)

	//if err := m.Update(4, 2, map[string]interface{}{"status": 99}); err != nil {
	//	t.Fatal(err.Error())
	//}

	//s1 := newStudent()
	//if err := m.FindOneByField("name", "name-3", s1); err != nil {
	//	t.Fatal(err.Error())
	//}
	//fmt.Println(s1)
	//
	//s1.SetStatus(88)
	//if err := m.UpdatePart(s1.Id, s1.Version, s1); err != nil {
	//	t.Fatal(err.Error())
	//}
	//fmt.Println(s1)
	//
	//s1.SetStatus(77)
	//if err := m.UpdatePart(s1.Id, s1.Version, s1); err != nil {
	//	t.Fatal(err.Error())
	//}
	//fmt.Println(s1)

	//searchMeta2 := NewSearchMete()
	//searchMeta2.SetPage(2)
	//searchMeta2.SetPageSize(3)
	//searchMeta2.UseCount(true)
	//searchMeta2.UseCursor(true)
	//searchMeta2.UsePage(true)
	//
	//mysqlBuilder2 := NewMysqlBuilder()
	//mysqlBuilder2.AppendWhere("version=?", 0)
	//

	q := make([]string, 0)
	for i := 0; i < 1000; i++ {
		q = append(q, fmt.Sprintf("%d", i))
		q = append(q, fmt.Sprintf("%d", i))
	}
	res := make([]student, 0)
	if err := m.FindIn("id", q, &res); err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println(len(res))
	fmt.Println(res)

	//res := make([]student, 0)
	//var cnt int64 = 0
	//if err := m.Search(searchMeta2, mysqlBuilder2, &res, &cnt); err != nil {
	//	t.Fatal(err.Error())
	//}
	//fmt.Println(res, cnt)
}
