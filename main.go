package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ken0911208818/demoTritonHo/handler"
	"github.com/ken0911208818/demoTritonHo/lib/config"
	"github.com/ken0911208818/demoTritonHo/setting"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

func main() {
	initDependency()

	// GOMAXPROCS :最大核心數 NumCPU: 當前編譯環境下的cpu核心數
	runtime.GOMAXPROCS(runtime.NumCPU())

	// 路由器庫
	router := mux.NewRouter()
	// uuid 正則表達式 若不符合則無法通過
	uuidRegexp := `[[:alnum:]]{8}-[[:alnum:]]{4}-4[[:alnum:]]{3}-[89AaBb][[:alnum:]]{3}-[[:alnum:]]{12}`

	router.HandleFunc("/v1/cats/", handler.CatGetAll).Methods("GET")
	router.HandleFunc("/v1/cats/{catId:"+uuidRegexp+"}", handler.CatGetOne).Methods("GET")
	router.HandleFunc("/v1/cats/{catId:"+uuidRegexp+"}", handler.CatUpdate).Methods("PUT")
	router.HandleFunc("/v1/cats/{catId:"+uuidRegexp+"}", handler.CatDelete).Methods("DELETE")
	router.HandleFunc("/v1/cats/", handler.CatCreate).Methods("POST")

	http.Handle("/", router)
	s := &http.Server{
		Addr:         ":7777",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Fatal(s.ListenAndServe())
}

func initDependency() {
	//the postgresql connection string
	connectStr := "host=" + config.GetStr(setting.DB_HOST) +
		" port=" + strconv.Itoa(config.GetInt(setting.DB_PORT)) +
		" dbname=" + config.GetStr(setting.DB_NAME) +
		" user=" + config.GetStr(setting.DB_USERNAME) +
		" password='" + config.GetStr(setting.DB_PASSWORD) + "'" +
		" sslmode=disable"

	db, err := sql.Open("postgres", connectStr)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("連線成功")
	handler.Init(db)
}
