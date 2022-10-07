package helper

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/burakkuru5534/src/auth"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var ExpireDuration = time.Minute * 1

var App *app

type app struct {
	DB *DbHandle
}

type conf struct {
	Auth      *auth.Auth
	JwtSecret string `json:"jwt_secret"`
}

var Conf *conf

func InitConf() {
	Conf = &conf{}
}

func (c *conf) SetAuth(auth *auth.Auth) {
	c.Auth = auth
}

func InitApp(db *DbHandle) error {
	App = &app{
		DB: db,
	}

	return nil
}

func BodyToJsonReq(r *http.Request, data interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return errors.New(fmt.Sprintf("Body unmarshall error %s", string(body)))
	}

	defer r.Body.Close()

	return nil
}

func StrToInt64(aval string) int64 {
	aval = strings.Trim(strings.TrimSpace(aval), "\n")
	i, err := strconv.ParseInt(aval, 10, 64)
	if err != nil {
		return 0
	}
	return i
}

func CheckIfShorteningVersionOfUrlExists(shorteningUrl string) (bool, error) {

	var isExists bool

	err := App.DB.Get(&isExists, "SELECT EXISTS (SELECT 1 FROM url WHERE shortening_version = $1)", shorteningUrl)
	if err != nil {
		return false, err
	}
	return isExists, nil

}

func GetUrlShorteningVersion(id int64) (string, error) {

	var name string

	err := App.DB.Get(&name, "SELECT shortening_version FROM url WHERE id = $1", id)
	if err != nil {
		return "", err
	}
	return name, nil

}

func CheckPass(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false
	}
	return true
}

func CheckPassBasic(dbPass string, password string) bool {

	if dbPass == password {
		return true
	} else {
		return false

	}
}

func HashPasswd(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), err
}

func SetLogJwt(jwt string, usrID int64) (int64, error) {
	var LogJwtID int64
	sq := `insert into logjwt (
			jwt
			, expires_on
			, usr_id
			, dttm				
		) values ( $1, $2, $3, current_timestamp) returning id`

	err := App.DB.QueryRow(sq, jwt, time.Now().Add(ExpireDuration), usrID).Scan(&LogJwtID)
	if err != nil {
		return 0, err
	}
	return LogJwtID, nil
}

func SetInvalidJwt(ctx context.Context, clID *int64) error {
	var err error

	sq := `update logjwt set is_invalid = true where cl_id = $1 and is_invalid = false`

	_, err = App.DB.ExecContext(ctx, sq, clID)
	if err != nil {
		return err
	}

	return nil
}

func getUsrShortenedUrlCount(usrID int64) (int64, error) {

	var count int64

	err := App.DB.Get(&count, "SELECT count(*) FROM url WHERE usr_id = $1", usrID)
	if err != nil {
		return 0, err
	}
	return count, nil

}

func getUsrAccountType(usrID int64) (string, error) {

	var accountType string

	err := App.DB.Get(&accountType, "SELECT account_type FROM usr WHERE id = $1", usrID)
	if err != nil {
		return "", err
	}
	return accountType, nil

}

func UsrUrlLimitCheck(usrID int64) error {

	count, err := getUsrShortenedUrlCount(usrID)
	if err != nil {
		return err
	}
	
	accountType, err := getUsrAccountType(usrID)
	if err != nil {
		return err
	}

	switch accountType {
	case "free":
		if count >= 1 {
			return errors.New("free account can shorten only 5 url")
		}

	case "premium":
		if count >= 10 {
			return errors.New("premium account can shorten only 100 url")
		}
	}

	return nil
}
