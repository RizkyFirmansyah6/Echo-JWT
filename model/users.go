package model

import (
	"EchoAPI/Helper"
	"fmt"
)

type User struct {
	ID       string `json:"id" form:"id"`
	Username string `json:"username" form:"username" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
	Nama     string `json:"nama" form:"nama"`
	Foto     string `json:"foto" form:"foto"`
}

func GetUsers() []User {
	db, err := Helper.OpenConn()
	if err != nil {
		fmt.Println(err)
	}
	var user []User
	db.Find(&user)
	return user
}

func PostUsers(username, password, nama, foto string) {
	db, err := Helper.OpenConn()
	if err != nil {
		fmt.Println(err)
	}
	user := User{Username: username, Password: password, Nama: nama, Foto: foto}
	db.Create(&user)
}

func PutUser(id, username, password, nama, foto string) {
	db, err := Helper.OpenConn()
	if err != nil {
		fmt.Println(err)
	}
	db.Model(User{ID: id}).Updates(User{Username: username, Password: password, Nama: nama, Foto: foto})
}

func DelUser(id string) {
	db, err := Helper.OpenConn()
	if err != nil {
		fmt.Println(err)
	}
	db.Delete(User{ID: id})
}

func Login(username, password string) User {
	db, err := Helper.OpenConn()
	if err != nil {
		fmt.Println(err)
	}
	var user User
	db.Where("username = ? AND password = ?", username, password).Find(&user)
	return user
}
