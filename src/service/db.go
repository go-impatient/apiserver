package service

import (
	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/lexkong/log"
	"fmt"
)

type Database struct {
	Self 		*gorm.DB
	// Docker 	*gorm.DB
}

var DB *Database

func realDSN(dbname, username, password, addr string) string {
	conf := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=%t&loc=%s",
		username,
		password,
		addr,
		dbname,
		true,
		"Local")
	return conf
}

func openDB(connStr string) *gorm.DB {
	m, err := gorm.Open("mysql", connStr)
	if err != nil {
		log.Errorf(err, "Database connection failed.")
	} else {
		log.Infof("Database connection succeed.")
	}

	// set for db connection
	setupDB(m)

	return m
}

func setupDB(db *gorm.DB) {
	db.LogMode(true)
	db.DB().SetMaxOpenConns(50) // 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误。
	db.DB().SetMaxIdleConns(10) // 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
}

// Init client storage.
func InitSelfDB(connStr string) *gorm.DB {
	return openDB(connStr)
}

func GetSelfDB(connStr string) *gorm.DB {
	return InitSelfDB(connStr)
}

//func InitDockerDB() *gorm.DB {
//	return openDB(conf.Username, conf.Password, conf.Addr, conf.Name)
//}
//
//func GetDockerDB() *gorm.DB {
//	return InitDockerDB()
//}

func (db *Database) Init(dbname, username, password, addr string) {
	DB = &Database{
		Self: GetSelfDB(realDSN(dbname, username, password, addr)),
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