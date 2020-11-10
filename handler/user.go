package handler

import (
	"errors"
	"github.com/gofrs/uuid"
	"github.com/ken0911208818/demoTritonHo/lib/auth"
	"github.com/ken0911208818/demoTritonHo/lib/httputil"
	"github.com/ken0911208818/demoTritonHo/lib/middleware"
	"github.com/ken0911208818/demoTritonHo/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
)

func UserCreate(w http.ResponseWriter, r *http.Request, urlValues map[string]string, db *gorm.DB) {
	user := struct {
		model.User
		Password string `json:"password" validate:"required"`
	}{}

	if err := httputil.Bind(r, &user); err != nil {
		middleware.SendResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	//generate new uuid for UserId
	u, _ := uuid.NewV4()
	user.Id = u.String()
	if digest, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost); err != nil {
		middleware.SendResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
	} else {
		user.PasswordDigest = string(digest)
	}
	q := ` insert into users(id, email, password_digest, first_name, last_name)
			select ?, ?, ?, ?, ?
			where not exists (select 1 from users where email = ?)`
	result := db.Exec(q, user.Id, user.Email, user.PasswordDigest, user.FirstName, user.LastName, user.Email)
	if err := result.Error; err != nil {
		middleware.SendResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	if attracted := result.RowsAffected; attracted == 0 {
		middleware.SendResponse(w, http.StatusForbidden, map[string]string{"error": "The email is already used."})
		return
	}

	if newToken, err := auth.Sign(user.Id); err != nil {
		middleware.SendResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	} else {
		// update JWT Token
		w.Header().Add("Authorization", newToken)
		//allow CORS
		w.Header().Set("Access-Control-Expose-Headers", "Authorization")
		middleware.SendResponse(w, http.StatusOK, map[string]string{"userId": user.Id})
	}
}

func UserUpdate(r *http.Request, urlValues map[string]string, db *gorm.DB, userId string) (statusCode int, err error, output interface{}) {
	id := urlValues[`userId`]
	if id != userId {
		return http.StatusForbidden, errors.New("Updating others account is forbidden"), nil
	}
	input := struct {
		model.User
		Password         *string `json:"password" validate:"omitempty,min=1"`
		OriginalPassword *string `json:"originalPassword" validate:"omitempty,min=1"`
	}{}

	err = httputil.BindForUpdate(r, &input)

	if err != nil {
		return http.StatusBadRequest, err, nil
	}
	user := model.User{}
	if input.Password != nil {
		if input.OriginalPassword == nil {
			return http.StatusForbidden, errors.New("Please provide the original password"), nil
		}
		// find user

		result := db.Where(`id = ?`, userId).First(&user)
		if err := result.Error; err != nil {
			return http.StatusInternalServerError, err, nil
		}
		//output the object , or any error
		// sql.ErrNoRows = not found any data offset id = url.id
		if result.RowsAffected == 0 {
			return http.StatusNotFound, errNotFound, nil
		}

		// check user.password  == input.password
		if bcrypt.CompareHashAndPassword([]byte(user.PasswordDigest), []byte(*input.OriginalPassword)) != nil {
			return http.StatusForbidden, errors.New(`The original password is invalid`), nil
		}

		// update user.password

		newPassword, err := bcrypt.GenerateFromPassword([]byte(*input.Password), bcrypt.DefaultCost)
		if err != nil {
			return http.StatusInternalServerError, err, nil
		}
		input.PasswordDigest = string(newPassword)
	}
	result := db.Debug().Model(&user).Where(`id = ?`, userId).Updates(input.User)
	if err := result.Error; err != nil {
		return http.StatusInternalServerError, err, nil
	}
	//output the object , or any error
	// sql.ErrNoRows = not found any data offset id = url.id
	if result.RowsAffected == 0 {
		return http.StatusNotFound, errNotFound, nil
	}

	return http.StatusNoContent, nil, nil
}
