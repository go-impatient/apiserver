package service

import (
	"fmt"
	// MySQL driver.
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/lexkong/log"
	"github.com/moocss/apiserver/src/config"
)


type Database struct {
	config config.ConfYaml
	Self 		*gorm.DB
	// Docker 	*gorm.DB
}

var DB *Database

func openDB(username, password, addr, name string) *gorm.DB {
	conf := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=%t&loc=%s",
		username,
		password,
		addr,
		name,
		true,
		//"Asia/Shanghai"),
		"Local")

	db, err := gorm.Open("mysql", conf)
	if err != nil {
		log.Errorf(err, "Database connection failed.")
	}

	log.Infof("Database connection succeed.")

	// set for db connection
	setupDB(db)

	return db
}

func setupDB(db *gorm.DB) {
	db.LogMode(true)
	//db.DB().SetMaxOpenConns(20000) // 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误。
	db.DB().SetMaxIdleConns(0) // 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
}

// Init client storage.
func (db *Database) InitSelfDB() *gorm.DB {
	return openDB(db.config.Db.Username, db.config.Db.Password, db.config.Db.Addr, db.config.Db.Name)
}

func (db *Database) GetSelfDB() *gorm.DB {
	return db.InitSelfDB()
}

//func (db *Database) InitDockerDB() *gorm.DB {
//	return openDB(db.config.DockerDb.Username, db.config.DockerDb.Password, db.config.DockerDb.Addr, db.config.DockerDb.Name)
//}
//
//func (db *Database) GetDockerDB() *gorm.DB {
//	return db.InitDockerDB()
//}

func (db *Database) Init(config config.ConfYaml) {
	DB = &Database{
		config: config,
		Self:  db.GetSelfDB(),
		// Docker: db.GetDockerDB(),
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