package handler

import (
	"errors"
	"github.com/ken0911208818/demoTritonHo/lib/httputil"
	"github.com/ken0911208818/demoTritonHo/model"
	"github.com/satori/go.uuid"
	"gorm.io/gorm"
	"net/http"
)

var errNotFound = errors.New("The record is not found.")

func CatGetAll(r *http.Request, values map[string]string, session *gorm.DB, userId string) (statusCode int, err error, output interface{}) {
	// create the object slice
	cats := []model.Cat{}

	//load the object data from database
	result := session.Where(``).Find(&cats)

	if err := result.Error; err != nil {
		return http.StatusInternalServerError, err, nil
	}

	//output the result
	return http.StatusOK, nil, cats
}

func CatGetOne(r *http.Request, values map[string]string, session *gorm.DB, UserId string) (statusCode int, err error, output interface{}) {
	//create the object and get the id from the url
	var cat model.Cat
	cat.Id = values[`catId`]

	//load the object date from database
	result := session.Where(`id = ?`, cat.Id).First(&cat)
	if err := result.Error; err != nil {
		return http.StatusInternalServerError, err, nil
	}
	//output the object , or any error
	// sql.ErrNoRows = not found any data offset id = url.id
	if result.RowsAffected == 0 {
		return http.StatusNotFound, errNotFound, nil
	}
	return http.StatusOK, nil, cat
}

func CatUpdate(r *http.Request, values map[string]string, session *gorm.DB, UserId string) (statusCode int, err error, output interface{}) {
	//create the object and get the id from url
	var cat model.Cat
	cat.Id = values[`catId`]

	//since we have to know which field is updated , thus we need to use structure with pointer attribute

	//bind the input
	if err := httputil.BindForUpdate(r, &cat); err != nil {
		return http.StatusBadRequest, err, nil
	}

	result := session.Model(&cat).Where(`id = ?`, cat.Id).Updates(cat)
	//output the result
	if err := result.Error; err != nil {
		return http.StatusInternalServerError, err, nil
	} else {
		if result.RowsAffected == 0 {
			return http.StatusNotFound, err, nil
		} else {
			return http.StatusNoContent, nil, nil
		}
	}
}

func CatCreate(r *http.Request, values map[string]string, session *gorm.DB, UserId string) (statusCode int, err error, output interface{}) {
	// bind the input
	cat := model.Cat{}
	if err := httputil.Bind(r, &cat); err != nil {
		return http.StatusBadRequest, err, nil
	}

	//generate the primary key for the cat
	u := uuid.NewV4()
	cat.Id = u.String()
	//perform the create to the database
	result := session.Create(&cat)

	if err := result.Error; err != nil {
		return http.StatusInternalServerError, err, nil
	} else {
		return http.StatusOK, nil, []byte(`{"id":"` + cat.Id + `"}`)
	}

}

func CatDelete(r *http.Request, values map[string]string, session *gorm.DB, UserId string) (statusCode int, err error, output interface{}) {
	// bind the input
	id := values[`catId`]

	//perform the delete to the database
	// db.Delete(new(model.Cat), id) // 帶入參數只支援整數
	result := session.Where(`id = ?`, id).Delete(new(model.Cat))

	//當 result.RowsAffected 進行 update insert or delete 若有資料進行更動時則會傳被引響的資料筆數 若沒有進行更動則回傳0
	if err := result.Error; err != nil {
		return http.StatusInternalServerError, err, nil
	} else {
		if result.RowsAffected == 0 {
			return http.StatusNotFound, err, nil
		} else {
			return http.StatusNoContent, nil, nil
		}
	}
}
