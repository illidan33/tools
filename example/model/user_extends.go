package model_test

import (
	"time"
)

// @def primary Id
// @def unique_index:uni_user_id UserId
type UserExtends struct {
	Id               uint64    `gorm:"column:id" json:"id"`                                 // id
	UserId           uint64    `gorm:"column:user_id" json:"user_id"`                       // users id
	IdType           uint8     `gorm:"column:id_type" json:"id_type"`                       // 1:ID card,2:passport
	IdNumber         string    `gorm:"column:id_number" json:"id_number"`                   // ID number
	NationalityId    uint64    `gorm:"column:nationality_id" json:"nationality_id"`         // user_nationality.id
	StateId          uint64    `gorm:"column:state_id" json:"state_id"`                     // user_state.id
	OccupationId     uint64    `gorm:"column:occupation_id" json:"occupation_id"`           // user_occupation.id
	NatureBusinessId uint64    `gorm:"column:nature_business_id" json:"nature_business_id"` // user_nature_business.id
	IsPassEkyc       uint8     `gorm:"column:is_pass_ekyc" json:"is_pass_ekyc"`             // is pass ekyc(0.no 1.yes)
	Latitude         string    `gorm:"column:latitude" json:"latitude"`                     // latitude
	Longitude        string    `gorm:"column:longitude" json:"longitude"`                   // longitude
	CreatedTime      time.Time `gorm:"column:created_time" json:"created_time"`             // create time
	UpdatedTime      time.Time `gorm:"column:updated_time" json:"updated_time"`             // update time
}
