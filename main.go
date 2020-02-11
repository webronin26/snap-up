package main

import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo"
	"github.com/webronin26/snap-up/config"
	"github.com/webronin26/snap-up/pkg/presenter"
	"github.com/webronin26/snap-up/pkg/store"
)

func main() {

	// 初始化設定檔案
	config.Init()
	// 初始化資料庫
	store.Init()
	// 初始化 route 需要的資源
	// 這邊先輸入 id = 1 的商品當作 sample
	presenter.InitSnapUp(1)

	e := echo.New()
	e.POST("/snap", presenter.SnapUp)
	e.Logger.Fatal(e.Start(":1323"))
}
