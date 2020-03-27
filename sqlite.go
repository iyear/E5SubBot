package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type MSData struct {
	tgId         int
	tgName       string
	refreshToken string
	uptime       int64
	other        string
}

func init() {
}
func AddData(db *sql.DB, u MSData) bool {
	sqlString := `
	INSERT INTO users (tg_id, tg_name, refresh_token, uptime,other)
	VALUES (?,?,?,?,?);`
	stmt, err := db.Prepare(sqlString)
	CheckErr(err)
	_, err = stmt.Exec(u.tgId, u.tgName, u.refreshToken, u.uptime, u.other)
	return CheckErr(err)
}
func CreateTB(db *sql.DB) {

	sqltable := `
    create table if not exists "users"
	(
	tg_id INTEGER,
	tg_name VARCHAR(255),
	refresh_token TEXT,
	uptime INTEGER,
	other TEXT
	);`
	db.Exec(sqltable)
}
