package sys

import (
	"database/sql"
	"encoding/json"
	"github.com/burakkuru5534/src/auth"
	"github.com/burakkuru5534/src/helper"
	"github.com/burakkuru5534/src/helper/messages"
	"github.com/burakkuru5534/src/helper/queries"
	"github.com/burakkuru5534/src/model"
	"log"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {

	var loginData model.LoginData
	var loginInfo model.LoginReqData

	err := helper.BodyToJsonReq(r, &loginInfo)
	if err != nil {
		log.Println(messages.BodyParseErrorMessage, err)
		http.Error(w, messages.BodyParseErrorMessage, http.StatusBadRequest)
		return
	}

	qs := queries.LoginQuery
	err = helper.App.DB.QueryRowx(qs, loginInfo.Email).Scan(&loginData.UserID, &loginData.UserEmail, &loginData.UpassFromDb)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println(messages.UserLoginUserDoesNotExistErrorMessage, err)
			http.Error(w, messages.UserLoginUserDoesNotExistErrorMessage, http.StatusBadRequest)
		} else {
			log.Println(messages.UserLoginErrorMessage, err)
			http.Error(w, messages.ServerErrorMessage, http.StatusInternalServerError)
		}
		return
	}

	// compare password hashes
	loginResult := helper.CheckPass(*loginData.UpassFromDb, loginInfo.Password)
	if !loginResult {
		log.Println(messages.UserLoginWrongPasswordErrorMessage, err)
		http.Error(w, messages.UserLoginWrongPasswordErrorMessage, http.StatusUnauthorized)
		return
	}

	tc, err := auth.NewTokenClaimsForUser(loginData.UserID, loginData.UserEmail)
	if err != nil {
		log.Println(messages.CreateTokenErrorMessage, err)
		http.Error(w, messages.ServerErrorMessage, http.StatusInternalServerError)
		return
	}

	var respStruct struct {
		UserID      int64
		UserEmail   string
		AccessToken string
	}
	respStruct.UserID = loginData.UserID
	respStruct.UserEmail = loginData.UserEmail
	respStruct.AccessToken = tc.Encode(helper.Conf.Auth.JWTAuth)

	_, err = helper.SetLogJwt(respStruct.AccessToken, loginData.UserID)
	if err != nil {
		log.Println(messages.SetLogJwtErrorMessage, err)
		http.Error(w, messages.ServerErrorMessage, http.StatusInternalServerError)
		return
	}

	resp := model.RespData{
		Success: messages.SuccessTrue,
		Message: messages.UserLoginSuccessMessage,
		Data:    respStruct,
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Println(messages.ReturnRespErrorMessage, err)
		http.Error(w, messages.ServerErrorMessage, http.StatusInternalServerError)
		return
	}

}
