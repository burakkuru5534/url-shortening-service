package helper

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "tayitkan"
	dbname   = "soft-robotics"
)

func ConnectDb() *sql.DB {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		panic(err)
	}
	// close database

	// check db
	err = db.Ping()

	fmt.Println("Connected!")

	return db
}

type PgConnectionInfo struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
	SSLMode  string
}

type queueStatement struct {
	sql          string
	params       []interface{}
	isNamedExec  bool
	namedExecArg interface{}
	// db *DbHandle
}

type DbHandle struct {
	*sqlx.DB
	LimitOfset string

	writeQueue chan queueStatement

	dateFormat           string
	dateTimeFormatMinute string
	dateTimeFormatSecond string
	timeFormat           string
}

//NewPgSqlx returns new DB connection for postgreSql with sqlx driver
func NewPgSqlx(conInfo PgConnectionInfo) (*sqlx.DB, error) {
	// todo: https://github.com/jackc/pgx/issues/81#issuecomment-296446179 ile aşağıdaki yöntemi etkinlik / performans açısından karşılaştır
	conninfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		conInfo.Host,
		conInfo.Port,
		conInfo.Username,
		conInfo.Password,
		conInfo.Database,
		conInfo.SSLMode)

	db, err := sqlx.Open("postgres", conninfo)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

//NewPgSqlxDbHandle returns DBHandle for PostgreSql Database with sqlx driver
func NewPgSqlxDbHandle(conInfo PgConnectionInfo, queueSize int) (*DbHandle, error) {
	db, err := NewPgSqlx(conInfo)
	if err != nil {
		return nil, err
	}

	dbHandleInstance := &DbHandle{
		DB:                   db,
		LimitOfset:           "offset %d limit %d",
		writeQueue:           make(chan queueStatement, queueSize),
		dateFormat:           "02.01.2006",
		dateTimeFormatMinute: "02.01.2006 15:04",
		dateTimeFormatSecond: "02.01.2006 15:04:05",
		timeFormat:           "15:04",
	}

	dbHandleInstance.dbqWriter()
	return dbHandleInstance, nil
}

func (db *DbHandle) dbqWriter() {
	go func() {
		for {
			select {
			case req := <-db.writeQueue:
				var err error
				if req.isNamedExec {
					_, err = db.NamedExec(req.sql, req.namedExecArg)
				} else {
					_, err = db.Exec(req.sql, req.params...)
				}
				if err != nil {
					errors.New("dbq")
				}
			}
		}
	}()
}
