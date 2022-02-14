package main

import (
	"fmt"
	"nonoSheep/model"
	"nonoSheep/routers"
)

func main() {
	err := model.InitDb()
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("连接数据库成功！")
	}
	routers.Start()
}
