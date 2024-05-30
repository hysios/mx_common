package models

type CasbinRule struct {
	ID    uint   `gorm:"primaryKey;autoIncrement"`
	Ptype string `gorm:"size:128;uniqueIndex:unique_index"`
	V0    string `gorm:"size:256;uniqueIndex:unique_index"`
	V1    string `gorm:"size:128;uniqueIndex:unique_index"`
	V2    string `gorm:"size:128;uniqueIndex:unique_index"`
	V3    string `gorm:"size:128"`
	V4    string `gorm:"size:128"`
	V5    string `gorm:"size:128"`
}
