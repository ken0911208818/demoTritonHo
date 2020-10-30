package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ken0911208818/demoTritonHo/lib/httputil"
	"github.com/ken0911208818/demoTritonHo/model"
	"github.com/satori/go.uuid"
	"net/http"
)

func CatGetAll(w http.ResponseWriter, r *http.Request) {
	// create the object slice
	cats := []model.Cat{}

	//load the object from database
	result := db.Find(&cats)

	if err := result.Error; err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{error":"` + err.Error() + `"}"`))
		return
	}

	//output the result
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cats)
}

func CatGetOne(w http.ResponseWriter, r *http.Request) {
	//create the object and get the id from the url
	var cat model.Cat
	cat.Id = mux.Vars(r)[`catId`]
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	//load the object date from database
	result := db.Where(`id = ?`, cat.Id).First(&cat)
	if err := result.Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"` + err.Error() + `"}`))
	}
	//output the object , or any error

	// sql.ErrNoRows = not found any data offset id = url.id
	fmt.Println(result.RowsAffected)
	if result.RowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(cat)
	}
}

func CatUpdate(w http.ResponseWriter, r *http.Request) {
	//create the object and get the id from url
	var cat model.Cat
	cat.Id = mux.Vars(r)[`catId`]

	//since we have to know which field is updated , thus we need to use structure with pointer attribute

	//bind the input
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := httputil.BindForUpdate(r, &cat); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"` + err.Error() + `"}`))
		return
	}
	//perform basic checking on gender
	if cat.Gender != "" && cat.Gender != `MALE` && cat.Gender != `FEMALE` {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"Gender must be MALE or FEMALE"}`))
		return
	}

	result := db.Save(&cat)
	//output the result
	if err := result.Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"` + err.Error() + `"}`))
	} else {
		if result.RowsAffected == 0 {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusNoContent)
		}
	}
}

func CatCreate(w http.ResponseWriter, r *http.Request) {
	// bind header basic
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// bind the input
	cat := model.Cat{}
	if err := httputil.Bind(r, &cat); err != nil {
		errors(w, http.StatusBadRequest, []byte(`{"error":"`+err.Error()+`"}`))
		return
	}

	//perform basic checking on gender
	if cat.Gender != `MALE` && cat.Gender != `FEMALE` {
		errors(w, http.StatusBadRequest, []byte(`{"error":"Gender must be MALE or FEMALE"}`))
	}

	//generate the primary key for the cat
	u := uuid.NewV4()
	cat.Id = u.String()
	//perform the create to the database
	result := db.Create(&cat)

	if err := result.Error; err != nil {
		errors(w, http.StatusInternalServerError, []byte(`{"error":"`+err.Error()+`"}`))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":"` + cat.Id + `"}`))
	}

}

func CatDelete(w http.ResponseWriter, r *http.Request) {
	// bind header basic
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// bind the input
	id := mux.Vars(r)[`catId`]

	//perform the delete to the database
	// db.Delete(new(model.Cat), id) // 帶入參數只支援整數
	result := db.Where(`id = ?`, id).Delete(new(model.Cat))

	//當 result.RowsAffected 進行 update insert or delete 若有資料進行更動時則會傳被引響的資料筆數 若沒有進行更動則回傳0
	if err := result.Error; err != nil {
		errors(w, http.StatusInternalServerError, []byte(`{"error":"`+err.Error()+`"}`))
	} else {
		if result.RowsAffected == 0 {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusNoContent)
		}
	}
}

func errors(w http.ResponseWriter, httpCode int, err []byte) {
	w.WriteHeader(httpCode)
	w.Write(err)
}
