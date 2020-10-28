package handler

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/ken0911208818/demoTritonHo/model"
	"github.com/satori/go.uuid"
	"net/http"
	"strconv"
)

func CatGetAll(w http.ResponseWriter, r *http.Request) {
	// create the object slice
	cats := []model.Cat{}

	//load the object from database
	rows, err := db.Query("SELECT id, name, gender, create_time, update_time FROM cats order by id desc")

	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{error":"` + err.Error() + `"}"`))
		return
	}
	defer rows.Close()

	for rows.Next() {
		var cat model.Cat
		if err := rows.Scan(&cat.Id, &cat.Name, &cat.Gender, &cat.CreateTime, &cat.UpdateTime); err != nil {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{error":"` + err.Error() + `"}"`))
			return
		}
		cats = append(cats, cat)
	}
	if err := rows.Err(); err != nil {
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

	//load the object date from database
	err := db.QueryRow("SELECT name, gender, create_time, update_time FROM cats WHERE id = $1::uuid", cat.Id).Scan(&cat.Name, &cat.Gender, &cat.CreateTime, &cat.UpdateTime)

	//output the object , or any error
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// sql.ErrNoRows = not found any data offset id = url.id
	switch err {
	case sql.ErrNoRows:
		w.WriteHeader(http.StatusNotFound)
	case nil:
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(cat)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error"":"` + err.Error() + `"}`))
	}
}

func CatUpdate(w http.ResponseWriter, r *http.Request) {
	//create the object and get the id from url
	var cat model.Cat
	cat.Id = mux.Vars(r)[`catId`]

	//since we have to know which field is updated , thus we need to use structure with pointer attribute
	input := struct {
		Name   *string `json:"name"`
		Gender *string `json:"gender"`
	}{}

	//bind the input
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"` + err.Error() + `"}`))
		return
	}
	//perform basic checking on gender
	if input.Gender != nil && *input.Gender != `MALE` && *input.Gender != `FEMALE` {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"Gender must be MALE or FEMALE"}`))
		return
	}

	//build the SQL for partial update
	columnNames := []string{}
	values := []interface{}{}
	if input.Name != nil {
		columnNames = append(columnNames, `name`)
		values = append(values, input.Name)
	}
	columnNames = append(columnNames, `gender`)
	values = append(values, input.Gender)
	colNamePart := ``
	for i, name := range columnNames {
		colNamePart = colNamePart + name + ` = $` + strconv.Itoa(i+1) + `, `
	}
	q := `UPDATE cats SET ` + colNamePart[0:len(colNamePart)-2] + ` WHERE id = $` + strconv.Itoa(len(columnNames)+1)
	values = append(values, cat.Id)

	//perform the update to the database
	result, err := db.Exec(q, values...)
	//output the result
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"` + err.Error() + `"}`))
	} else {
		if affected, _ := result.RowsAffected(); affected == 0 {
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
	if err := json.NewDecoder(r.Body).Decode(&cat); err != nil {
		errors(w, http.StatusBadRequest, []byte(`{"error":"`+err.Error()+`"}`))
	}

	//perform basic checking on gender
	if cat.Gender != `MALE` && cat.Gender != `FEMALE` {
		errors(w, http.StatusBadRequest, []byte(`{"error":"Gender must be MALE or FEMALE"}`))
	}

	//generate the primary key for the cat
	u := uuid.NewV4()
	cat.Id = u.String()
	//perform the create to the database
	_, err := db.Exec(`insert into cats(id, name, gender) values ($1, $2, $3)`, cat.Id, cat.Name, cat.Gender)

	if err != nil {
		errors(w, http.StatusInternalServerError, []byte(`{"error":"`+err.Error()+`"}`))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":"` + cat.Id + `"}`))
	}

}

func errors(w http.ResponseWriter, httpCode int, err []byte) {
	w.WriteHeader(httpCode)
	w.Write(err)
}
