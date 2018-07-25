package service

import (
	"fmt"
	// MySQL driver.
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/lexkong/log"
)



type Database struct {
	Self 		*gorm.DB
	// Docker 	*gorm.DB
}

type confing struct {
	Name     string
	Addr     string
	Username string
	Password string
}

var DB *Database
var conf * confing

func openDB(username, password, addr, name string) *gorm.DB {
	conf := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=%t&loc=%s",
		username,
		password,
		addr,
		name,
		true,
		"Local")

	db, err := gorm.Open("mysql", conf)
	if err != nil {
		log.Errorf(err, "Database connection failed. Database name: %s", name)
	} else {
		log.Infof("Database connection succeed. Database name: %s", name)
	}

	// set for db connection
	setupDB(db)

	return db
}

func setupDB(db *gorm.DB) {
	db.LogMode(true)
	db.DB().SetMaxOpenConns(50) // 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误。
	db.DB().SetMaxIdleConns(10) // 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
}

// Init client storage.
func InitSelfDB() *gorm.DB {
	return openDB(conf.Username, conf.Password, conf.Addr, conf.Name)
}

func GetSelfDB() *gorm.DB {
	return InitSelfDB()
}

//func InitDockerDB() *gorm.DB {
//	return openDB(conf.Username, conf.Password, conf.Addr, conf.Name)
//}
//
//func GetDockerDB() *gorm.DB {
//	return InitDockerDB()
//}

func (db *Database) Init(username, password, addr, name string) {
	conf = &confing{
		Username: username,
		Password: password,
		Addr: addr,
		Name: name,
	}
	DB = &Database{
		Self:  GetSelfDB(),
		// Docker: GetDockerDB(),
	}
}

func (db *Database) Close() {
	if err := DB.Self.Close(); nil != err {
		log.Error("Disconnect from database failed: ", err)
	}

	//if err := DB.Docker.Close(); nil != err {
	//	log.Error("Disconnect from database failed: ", err)
	//}
}