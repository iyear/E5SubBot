package main

import (
	"fmt"
	"log"
	"os"
)

func CheckErr(err error) bool {
	if err != nil {
		log.Println(err)
		fmt.Println("error: ", err.Error())
		panic(err)
		return false
	}
	return true
}
func FileExist(Path string) bool {
	if _, err := os.Stat(Path); err != nil {
		if os.IsNotExist(err) {
			return false
		} else {
			CheckErr(err)
		}
	}
	return true
}
