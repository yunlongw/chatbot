package models

import "github.com/jinzhu/gorm"

type AdminsGroups struct {
	GroupID int64 `json:"group_id"`
	AdminID int   `json:"admin_id"`
}

func ExistAdminsGroups(maps map[string]interface{}) (bool, error) {
	var adminsGroups AdminsGroups
	err := db.Where(maps).First(&adminsGroups).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if adminsGroups.GroupID > 0 {
		return true, nil
	}
	return false, nil
}

func AddAdminsGroups(data map[string]interface{}) error  {
	d := AdminsGroups{
		GroupID: data["group_id"].(int64),
		AdminID: data["admin_id"].(int),
	}
	if err := db.Create(&d).Error ; err != nil{
		return err
	}
	return nil
}

