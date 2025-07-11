package entities

type IconMenu struct {
	IconActive    string `json:"iconActive"`
	IconNonActive string `json:"iconNonActive"`
}

type ChildMenu struct {
	Code     string `json:"code"`
	Name     string `json:"name"`
	Priority string `json:"priority"`
}

type Menus struct {
	Id          int    `gorm:"primary_key;Column:id" json:"id"`
	Code        string `gorm:"Column:code" json:"code"`
	Name        string `gorm:"Column:name" json:"name"`
	Path        string `gorm:"Column:path" json:"path"`
	Icon        string `gorm:"type:json;Column:icons" json:"icons"`
	HasChild    bool   `gorm:"has_child" json:"hasChild"`
	ChildMenus  string `gorm:"child_menus" json:"childMenus"`
	Priority    int    `gorm:"priority" json:"priority"`
	CanAccessBy string `gorm:"can_access_by" json:"-"`
}
