package model

import (
	"fmt"
	"nonoSheep/util/errormsg"
)

type Topic struct {
	Id           int    `json:"id"`
	LogoUrl      string `json:"logoUrl"`
	TopicName    string `json:"topicName"`
	Introduction string `json:"introduction"`
}

//获取所有主题
func GetTopics() ([]Topic, int) {
	var len int
	r, _ := Db.Query("select count(*)from topics")
	for r.Next() {
		_ = r.Scan(&len)

	}
	var topics = make([]Topic, len)
	rows, _ := Db.Query("select id,logoUrl,topicName,introduction from topics")
	if rows != nil {
		var i int
		for rows.Next() {
			err := rows.Scan(&topics[i].Id, &topics[i].LogoUrl, &topics[i].TopicName, &topics[i].Introduction)
			i++
			if err != nil {
				return nil, errormsg.ERROR
			}
		}
	}

	return topics, errormsg.SUCCESS
}

//创建主题
func CreatTopic(data *Topic) int {
	_, err := Db.Exec("insert into topics (logoUrl,topicName,introduction) values (?,?,?)", data.LogoUrl, data.TopicName, data.Introduction)
	if err != nil {
		fmt.Printf("err is:%v\n", err)
		return errormsg.ERROR
	}
	return errormsg.SUCCESS
}
