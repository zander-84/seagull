package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
	"time"
)

// go test -v  -run TestMongo  -args 172.16.86.150  27017 dbtest
func TestMongo(t *testing.T) {
	//if !flag.Parsed() {
	//	flag.Parse()
	//}
	//argList := flag.Args()
	//user := ""
	//if len(argList) > 3 {
	//	user = argList[3]
	//}
	//pwd := ""
	//if len(argList) > 4 {
	//	pwd = argList[4]
	//}

	mdb := NewMongo(Conf{
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

	defer mdb.Stop()

	type phone struct {
		Number string
		Status int
	}
	//for i := 0; i < 100; i++ {
	//
	//	if _, err := mdb.Collection("m_test").InsertOne(context.Background(), phone{fmt.Sprintf("%d", i), 2}); err != nil {
	//		t.Fatal(err.Error())
	//	}
	//}
	for i := 0; i < 10; i++ {
		go func() {
			cmd := mdb.Collection("m_test").FindOneAndUpdate(context.Background(), bson.D{{"status", 2}}, bson.D{{"$set", bson.D{{"status", 5}}}})
			if cmd.Err() != nil {
				if cmd.Err() == mongo.ErrNoDocuments {
					fmt.Println("mongo.ErrNilDocument")
				}
				fmt.Println(cmd.Err())
				return
			}
			var out = new(phone)
			if err := cmd.Decode(out); err != nil {
				fmt.Println("Decode: ", err.Error())
				return
			}
			fmt.Println("success:", out)
		}()
	}
	time.Sleep(time.Second * 5)
	fmt.Println("fin")
}
