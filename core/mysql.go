package core

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"main/logger"
	"main/outlook"
	"main/util"
	"time"
)

var db *gorm.DB

func InitDB() error {
	var err error
	db, err = gorm.Open("11", &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return err
	}
	return nil
}

//update data by msId
func UpdateData(u Client) (bool, error) {
	db, err := sql.Open(dbDriverName, dbPath)
	if err != nil {
		logger.Println(err)
	}
	defer db.Close()
	sqlString := `UPDATE Clients set tg_id=?,refresh_token=?,uptime=?,alias=?,client_id=?,client_secret=?,other=?  where ms_id=?`
	stmt, err := db.Prepare(sqlString)
	if err != nil {
		return false, err
	}
	_, err = stmt.Exec(u.TgId, u.RefreshToken, u.Uptime, u.Alias, u.ClientId, u.ClientSecret, u.Other, u.MsId)
	if err != nil {
		return false, err
	}
	return true, nil
}

//add data
func AddData(u Client) (bool, error) {
	db, err := sql.Open(dbDriverName, dbPath)
	if err != nil {
		logger.Println(err)
	}
	defer db.Close()
	sqlString := `
	INSERT INTO Clients (tg_id, refresh_token,ms_id, uptime,alias,client_id,client_secret,other)
	VALUES (?,?,?,?,?,?,?,?)`
	stmt, err := db.Prepare(sqlString)
	if err != nil {
		return false, err
	}
	_, err = stmt.Exec(u.TgId, u.RefreshToken, u.MsId, u.Uptime, u.Alias, u.ClientId, u.ClientSecret, u.Other)
	if err != nil {
		return false, err
	}
	return true, nil
}

//del data by ms_id
func DelData(msId string) (bool, error) {
	db, err := sql.Open(dbDriverName, dbPath)
	if err != nil {
		logger.Println(err)
	}
	defer db.Close()
	sqlString := `delete from Clients where ms_id=?`
	stmt, err := db.Prepare(sqlString)
	if err != nil {
		return false, err
	}
	res, err := stmt.Exec(msId)
	if err != nil {
		return false, err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return false, err
	}
	return true, nil
}
func QueryData(rows *sql.Rows) []Client {

	var result = make([]Client, 0)
	defer rows.Close()
	for rows.Next() {
		var (
			tgIdt, uptimet                                        int64
			refresht, othert, msidt, aliast, clientIdt, clientSet string
		)
		rows.Scan(&tgIdt, &refresht, &msidt, &uptimet, &aliast, &clientIdt, &clientSet, &othert)
		//fmt.Println(string(tgNamet) + "=>" + uptimet.Format("2006-01-02 15:04:05"))
		result = append(result, Client{tgIdt, refresht, msidt, uptimet, aliast, clientIdt, clientSet, othert})
	}
	return result
}
func QueryDataByMS(msId string) []Client {
	db, err := sql.Open(dbDriverName, dbPath)
	if err != nil {
		logger.Println(err)
	}
	defer db.Close()
	rows, err := db.Query("select  * from Clients where ms_id = ?", msId)
	util.CheckErr(err)
	return QueryData(rows)
}

func QueryDataAll() []Client {
	db, err := sql.Open(dbDriverName, dbPath)
	if err != nil {
		logger.Println(err)
	}
	defer db.Close()
	rows, err := db.Query("select  * from Clients ")
	util.CheckErr(err)
	return QueryData(rows)
}

//query data by tg_id
func QueryDataByTG(tgId int64) []Client {
	db, err := sql.Open(dbDriverName, dbPath)
	if err != nil {
		logger.Println(err)
	}
	defer db.Close()
	rows, err := db.Query("select  * from Clients where tg_id = ?", tgId)
	util.CheckErr(err)
	return QueryData(rows)
}
func CreateTB() (bool, error) {
	db, err := sql.Open(dbDriverName, dbPath)
	if err != nil {
		logger.Println(err)
	}
	defer db.Close()
	sqltable := `
    create table if not exists Clients
	(
	tg_id INTEGER,
	refresh_token TEXT,
	ms_id VARCHAR(255),
	uptime INTEGER,
	alias VARCHAR(255),
	client_id VARCHAR(255),
	client_secret VARCHAR(255),
	other TEXT
	);`
	_, err = db.Exec(sqltable)
	if err != nil {
		return false, err
	}
	return true, nil
}
