package model

import (
	"errors"
	"github.com/burakkuru5534/src/helper"
	"github.com/burakkuru5534/src/shortener"
)

type Url struct {
	ID               int    `db:"id"`
	LongVersion      string `db:"long_version"`
	ShortenedVersion string `db:"shortened_version"`
	UsrID            int64  `db:"usr_id"`
}

func (u *Url) Create(usrID int64) error {

	u.UsrID = usrID
	u.ShortenedVersion = shortener.GenerateShortLink(u.LongVersion, u.UsrID)

	_, err := helper.App.DB.Exec("INSERT INTO url (long_version, shortened_version, usr_id, zlins_dttm) VALUES ($1, $2, $3, current_timestamp)", u.LongVersion, u.ShortenedVersion, u.UsrID)
	return err
}

func (u *Url) List(usrID int64) ([]Url, error) {

	var urls []Url
	err := helper.App.DB.Select(&urls, "SELECT id,long_version,shortened_version, usr_id FROM url WHERE usr_id = $1", usrID)
	return urls, err
}

func (u *Url) Delete(usrID, id int64) error {

	isExists, err := checkIfUrlExists(id)
	if err != nil {
		return err
	}

	if !isExists {
		return errors.New("Url not found")
	}
	_, err = helper.App.DB.Exec("DELETE FROM url WHERE usr_id = $1 AND id = $2", usrID, id)
	return err
}

func checkIfUrlExists(id int64) (bool, error) {

	var isExists bool

	err := helper.App.DB.Get(&isExists, "SELECT EXISTS (SELECT 1 FROM url WHERE id = $1)", id)
	if err != nil {
		return false, err
	}
	return isExists, nil
}

func (u *Url) GetLongUrlFromShortened(usrID int64, shortenedUrl string) (string, error) {

	var longUrl string

	err := helper.App.DB.Get(&longUrl, "SELECT long_version FROM url WHERE shortened_version = $1 and usr_id = $2", shortenedUrl, usrID)
	if err != nil {
		return "", err
	}
	return longUrl, nil
}
