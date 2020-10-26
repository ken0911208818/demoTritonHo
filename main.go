package main

import (
	"database/sql"
	"github.com/ken0911208818/demoTritonHo/handler"
	"github.com/ken0911208818/demoTritonHo/lib/config"
	"github.com/ken0911208818/demoTritonHo/setting"
	_ "github.com/lib/pq"
	"log"
	"strconv"
)

func main() {
	initDependency()
}

func initDependency() {
	//the postgresql connection string
	connectStr := "host=" + config.GetStr(setting.DB_HOST) +
		" ports=" + strconv.Itoa(config.GetInt(setting.DB_PORT)) +
		" dbname=" + config.GetStr(setting.DB_NAME) +
		" user=" + config.GetStr(setting.DB_USERNAME) +
		" password=" + config.GetStr(setting.DB_PASSWORD) +
		" sslmode=disable"

	db, err := sql.Open("postgres", connectStr)
	if err != nil {
		log.Panic(err)
	}
	handler.Init(db)
}