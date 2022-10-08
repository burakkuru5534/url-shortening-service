package api

import (
	"encoding/json"
	"github.com/Shyp/go-dberror"
	"github.com/burakkuru5534/src/auth"
	"github.com/burakkuru5534/src/helper"
	"github.com/burakkuru5534/src/helper/messages"
	"github.com/burakkuru5534/src/model"
	_ "github.com/letsencrypt/boulder/db"
	"net/http"

	"log"
)

func UrlCreate(w http.ResponseWriter, r *http.Request) {

	var Url model.Url

	tc := auth.TokenClaimsFromRequest(r)

	err := helper.UsrUrlLimitCheck(tc.UserID)
	if err != nil {
		log.Println(messages.UsrLimitErrorMessage, err)
		http.Error(w, messages.UsrLimitErrorMessage, http.StatusBadRequest)
		return
	}

	err = helper.BodyToJsonReq(r, &Url)
	if err != nil {
		log.Println(messages.BodyParseErrorMessage, err)
		http.Error(w, messages.BodyParseErrorMessage, http.StatusBadRequest)
		return
	}

	err = Url.Create(tc.UserID)
	if err != nil {
		dberr := dberror.GetError(err)
		switch e := dberr.(type) {
		case *dberror.Error:
			if e.Code == "23505" {
				log.Println(messages.UrlAlreadyExistErrorMessage, err)
				http.Error(w, messages.UrlAlreadyExistErrorMessage, http.StatusForbidden)
				return
			}
		}
		log.Println(messages.UrlCreateErrorMessage, err)
		http.Error(w, messages.ServerErrorMessage, http.StatusInternalServerError)
		return
	}

	resp := model.RespData{
		Success: messages.SuccessTrue,
		Message: messages.UrlCreateSuccessMessage,
		Data:    Url.ShortenedVersion,
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Println(messages.ReturnRespErrorMessage, err)
		http.Error(w, messages.ServerErrorMessage, http.StatusInternalServerError)
		return
	}

}

func UrlList(w http.ResponseWriter, r *http.Request) {

	var url model.Url
	tc := auth.TokenClaimsFromRequest(r)

	Urls, err := url.List(tc.UserID)
	if err != nil {
		log.Println("Url list error: ", err)
		http.Error(w, "{\"error\": \"server error\"}", http.StatusInternalServerError)
		return
	}

	resp := model.RespData{
		Success: messages.SuccessTrue,
		Message: messages.UrlListSuccessMessage,
		Data:    Urls,
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Println(messages.ReturnRespErrorMessage, err)
		http.Error(w, messages.ServerErrorMessage, http.StatusInternalServerError)
		return
	}

}

func UrlDelete(w http.ResponseWriter, r *http.Request) {

	var url model.Url
	tc := auth.TokenClaimsFromRequest(r)

	id := helper.StrToInt64(r.URL.Query().Get("id"))

	err := url.Delete(tc.UserID, id)
	if err != nil {
		log.Println(messages.UrlDeleteErrorMessage, err)
		http.Error(w, messages.ServerErrorMessage, http.StatusInternalServerError)
		return
	}

	resp := model.RespData{
		Success: messages.SuccessTrue,
		Message: messages.UrlDeleteSuccessMessage,
		Data:    id,
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Println(messages.ReturnRespErrorMessage, err)
		http.Error(w, messages.ServerErrorMessage, http.StatusInternalServerError)
		return
	}

}

func UrlGet(w http.ResponseWriter, r *http.Request) {

	var url model.Url
	tc := auth.TokenClaimsFromRequest(r)

	shortenedVersion := r.URL.Query().Get("shortenedVersion")

	longVersion, err := url.GetLongUrlFromShortened(tc.UserID, shortenedVersion)
	if err != nil {
		log.Println(messages.GetLongUrlErrorMessage, err)
		http.Error(w, messages.ServerErrorMessage, http.StatusInternalServerError)
		return
	}

	resp := model.RespData{
		Success: messages.SuccessTrue,
		Message: messages.GetLongUrlSuccessMessage,
		Data:    longVersion,
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Println(messages.ReturnRespErrorMessage, err)
		http.Error(w, messages.ServerErrorMessage, http.StatusInternalServerError)
		return
	}
}
