package handler

import (
	"encoding/json"
	"github.com/ken0911208818/demoTritonHo/model"
	"net/http"
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
