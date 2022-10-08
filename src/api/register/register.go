package register

import (
	"encoding/json"
	"fmt"
	"github.com/burakkuru5534/src/helper"
	"github.com/burakkuru5534/src/helper/messages"
	"github.com/burakkuru5534/src/helper/queries"
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
		log.Println(messages.BodyParseErrorMessage, err)
		http.Error(w, messages.BodyParseErrorMessage, http.StatusBadRequest)
		return
	}

	validate = validator.New()
	err = validate.Struct(data)
	if err != nil {
		log.Println(messages.BodyValidateErrorMessage, err)
		http.Error(w, messages.BadDataMessage, http.StatusBadRequest)
		return
	}

	data.Password, err = helper.HashPasswd(data.Password)
	if err != nil {
		log.Println(messages.PasswordHashingErrorMessage, err)
		http.Error(w, messages.ServerErrorMessage, http.StatusInternalServerError)
		return
	}

	err = createUsr(data)
	if err != nil {
		log.Println(messages.UserCreateErrorMessage, err)
		http.Error(w, messages.ServerErrorMessage, http.StatusInternalServerError)
		return
	}

	resp := model.RespData{
		Success: messages.SuccessTrue,
		Message: messages.UserCreateSuccessMessage,
		Data:    fmt.Sprintf("Usr %s created:", data.Email),
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Println(messages.ReturnRespErrorMessage, err)
		http.Error(w, messages.ServerErrorMessage, http.StatusInternalServerError)
		return
	}
}

func createUsr(data *model.RegisterRequest) error {

	var err error
	shorteningUrlLimit := messages.FreeAccountShorteningUrlLimit
	accountType := messages.FreeAccountType
	if data.IsPremium {
		shorteningUrlLimit = messages.PremiumAccountShorteningUrlLimit
		accountType = messages.PremiumAccountType
	}
	qs := queries.CreateUsrQuery
	err = helper.App.DB.QueryRowx(qs, data.Password, data.Email, shorteningUrlLimit, accountType).Scan(&data.ID)
	if err != nil {
		return err
	}

	return nil
}
