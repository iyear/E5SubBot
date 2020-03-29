package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"strings"
)

type MSData struct {
	tgId         int64
	refreshToken string
	msId         string
	uptime       int64
	other        string
}

func init() {
	var err error
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	CheckErr(err)
	host := viper.GetString("mysql.host")
	user := viper.GetString("mysql.user")
	port := viper.GetString("mysql.port")
	pwd := viper.GetString("mysql.user")
	database := viper.GetString("mysql.database")
	path := strings.Join([]string{user, ":", pwd, "@tcp(", host, ":", port, ")/", database, "?charset=utf8"}, "")
	//fmt.Println(path)
	db, err = sql.Open(dbDriverName, path)
	fmt.Println(db.Ping())
	CheckErr(err)
	CreateTB(db)
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
		var (
			tgIdt, uptimet          int64
			refresht, othert, msidt string
		)
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
		var tgIdt, uptimet int64
		rows.Scan(&tgIdt, &refresht, &msidt, &uptimet, &othert)
		//fmt.Println(string(tgNamet) + "=>" + uptimet.Format("2006-01-02 15:04:05"))
		result = append(result, MSData{tgIdt, refresht, msidt, uptimet, othert})
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
		var tgIdt, uptimet int64
		rows.Scan(&tgIdt, &refresht, &msidt, &uptimet, &othert)
		result = append(result, MSData{tgIdt, refresht, msidt, uptimet, othert})
	}
	return result
}
func CreateTB(db *sql.DB) (bool, error) {

	sqltable := `
    create table if not exists users
	(
	tg_id INTEGER,
	refresh_token TEXT,
	ms_id TEXT,
	uptime INTEGER,
	other TEXT
	);`
	_, err := db.Exec(sqltable)
	if err != nil {
		return false, err
	}
	return true, nil
}
