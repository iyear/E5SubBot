package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tidwall/gjson"
	"time"
)

type MSData struct {
	tgId         int64
	refreshToken string
	msId         string
	uptime       time.Time
	other        string
}

func init() {
}

//update data by msId
func UpdateData(db *sql.DB, u MSData) (bool, error) {
	sqlString := `UPDATE users set tg_id=?,refresh_token=?,uptime=?,other=?  where ms_id=?`
	stmt, err := db.Prepare(sqlString)
	if err != nil {
		return false, err
	}
	res, err := stmt.Exec(u.tgId, u.refreshToken, u.uptime, u.other, u.msId)
	if err != nil {
		return false, err
	}
	fmt.Println("Update Data Successd:", res)
	return true, nil
}

//add data
func AddData(db *sql.DB, u MSData) (bool, error) {
	sqlString := `
	INSERT INTO users (tg_id, refresh_token,ms_id, uptime,other)
	VALUES (?,?,?,?,?)`
	stmt, err := db.Prepare(sqlString)
	if err != nil {
		return false, err
	}
	_, err = stmt.Exec(u.tgId, u.refreshToken, u.msId, u.uptime, u.other)
	if err != nil {
		return false, err
	}
	return true, nil
}

//del data by ms_id
func DelData(db *sql.DB, msId string) (bool, error) {
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
func QueryDataByMS(db *sql.DB, msId string) []MSData {
	rows, err := db.Query("select  * from users where ms_id = ?", msId)
	CheckErr(err)
	var result = make([]MSData, 0)
	defer rows.Close()
	for rows.Next() {
		var refresht, othert, msidt string
		var tgIdt int64
		var uptimet time.Time
		rows.Scan(&tgIdt, &refresht, &msidt, &uptimet, &othert)
		//fmt.Println(string(tgNamet) + "=>" + uptimet.Format("2006-01-02 15:04:05"))
		result = append(result, MSData{tgIdt, refresht, msidt, uptimet, othert})
	}
	return result
}

func QueryDataAll(db *sql.DB) []MSData {
	rows, err := db.Query("select  * from users ")
	CheckErr(err)
	var result = make([]MSData, 0)
	defer rows.Close()
	for rows.Next() {
		var refresht, othert, msidt string
		var tgIdt int64
		var uptimet time.Time
		rows.Scan(&tgIdt, &refresht, &msidt, &uptimet, &othert)
		//fmt.Println(string(tgNamet) + "=>" + uptimet.Format("2006-01-02 15:04:05"))
		result = append(result, MSData{tgIdt, refresht, msidt, uptimet, othert})
	}
	return result
}
func QueryDataBySign(db *sql.DB, tgId int64, sign string) []MSData {
	rows, err := db.Query("select  * from users where tg_id = ?", tgId)
	CheckErr(err)
	var result = make([]MSData, 0)
	defer rows.Close()
	for rows.Next() {
		var refresht, othert, msidt string
		var tgIdt int64
		var uptimet time.Time
		rows.Scan(&tgIdt, &refresht, &msidt, &uptimet, &othert)
		if gjson.Get(othert, "sign").String() == sign {
			result = append(result, MSData{tgIdt, refresht, msidt, uptimet, othert})
		}
	}
	return result
}

//query data by tg_id
func QueryDataByTG(db *sql.DB, tgId int64) []MSData {
	rows, err := db.Query("select  * from users where tg_id = ?", tgId)
	CheckErr(err)
	var result = make([]MSData, 0)
	defer rows.Close()
	for rows.Next() {
		var refresht, othert, msidt string
		var tgIdt int64
		var uptimet time.Time
		rows.Scan(&tgIdt, &refresht, &msidt, &uptimet, &othert)
		result = append(result, MSData{tgIdt, refresht, msidt, uptimet, othert})
	}
	return result
}
func CreateTB(db *sql.DB) (bool, error) {

	sqltable := `
    create table if not exists "users"
	(
	tg_id INTEGER,
	refresh_token TEXT,
	ms_id TEXT,
	uptime DATE,
	other TEXT
	);`
	_, err := db.Exec(sqltable)
	if err != nil {
		return false, err
	}
	return true, nil
}
