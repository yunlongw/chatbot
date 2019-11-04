package models

import (
	"github.com/jinzhu/gorm"
	"log"
)

type Group struct {
	Model
	GroupID int64  `json:"group_id"`
	Title   string `json:"title"`
}

type GroupSetting struct {
	Model
	Group   Group  `json:"group"`
	GroupID int64  `json:"group_id"`
	Key     string `json:"key"`
	Values  string `json:"values"`
}

func SetGroupSetting(groupId int64, key string, val string) (bool, error) {
	var groupSetting GroupSetting
	maps := make(map[string]interface{})
	maps["group_id"] = groupId
	maps["key"] = key
	groupSetting, err := FindGroupSetting(maps)
	if err != nil {
		log.Println(err)
		return false, err
	}

	if groupSetting.Key != "" {
		err := UpdateGroupSetting(maps, val)
		if err != nil {
			return false, err
		}
		return true, nil
	} else {
		err := AddGroupSetting(groupId, key, val)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

func UpdateGroupSetting(maps map[string]interface{}, val string) error {
	if err := db.Model(&GroupSetting{}).Where(maps).Update("values", val).Error; err != nil {
		return err
	}
	return nil
}

func AddGroupSetting(groupId int64, key string, val interface{}) error {
	groupSetting := &GroupSetting{
		GroupID: groupId,
		Key:     key,
		Values:  val.(string),
	}
	err := db.Create(&groupSetting).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	return nil
}

func FindGroupSetting(maps map[string]interface{}) (GroupSetting, error) {
	var groupSetting GroupSetting
	err := db.Where(maps).First(&groupSetting).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return groupSetting, err
	}
	return groupSetting, nil
}

func GetGroupSettings(maps map[string]interface{}) ([]GroupSetting, error) {
	var groupSetting []GroupSetting
	err := db.Where(maps).Find(&groupSetting).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return groupSetting, err
	}
	return groupSetting, nil
}

func GetGroupSettingByGroupIDToHashMap(GroupId int64) (map[string]interface{} , error) {
	maps := make(map[string]interface{})
	sets := make(map[string]interface{})
	maps["group_id"] = GroupId
	groupSetting, err := GetGroupSettings(maps)
	if err != nil {
		return sets, err
	}
	for _, v := range groupSetting {
		sets[v.Key] = v.Values
	}
	return sets, nil
}

func ExistGroupSetting(groupId int64, key string) (bool, error) {
	var groupSetting *GroupSetting
	maps := make(map[string]interface{})
	maps["group_id"] = groupId
	maps["key"] = key
	err := db.Where(maps).First(&groupSetting).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if groupSetting.GroupID > 0 {
		return true, nil
	}
	return false, nil
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

func ExistGroupsByGroupId(GroupId int64) (bool, error) {
	var group Group
	err := db.Select("id").Where("group_id=?", GroupId).First(&group).Error
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
		GroupID: data["group_id"].(int64),
		Title:   data["title"].(string),
	}
	if err := db.Create(&group).Error; err != nil {
		return err
	}
	return nil
}

func GetTotalGroup(maps interface{}) (int, error) {
	var total int
	if err := db.Model(&Group{}).Where(maps).Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}
