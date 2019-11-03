package models

import (
	"github.com/jinzhu/gorm"
)

type Group struct {
	Model
	GroupID int    `json:"group_id"`
	Title   string `json:"title"`
}

type GroupSetting struct {
	Group  Group  `json:"group"`
	Key    string `json:"key"`
	Values string `json:"values"`
}

func GetGroups(pageNum int, pageSize int, maps interface{}) ([]*Group, error) {
	var groups []*Group
	err := db.Where(maps).Offset(pageNum).Order("id desc").Limit(pageSize).Find(&groups).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return groups, nil
}

func ExistGroups(id int) (bool, error) {
	var group Group
	err := db.Select("id").Where("id=?", id).First(&group).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if group.ID > 0 {
		return true, nil
	}
	return false, nil
}

func AddGroup(data map[string]interface{}) error {
	group := Group{
		GroupID: data["group_id"].(int),
		Title:   data["title"].(string),
	}
	if err := db.Create(&group).Error; err != nil {
		return err
	}
	return nil
}
