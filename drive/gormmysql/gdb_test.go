package gormmysql

import (
	"gorm.io/gorm"
	"testing"
	"time"
)

// go test -v  -run TestGdb  -args 172.16.86.150  3307 zander zander test2
func TestGdb(t *testing.T) {

	gdb := NewGdb(Conf{
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
	type Product struct {
		ID        uint `gorm:"primarykey"`
		CreatedAt time.Time
		UpdatedAt time.Time
		Code      string
		Price     uint
		Version   int64
	}
	//if err := gdb.Engine().AutoMigrate(&Product{}); err != nil {
	//	t.Fatal(err.Error())
	//}

	//if err := gdb.Engine().Create(&Product{Code: "D42", Price: 100}).Error; err != nil {
	//	t.Fatal(err.Error())
	//}

	//var product Product
	//gdb.Engine().First(&product, 1) // 根据整形主键查找
	//t.Log(product)
	//gdb.Engine().First(&product, "code = ?", "D42") // 查找 code 字段值为 D42 的记录
	//t.Log(product)

	//if err := gdb.Engine().Session(&gorm.Session{SkipHooks: true}).Model(&product).Update("Price", 200).Error; err != nil {
	//	t.Fatal(err.Error())
	//}
	//
	//if err := gdb.Engine().Model(&product).Updates(Product{Price: 200, Code: "F42"}).Error; err != nil {
	//	t.Fatal(err.Error())
	//}
	//
	//if err := gdb.Engine().Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"}).Error; err != nil {
	//	t.Fatal(err.Error())
	//}

	//if err := gdb.Engine().Delete(&product, 1).Error; err != nil {
	//	t.Fatal(err.Error())
	//}
	data := make(map[string]interface{})
	//data["version"] = gorm.Expr("version+?", 1)
	data["code"] = "88"

	data2 := Product{
		Code: "66",
	}
	db := gdb.Engine().Session(&gorm.Session{SkipHooks: true}).Table("product").Where("id=?", 1)
	db.Statement.SetColumn("version", gorm.Expr("version+?", 1), false)
	db.Updates(data2)
	t.Log("success")
}
