package models_test

import (
	"telegram-assistant-bot/models"
	"testing"
)

func TestAddAdminsGroups(t *testing.T) {
	var group_id int64 = -565656
	//var admin_id int = 56564165
	maps := make(map[string]interface{})
	maps["group_id"] = group_id
	maps["admin_id"] = 5656
	err := models.AddAdminsGroups(maps)
	if err != nil {
		t.Error(err)
	}
}

func TestExistAdminsGroups(t *testing.T) {
	var group_id int64 = -565656
	var admin_id int = 56564165
	maps := models.AdminsGroups{
		GroupID: group_id,
		AdminID: admin_id,
	}
	_, err := models.ExistAdminsGroups(maps)
	if err != nil {
		t.Error(err)
	}
}
