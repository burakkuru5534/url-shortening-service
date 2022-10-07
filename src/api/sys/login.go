package sys

import (
	"database/sql"
	"encoding/json"
	"github.com/burakkuru5534/src/auth"
	"github.com/burakkuru5534/src/helper"
	"github.com/burakkuru5534/src/model"
	"log"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {

	var loginData model.LoginData
	var loginInfo model.LoginReqData

	err := helper.BodyToJsonReq(r, &loginInfo)
	if err != nil {
		log.Println("Login body to json error: ", err)
		http.Error(w, "{\"error\": \"server error\"}", http.StatusInternalServerError)
		return
	}

	qs := "select id, email, password from usr where  email = $1 and is_active"
	err = helper.App.DB.QueryRowx(qs, loginInfo.Email).Scan(&loginData.UserID, &loginData.UserEmail, &loginData.UpassFromDb)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Login user not found: ", err)
			http.Error(w, "{\"error\": \"Bad request\"}", http.StatusBadRequest)
		} else {
			log.Println("Login user query error: ", err)
			http.Error(w, "{\"error\": \"server error\"}", http.StatusInternalServerError)
		}
		return
	}

	// compare password hashes
	loginResult := helper.CheckPass(*loginData.UpassFromDb, loginInfo.Password)
	if !loginResult {
		log.Println("Login password not match: ", err)
		http.Error(w, "{\"error\": \"Bad request\"}", http.StatusBadRequest)
		return
	}

	tc, err := auth.NewTokenClaimsForUser(loginData.UserID, loginData.UserEmail)
	if err != nil {
		log.Println("create new token claims for usr error ", err)
		http.Error(w, "{\"error\": \"server error\"}", http.StatusInternalServerError)
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
		log.Println("set log jwt error: ", err)
		http.Error(w, "{\"error\": \"server error\"}", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(respStruct)
	if err != nil {
		log.Println("return resp: ", err)
		http.Error(w, "{\"error\": \"server error\"}", http.StatusInternalServerError)
		return
	}

}
