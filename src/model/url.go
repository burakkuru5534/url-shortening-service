package model

import (
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
