package api

import (
	"encoding/json"
	"github.com/Shyp/go-dberror"
	"github.com/burakkuru5534/src/auth"
	"github.com/burakkuru5534/src/helper"
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
		log.Println("usr url limit check error ", err)
		http.Error(w, "{\"error\": \"You reached to the limit.\"}", http.StatusBadRequest)
		return
	}

	err = helper.BodyToJsonReq(r, &Url)
	if err != nil {
		log.Println("Url create body to json error: ", err)
		http.Error(w, "{\"error\": \"Bad request\"}", http.StatusBadRequest)
		return
	}

	err = Url.Create(tc.UserID)
	if err != nil {
		dberr := dberror.GetError(err)
		switch e := dberr.(type) {
		case *dberror.Error:
			if e.Code == "23505" {
				log.Println("Url already exist error: ", err)
				http.Error(w, "{\"error\": \"Url with that name already exists\"}", http.StatusForbidden)
				return
			}
		}
		log.Println("Url create error: ", err)
		http.Error(w, "{\"error\": \"server error\"}", http.StatusInternalServerError)
		return
	}

	respBody := struct {
		ShortenedVersion string `json:"ShortenedVersion"`
	}{

		ShortenedVersion: Url.ShortenedVersion,
	}

	json.NewEncoder(w).Encode(respBody)

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

	json.NewEncoder(w).Encode(Urls)

}
