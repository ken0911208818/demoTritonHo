package middleware

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/ken0911208818/demoTritonHo/lib/auth"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

var (
	db *gorm.DB
)

func Init(database *gorm.DB) {
	db = database
}

type Handler func(r *http.Request, urlValues map[string]string, session *gorm.DB, userId string) (statusCode int, err error, output interface{})
type PlainHandler func(res http.ResponseWriter, req *http.Request, urlValues map[string]string, db *gorm.DB)

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
		authorization := req.Header.Get("Authorization")
		jwtToken := strings.Split(authorization, "Bearer ")[1]
		userId, err := auth.Verify(jwtToken)
		if err != nil {
			SendResponse(res, http.StatusUnauthorized, map[string]string{"error": err.Error()})
			return
		} else {
			//please think carefully on this design, as it has potential security problem
			if newToken, err := auth.Sign(userId); err != nil {
				SendResponse(res, http.StatusInternalServerError, map[string]string{"error": err.Error()})
				return
			} else {
				res.Header().Add("Authorization", newToken) // update JWT Token
			}
		}
		session := db.Session(&gorm.Session{PrepareStmt: true})
		if statusCode, err, output := f(req, mux.Vars(req), session, userId); err == nil {
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

//do nothing and provide injection of database object only
//normally it is used by public endpoint
func Plain(f PlainHandler) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		f(res, req, mux.Vars(req), db)

	}
}
