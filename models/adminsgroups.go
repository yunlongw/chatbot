package models

import "github.com/jinzhu/gorm"

type AdminsGroups struct {
	GroupID int64 `json:"group_id"`
	AdminID int64   `json:"admin_id"`
}

func ExistAdminsGroups(ags AdminsGroups) (bool, error) {
	var ag AdminsGroups
	maps := make(map[string]string)
	maps["group_id"] = string(ags.GroupID)
	maps["admin_id"] = string(ags.AdminID)
	err := db.Select("id").Where(maps).First(&ag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if ag.GroupID > 0 {
		return true, nil
	}
	return false, nil
}

func AddAdminsGroups(data map[string]interface{}) error  {
	d := AdminsGroups{
		GroupID: data["group_id"].(int64),
		AdminID: data["admin_id"].(int64),
	}
	if err := db.Create(&d).Error ; err != nil{
		return err
	}
	return nil
}

