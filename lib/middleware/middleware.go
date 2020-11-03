package middleware

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
)

var (
	db *gorm.DB
)

func Init(database *gorm.DB) {
	db = database
}

type Handler func(r *http.Request, urlValues map[string]string, session *gorm.DB) (statusCode int, err error, output interface{})

//type PlainHandler func(res http.ResponseWriter, req *http.Request, urlValues map[string]string)

// send a http response to the user with JSON format
func SendResponse(res http.ResponseWriter, statusCode int, data interface{}) {
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.WriteHeader(statusCode)
	if d, ok := data.([]byte); ok {
		res.Write(d)
	} else {
		json.NewEncoder(res).Encode(data)
	}
}

//a middleware to handle user authorization
func Wrap(f Handler) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		session := db.Session(&gorm.Session{PrepareStmt: true})
		if statusCode, err, output := f(req, mux.Vars(req), session); err == nil {
			//the business logic handler return no error, then try to commit the db session
			if err := session.Error; err != nil {
				SendResponse(res, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			} else {
				SendResponse(res, statusCode, output)
			}
		} else {
			SendResponse(res, statusCode, map[string]string{"error": err.Error()})
		}
	}
}
