package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type MSData struct {
	tgId         int64
	refreshToken string
	msId         string
	uptime       int64
	alias        string
	clientId     string
	clientSecret string
	other        string
}

//update data by msId
func UpdateData(u MSData) (bool, error) {
	db, err := sql.Open(dbDriverName, dbPath)
	if err != nil {
		logger.Println(err)
	}
	defer db.Close()
	sqlString := `UPDATE users set tg_id=?,refresh_token=?,uptime=?,alias=?,client_id=?,client_secret=?,other=?  where ms_id=?`
	stmt, err := db.Prepare(sqlString)
	if err != nil {
		return false, err
	}
	_, err = stmt.Exec(u.tgId, u.refreshToken, u.uptime, u.alias, u.clientId, u.clientSecret, u.other, u.msId)
	if err != nil {
		return false, err
	}
	return true, nil
}

//add data
func AddData(u MSData) (bool, error) {
	db, err := sql.Open(dbDriverName, dbPath)
	if err != nil {
		logger.Println(err)
	}
	defer db.Close()
	sqlString := `
	INSERT INTO users (tg_id, refresh_token,ms_id, uptime,alias,client_id,client_secret,other)
	VALUES (?,?,?,?,?,?,?,?)`
	stmt, err := db.Prepare(sqlString)
	if err != nil {
		return false, err
	}
	_, err = stmt.Exec(u.tgId, u.refreshToken, u.msId, u.uptime, u.alias, u.clientId, u.clientSecret, u.other)
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
	sqlString := `delete from users where ms_id=?`
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
func QueryData(rows *sql.Rows) []MSData {

	var result = make([]MSData, 0)
	defer rows.Close()
	for rows.Next() {
		var (
			tgIdt, uptimet                                        int64
			refresht, othert, msidt, aliast, clientIdt, clientSet string
		)
		rows.Scan(&tgIdt, &refresht, &msidt, &uptimet, &aliast, &clientIdt, &clientSet, &othert)
		//fmt.Println(string(tgNamet) + "=>" + uptimet.Format("2006-01-02 15:04:05"))
		result = append(result, MSData{tgIdt, refresht, msidt, uptimet, aliast, clientIdt, clientSet, othert})
	}
	return result
}
func QueryDataByMS(msId string) []MSData {
	db, err := sql.Open(dbDriverName, dbPath)
	if err != nil {
		logger.Println(err)
	}
	defer db.Close()
	rows, err := db.Query("select  * from users where ms_id = ?", msId)
	CheckErr(err)
	return QueryData(rows)
}

func QueryDataAll() []MSData {
	db, err := sql.Open(dbDriverName, dbPath)
	if err != nil {
		logger.Println(err)
	}
	defer db.Close()
	rows, err := db.Query("select  * from users ")
	CheckErr(err)
	return QueryData(rows)
}

//query data by tg_id
func QueryDataByTG(tgId int64) []MSData {
	db, err := sql.Open(dbDriverName, dbPath)
	if err != nil {
		logger.Println(err)
	}
	defer db.Close()
	rows, err := db.Query("select  * from users where tg_id = ?", tgId)
	CheckErr(err)
	return QueryData(rows)
}
func CreateTB() (bool, error) {
	db, err := sql.Open(dbDriverName, dbPath)
	if err != nil {
		logger.Println(err)
	}
	defer db.Close()
	sqltable := `
    create table if not exists users
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
