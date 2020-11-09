package handler

import (
	"github.com/ken0911208818/demoTritonHo/lib/auth"
	"github.com/ken0911208818/demoTritonHo/lib/httputil"
	"github.com/ken0911208818/demoTritonHo/lib/middleware"
	"github.com/ken0911208818/demoTritonHo/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request, urlValues map[string]string, db *gorm.DB) {
	var input struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	if err := httputil.Bind(r, &input); err != nil {
		middleware.SendResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	user := model.User{}
	result := db.Where("email = ?", input.Email).First(&user)
	if err := result.Error; err != nil {
		middleware.SendResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	//密碼比對 || 確認是否有無使用者
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordDigest), []byte(input.Password)); err != nil || result.RowsAffected == 0 {
		middleware.SendResponse(w, http.StatusUnauthorized, map[string]string{"error": "Incorrect Email / Password"})
		return
	}
	if newToken, err := auth.Sign(user.Id); err != nil {
		middleware.SendResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
	} else {
		// update JWT Token
		w.Header().Add("Authorization", newToken)
		//allow CORS
		w.Header().Set("Access-Control-Expose-Headers", "Authorization")
		middleware.SendResponse(w, http.StatusOK, map[string]string{"userId": user.Id, "token": newToken})
	}
}
