package models_test

import (
	"telegram-assistant-bot/models"
	"testing"
)

func TestAddAdminsGroups(t *testing.T) {
	var groupId int64 = -565656
	var adminId int = 56564165
	maps := make(map[string]interface{})
	maps["group_id"] = groupId
	maps["admin_id"] = adminId
	err := models.AddAdminsGroups(maps)
	if err != nil {
		t.Error(err)
	}
}

func TestExistAdminsGroups(t *testing.T) {
	var groupId int64 = -565656
	var adminId int = 56564165
	maps := make(map[string]interface{})
	maps["group_id"] = groupId
	maps["admin_id"] = adminId
	_, err := models.ExistAdminsGroups(maps)
	if err != nil {
		t.Error(err)
	}
}
