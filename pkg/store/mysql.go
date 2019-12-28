package store

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/webronin26/snap-up/config"
)

var DB *gorm.DB

func Init() {

	databasesConfig := config.GetDatabaseConfig()

	var err error
	DB, err = gorm.Open("mysql", mysqlSource(&databasesConfig))
	if err != nil {
		panic("init mysql failed: " + err.Error())
	}

	DB.LogMode(databasesConfig.Logmode)
}

func mysqlSource(config *config.DatabaseConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=UTC&multiStatements=true",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Name,
		config.Encoding,
	)
}
