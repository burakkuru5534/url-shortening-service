package register

import (
	"encoding/json"
	"fmt"
	"github.com/burakkuru5534/src/helper"
	"github.com/burakkuru5534/src/model"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
)

var (
	validate *validator.Validate
)

func NewRegister(w http.ResponseWriter, r *http.Request) {
	data := model.RegisterRequestData()

	err := helper.BodyToJsonReq(r, &data)
	if err != nil {
		log.Println("Register body to json error: ", err)
		http.Error(w, "{\"error\": \"server error\"}", http.StatusInternalServerError)
		return
	}

	validate = validator.New()
	err = validate.Struct(data)
	if err != nil {
		log.Println("Register body validate error: ", err)
		http.Error(w, "{\"error\": \"Bad request\"}", http.StatusBadRequest)
		return
	}

	data.Password, err = helper.HashPasswd(data.Password)
	if err != nil {
		log.Println("hash password error: ", err)
		http.Error(w, "{\"error\": \"server error\"}", http.StatusInternalServerError)
		return
	}

	err = createSysUsr(data)
	if err != nil {
		log.Println("create sysusr error: ", err)
		http.Error(w, "{\"error\": \"server error\"}", http.StatusInternalServerError)
		return
	}

	resp := fmt.Sprintf("Usr %s created:", data.Email)
	json.NewEncoder(w).Encode(resp)
}

func createSysUsr(data *model.RegisterRequest) error {

	var err error
	shorteningUrlLimit := 1
	accountType := "free"
	if data.IsPremium {
		shorteningUrlLimit = 10
		accountType = "premium"
	}
	qs := "insert into usr ( password, email, shortening_url_limit, account_type, zlins_dttm) values ($1, $2, $3, $4, current_timestamp) returning id"
	err = helper.App.DB.QueryRowx(qs, data.Password, data.Email, shorteningUrlLimit, accountType).Scan(&data.ID)
	if err != nil {
		return err
	}

	return nil
}
