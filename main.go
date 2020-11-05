package main

import (
	"crypto/rsa"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/ken0911208818/demoTritonHo/handler"
	"github.com/ken0911208818/demoTritonHo/lib/auth"
	"github.com/ken0911208818/demoTritonHo/lib/config"
	"github.com/ken0911208818/demoTritonHo/lib/middleware"
	"github.com/ken0911208818/demoTritonHo/setting"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io/ioutil"
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

	//auth
	router.HandleFunc("/v1/auth/", middleware.Plain(handler.Login)).Methods("POST")

	//user
	router.HandleFunc("/v1/users/", middleware.Plain(handler.UserCreate)).Methods("POST")
	router.HandleFunc("/v1/users/{userId:"+uuidRegexp+"}", middleware.Wrap(handler.UserUpdate)).Methods("PUT")

	router.HandleFunc("/v1/cats/", middleware.Wrap(handler.CatGetAll)).Methods("GET")
	router.HandleFunc("/v1/cats/{catId:"+uuidRegexp+"}", middleware.Wrap(handler.CatGetOne)).Methods("GET")
	router.HandleFunc("/v1/cats/{catId:"+uuidRegexp+"}", middleware.Wrap(handler.CatUpdate)).Methods("PUT")
	router.HandleFunc("/v1/cats/{catId:"+uuidRegexp+"}", middleware.Wrap(handler.CatDelete)).Methods("DELETE")
	router.HandleFunc("/v1/cats/", middleware.Wrap(handler.CatCreate)).Methods("POST")

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

	//db, err := xorm.NewEngine("postgres", connectStr) //xorm
	db, err := gorm.Open(postgres.New(postgres.Config{DSN: connectStr, PreferSimpleProtocol: true}), &gorm.Config{})
	if err != nil {
		log.Panic("DB connection initialization failed", err)
	}
	sql, _ := db.DB()

	sql.SetMaxIdleConns(config.GetInt(setting.DB_MAX_IDLE_CONN))
	sql.SetMaxOpenConns(config.GetInt(setting.DB_MAX_OPEN_CONN))
	sql.SetConnMaxLifetime(time.Hour)
	sql.Ping()
	//設定連線池數量
	//db.SetMaxIdleConns(config.GetInt(setting.DB_MAX_IDLE_CONN)) //xorm
	//db.SetMaxOpenConns(config.GetInt(setting.DB_MAX_OPEN_CONN)) //xorm
	//db.SetColumnMapper(xormCore.SnakeMapper{}) //xorm
	//uncomment it if you want to debug
	//db.ShowSQL = true
	//db.ShowErr = true
	fmt.Println("連線成功")

	//load the RSA key from the file system ,for the jwt auth
	var err1 error
	var currentKey *rsa.PrivateKey = nil
	var oldKey *rsa.PrivateKey = nil

	currentKeyBytes, _ := ioutil.ReadFile(config.GetStr(setting.JWT_RSA_KEY_LOCATION))
	//解析RSAkey from .pem
	currentKey, err1 = jwt.ParseRSAPrivateKeyFromPEM(currentKeyBytes)
	if err != nil {
		log.Panic(err1)
	}
	if location := config.GetStr(setting.JWT_OLD_RSA_KEY_LOCATION); location != `` {
		oldKeyBytes, _ := ioutil.ReadFile(location)
		oldKey, err1 = jwt.ParseRSAPrivateKeyFromPEM(oldKeyBytes)
		if err1 != nil {
			log.Panic(err1)
		}
	}
	lifetime := time.Duration(config.GetInt(setting.JWT_TOKEN_LIFETIME)) * time.Minute
	auth.Init(currentKey, oldKey, lifetime)
	middleware.Init(db)
}
