package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/burakkuru5534/src/api"
	"github.com/burakkuru5534/src/api/register"
	"github.com/burakkuru5534/src/api/sys"
	"github.com/burakkuru5534/src/auth"
	"github.com/burakkuru5534/src/helper"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegister(t *testing.T) {

	helper.InitConf()
	conInfo := helper.PgConnectionInfo{
		Host:     "127.0.0.1",
		Port:     5432,
		Database: "soft-robotics",
		Username: "postgres",
		Password: "tayitkan",
		SSLMode:  "disable",
	}

	db, err := helper.NewPgSqlxDbHandle(conInfo, 10)
	if err != nil {
		errors.New("create db handle error.")
	}
	err = db.Ping()
	if err != nil {
		errors.New("ping db error.")
	}

	// Create Appplication Service
	err = helper.InitApp(db)
	if err != nil {
		errors.New("init app error.")
	}

	reqUserName := "Crazy"
	ReqUserLastName := "Tom"

	var jsonStrForRegister = []byte(fmt.Sprintf(`{"FirstName":"%s","MiddleName":"","LastName":"%s","Email":"testcrazytom@gmail.com","Password":"tomhanks123"}`, reqUserName, ReqUserLastName))
	reqRegister, err := http.NewRequest("POST", "/sys/register", bytes.NewBuffer(jsonStrForRegister))
	if err != nil {
		t.Fatal(err)
	}
	reqRegister.Header.Set("Content-Type", "application/json")

	rrRegister := httptest.NewRecorder()
	handlerRegister := http.HandlerFunc(register.NewRegister)
	handlerRegister.ServeHTTP(rrRegister, reqRegister)

	if status := rrRegister.Code; status != http.StatusOK {
		if status == http.StatusBadRequest {
			expected := `{"error": "Bad request"}
`
			if rrRegister.Body.String() != expected {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rrRegister.Body.String(), expected)
			}
		} else if status == http.StatusForbidden {
			expected := `{"error": "Movie with that email already exists"}
`
			if rrRegister.Body.String() != expected {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rrRegister.Body.String(), expected)
			}
		} else if status == http.StatusInternalServerError {
			expected := `{"error": "Internal server error"}
`
			if rrRegister.Body.String() != expected {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rrRegister.Body.String(), expected)
			}
		} else if status == http.StatusNotFound {
			expected := `{"error": "Movie with that id does not exist"}
`
			if rrRegister.Body.String() != expected {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rrRegister.Body.String(), expected)
			}
		} else {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
	} else {
		expected := fmt.Sprintf(`"User %s.%s created:"
`, reqUserName, ReqUserLastName)
		if rrRegister.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rrRegister.Body.String(), expected)
		}
	}
	// Check the response body is what we expect.

}

func TestLogin(t *testing.T) {

	helper.InitConf()
	helper.Conf.Auth = auth.NewAuth("2GcQCe7SuKxbaA3NSMBy8ztBPbfDsXJ4", false)
	conInfo := helper.PgConnectionInfo{
		Host:     "127.0.0.1",
		Port:     5432,
		Database: "soft-robotics",
		Username: "postgres",
		Password: "tayitkan",
		SSLMode:  "disable",
	}

	db, err := helper.NewPgSqlxDbHandle(conInfo, 10)
	if err != nil {
		errors.New("create db handle error.")
	}
	err = db.Ping()
	if err != nil {
		errors.New("ping db error.")
	}

	// Create Appplication Service
	err = helper.InitApp(db)
	if err != nil {
		errors.New("init app error.")
	}

	var jsonStrForLogin = []byte(`{"username":"Burak.Kuru","password":"burakkuru123"}`)
	reqLogin, err := http.NewRequest("POST", "/sys/login", bytes.NewBuffer(jsonStrForLogin))
	if err != nil {
		t.Fatal(err)
	}
	reqLogin.Header.Set("Content-Type", "application/json")

	rrLogin := httptest.NewRecorder()
	handlerLogin := http.HandlerFunc(sys.Login)
	handlerLogin.ServeHTTP(rrLogin, reqLogin)
	var respStruct struct {
		UserID       int64
		UserName     string
		UserEmail    string
		UserFullName string
		AccessToken  string
	}

	err = json.Unmarshal(rrLogin.Body.Bytes(), &respStruct)
	if err != nil {
		t.Fatal(err)
	}

	if status := rrLogin.Code; status != http.StatusOK {
		if status == http.StatusBadRequest {
			expected := `{"error": "Bad request"}
`
			if rrLogin.Body.String() != expected {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rrLogin.Body.String(), expected)
			}
		} else if status == http.StatusForbidden {
			expected := `{"error": "Movie with that email already exists"}
`
			if rrLogin.Body.String() != expected {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rrLogin.Body.String(), expected)
			}
		} else if status == http.StatusInternalServerError {
			expected := `{"error": "Internal server error"}
`
			if rrLogin.Body.String() != expected {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rrLogin.Body.String(), expected)
			}
		} else if status == http.StatusNotFound {
			expected := `{"error": "Movie with that id does not exist"}
`
			if rrLogin.Body.String() != expected {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rrLogin.Body.String(), expected)
			}
		} else {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
	} else {
		expected := fmt.Sprintf(`{"UserID":6,"UserName":"Burak.Kuru","UserEmail":"brkkr5534@gmail.com","UserFullName":"Burak Kuru","AccessToken":"%s"}
`, respStruct.AccessToken)
		if rrLogin.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rrLogin.Body.String(), expected)
		}
	}
	// Check the response body is what we expect.

}

func TestList(t *testing.T) {

	conInfo := helper.PgConnectionInfo{
		Host:     "127.0.0.1",
		Port:     5432,
		Database: "soft-robotics",
		Username: "postgres",
		Password: "tayitkan",
		SSLMode:  "disable",
	}

	db, err := helper.NewPgSqlxDbHandle(conInfo, 10)
	if err != nil {
		errors.New("create db handle error.")
	}
	err = db.Ping()
	if err != nil {
		errors.New("ping db error.")
	}

	// Create Appplication Service
	err = helper.InitApp(db)
	if err != nil {
		errors.New("init app error.")
	}

	req, err := http.NewRequest("GET", "/api/movies", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.MovieList)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		if status == http.StatusBadRequest {
			expected := `{"error": "Bad request"}
`
			if rr.Body.String() != expected {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), expected)
			}
		} else if status == http.StatusForbidden {
			expected := `{"error": "Movie with that email already exists"}
`
			if rr.Body.String() != expected {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), expected)
			}
		} else if status == http.StatusInternalServerError {
			expected := `{"error": "Internal server error"}
`
			if rr.Body.String() != expected {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), expected)
			}
		} else if status == http.StatusNotFound {
			expected := `{"error": "Movie with that id does not exist"}
`
			if rr.Body.String() != expected {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), expected)
			}
		} else {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
	} else {
		expected := `[{"id":4,"name":"The Lord Of The Rings","description":"desc","typ":"fantasy"},{"id":5,"name":"Harry Potter","description":"desc","typ":"fantasy"}]
`
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), expected)
		}
	}
	// Check the response body is what we expect.

}

func TestCreate(t *testing.T) {
	helper.InitConf()
	helper.Conf.Auth = auth.NewAuth("2GcQCe7SuKxbaA3NSMBy8ztBPbfDsXJ4", false)
	conInfo := helper.PgConnectionInfo{
		Host:     "127.0.0.1",
		Port:     5432,
		Database: "soft-robotics",
		Username: "postgres",
		Password: "tayitkan",
		SSLMode:  "disable",
	}

	db, err := helper.NewPgSqlxDbHandle(conInfo, 10)
	if err != nil {
		errors.New("create db handle error.")
	}
	err = db.Ping()
	if err != nil {
		errors.New("ping db error.")
	}

	// Create Appplication Service
	err = helper.InitApp(db)
	if err != nil {
		errors.New("init app error.")
	}

	var jsonStrForLogin = []byte(`{"username":"Burak.Kuru","password":"burakkuru123"}`)
	reqLogin, err := http.NewRequest("POST", "/sys/login", bytes.NewBuffer(jsonStrForLogin))
	if err != nil {
		t.Fatal(err)
	}
	reqLogin.Header.Set("Content-Type", "application/json")

	rrLogin := httptest.NewRecorder()
	handlerLogin := http.HandlerFunc(sys.Login)
	handlerLogin.ServeHTTP(rrLogin, reqLogin)
	var respStruct struct {
		UserID       int64
		UserName     string
		UserEmail    string
		UserFullName string
		AccessToken  string
	}

	err = json.Unmarshal(rrLogin.Body.Bytes(), &respStruct)
	if err != nil {
		t.Fatal(err)
	}

	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + respStruct.AccessToken
	// add authorization header to the req
	var jsonStr = []byte(`{"name":"MYTESTMOVIE","description":"desc","typ":"action"}`)
	req, err := http.NewRequest("POST", "/api/movie", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", bearer)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.MovieCreate)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		if status == http.StatusBadRequest {
			expected := `{"error": "Bad request"}
`
			if rr.Body.String() != expected {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), expected)
			}
		} else if status == http.StatusForbidden {
			expected := `{"error": "Movie with that email already exists"}
`
			if rr.Body.String() != expected {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), expected)
			}
		} else if status == http.StatusInternalServerError {
			expected := `{"error": "server error"}
`
			if rr.Body.String() != expected {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), expected)
			}
		} else if status == http.StatusNotFound {
			expected := `{"error": "Movie with that id does not exist"}
`
			if rr.Body.String() != expected {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), expected)
			}
		} else {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
	} else {
		var id int64

		err = db.Get(&id, "SELECT id from movie order by id desc limit 1")
		if err != nil {
			errors.New("get id error.")
		}

		expected := fmt.Sprintf(`{"id":%d,"name":"MYTESTMOVIE","description":"desc","typ":"action"}
`, id)
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), expected)
		}
	}

}

func TestGet(t *testing.T) {

	helper.InitConf()
	helper.Conf.Auth = auth.NewAuth("2GcQCe7SuKxbaA3NSMBy8ztBPbfDsXJ4", false)
	conInfo := helper.PgConnectionInfo{
		Host:     "127.0.0.1",
		Port:     5432,
		Database: "soft-robotics",
		Username: "postgres",
		Password: "tayitkan",
		SSLMode:  "disable",
	}

	db, err := helper.NewPgSqlxDbHandle(conInfo, 10)
	if err != nil {
		errors.New("create db handle error.")
	}
	err = db.Ping()
	if err != nil {
		errors.New("ping db error.")
	}

	// Create Appplication Service
	err = helper.InitApp(db)
	if err != nil {
		errors.New("init app error.")
	}

	var jsonStrForLogin = []byte(`{"username":"Burak.Kuru","password":"burakkuru123"}`)
	reqLogin, err := http.NewRequest("POST", "/sys/login", bytes.NewBuffer(jsonStrForLogin))
	if err != nil {
		t.Fatal(err)
	}
	reqLogin.Header.Set("Content-Type", "application/json")

	rrLogin := httptest.NewRecorder()
	handlerLogin := http.HandlerFunc(sys.Login)
	handlerLogin.ServeHTTP(rrLogin, reqLogin)
	var respStruct struct {
		UserID       int64
		UserName     string
		UserEmail    string
		UserFullName string
		AccessToken  string
	}

	err = json.Unmarshal(rrLogin.Body.Bytes(), &respStruct)
	if err != nil {
		t.Fatal(err)
	}

	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + respStruct.AccessToken

	req, err := http.NewRequest("GET", "/api/movie", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Authorization", bearer)
	q := req.URL.Query()
	q.Add("id", "27")
	req.URL.RawQuery = q.Encode()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.MovieGet)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		if status == http.StatusBadRequest {
			expected := `{"error": "Bad request"}
`
			if rr.Body.String() != expected {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), expected)
			}
		} else if status == http.StatusForbidden {
			expected := `{"error": "Movie with that email already exists"}
`
			if rr.Body.String() != expected {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), expected)
			}
		} else if status == http.StatusInternalServerError {
			expected := `{"error": "Internal server error"}
`
			if rr.Body.String() != expected {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), expected)
			}
		} else if status == http.StatusNotFound {
			expected := `{"error": "Movie with that id does not exist"}
`
			if rr.Body.String() != expected {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), expected)
			}
		} else {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
	} else {
		expected := `{"id":27,"name":"MYTESTMOVIE","description":"desc","typ":"action"}
`
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), expected)
		}
	}
	// Check the response body is what we expect.

}

func TestDelete(t *testing.T) {

	helper.InitConf()
	helper.Conf.Auth = auth.NewAuth("2GcQCe7SuKxbaA3NSMBy8ztBPbfDsXJ4", false)
	conInfo := helper.PgConnectionInfo{
		Host:     "127.0.0.1",
		Port:     5432,
		Database: "soft-robotics",
		Username: "postgres",
		Password: "tayitkan",
		SSLMode:  "disable",
	}

	db, err := helper.NewPgSqlxDbHandle(conInfo, 10)
	if err != nil {
		errors.New("create db handle error.")
	}
	err = db.Ping()
	if err != nil {
		errors.New("ping db error.")
	}

	// Create Appplication Service
	err = helper.InitApp(db)
	if err != nil {
		errors.New("init app error.")
	}

	var jsonStrForLogin = []byte(`{"username":"Burak.Kuru","password":"burakkuru123"}`)
	reqLogin, err := http.NewRequest("POST", "/sys/login", bytes.NewBuffer(jsonStrForLogin))
	if err != nil {
		t.Fatal(err)
	}
	reqLogin.Header.Set("Content-Type", "application/json")

	rrLogin := httptest.NewRecorder()
	handlerLogin := http.HandlerFunc(sys.Login)
	handlerLogin.ServeHTTP(rrLogin, reqLogin)
	var respStruct struct {
		UserID       int64
		UserName     string
		UserEmail    string
		UserFullName string
		AccessToken  string
	}

	err = json.Unmarshal(rrLogin.Body.Bytes(), &respStruct)
	if err != nil {
		t.Fatal(err)
	}

	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + respStruct.AccessToken
	req, err := http.NewRequest("DELETE", "/api/movie", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Authorization", bearer)

	q := req.URL.Query()
	q.Add("id", "26")
	req.URL.RawQuery = q.Encode()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.MovieDelete)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		if status == http.StatusBadRequest {
			expected := `{"error": "Bad request"}
`
			if rr.Body.String() != expected {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), expected)
			}
		} else if status == http.StatusForbidden {
			expected := `{"error": "Movie with that email already exists"}
`
			if rr.Body.String() != expected {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), expected)
			}
		} else if status == http.StatusInternalServerError {
			expected := `{"error": "Internal server error"}
`
			if rr.Body.String() != expected {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), expected)
			}
		} else if status == http.StatusNotFound {
			expected := `{"error": "Movie with that id does not exist"}
`
			if rr.Body.String() != expected {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), expected)
			}
		} else {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

	} else {
		// Check the response body is what we expect.
		expected := `"Movie deleted"
`
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), expected)
		}
	}

}

func TestUpdate(t *testing.T) {

	helper.InitConf()
	helper.Conf.Auth = auth.NewAuth("2GcQCe7SuKxbaA3NSMBy8ztBPbfDsXJ4", false)
	conInfo := helper.PgConnectionInfo{
		Host:     "127.0.0.1",
		Port:     5432,
		Database: "soft-robotics",
		Username: "postgres",
		Password: "tayitkan",
		SSLMode:  "disable",
	}

	db, err := helper.NewPgSqlxDbHandle(conInfo, 10)
	if err != nil {
		errors.New("create db handle error.")
	}
	err = db.Ping()
	if err != nil {
		errors.New("ping db error.")
	}

	// Create Appplication Service
	err = helper.InitApp(db)
	if err != nil {
		errors.New("init app error.")
	}

	var jsonStrForLogin = []byte(`{"username":"Burak.Kuru","password":"burakkuru123"}`)
	reqLogin, err := http.NewRequest("POST", "/sys/login", bytes.NewBuffer(jsonStrForLogin))
	if err != nil {
		t.Fatal(err)
	}
	reqLogin.Header.Set("Content-Type", "application/json")

	rrLogin := httptest.NewRecorder()
	handlerLogin := http.HandlerFunc(sys.Login)
	handlerLogin.ServeHTTP(rrLogin, reqLogin)
	var respStruct struct {
		UserID       int64
		UserName     string
		UserEmail    string
		UserFullName string
		AccessToken  string
	}

	err = json.Unmarshal(rrLogin.Body.Bytes(), &respStruct)
	if err != nil {
		t.Fatal(err)
	}

	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + respStruct.AccessToken
	var jsonStr = []byte(`{"name":"TestMovieUpdatedName","description":"desc","typ":"fantasy"}`)

	req, err := http.NewRequest("PATCH", "/api/movie", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Authorization", bearer)
	q := req.URL.Query()
	q.Add("id", "27")
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.MovieUpdate)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		if status == http.StatusBadRequest {
			expected := `{"error": "Bad request"}
`
			if rr.Body.String() != expected {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), expected)
			}
		} else if status == http.StatusForbidden {
			expected := `{"error": "Movie with that email already exists"}
`
			if rr.Body.String() != expected {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), expected)
			}
		} else if status == http.StatusInternalServerError {
			expected := `{"error": "Internal server error"}
`
			if rr.Body.String() != expected {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), expected)
			}
		} else if status == http.StatusNotFound {
			expected := `{"error": "Movie with that id does not exist"}
`
			if rr.Body.String() != expected {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), expected)
			}
		} else {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
	} else {
		expected := fmt.Sprintf(`{"id":27,"name":"TestMovieUpdatedName","description":"desc","typ":"fantasy"}
`)
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), expected)
		}
	}

}
