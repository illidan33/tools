package model_test

import (
	"time"
)

// @def primary Uuid
// @def index:photo Photo
// @def index:user_profile_id UserProfileId
type UserProfiles struct {
	Uuid             string    `gorm:"column:uuid" json:"uuid"`
	Name             string    `gorm:"column:name" json:"name"`
	Email            string    `gorm:"column:email" json:"email"`
	OfficeEmail      string    `gorm:"column:office_email" json:"office_email"`
	PhoneCountryCode string    `gorm:"column:phone_country_code" json:"phone_country_code"`
	Phone            string    `gorm:"column:phone" json:"phone"`
	AltCountryCode   string    `gorm:"column:alt_country_code" json:"alt_country_code"`
	OfficePhone      string    `gorm:"column:office_phone" json:"office_phone"`
	PhotoUrl         string    `gorm:"column:photo_url" json:"photo_url"`
	CreatedAt        time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at" json:"updated_at"`
	Photo            string    `gorm:"column:photo" json:"photo"`
	DeletedAt        time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	FirstLoginAt     time.Time `gorm:"column:first_login_at" json:"first_login_at"`
	OldPhone         string    `gorm:"column:old_phone" json:"old_phone"`
	IsGroupAdmin     int8      `gorm:"column:is_group_admin" json:"is_group_admin"`
	IsDemoAccount    int8      `gorm:"column:is_demo_account" json:"is_demo_account"`
	QrCode           string    `gorm:"column:qr_code" json:"qr_code"`
	QrcodeExpiredAt  time.Time `gorm:"column:qrcode_expired_at" json:"qrcode_expired_at"`
	QrcodeUpdatedAt  time.Time `gorm:"column:qrcode_updated_at" json:"qrcode_updated_at"`
	UserProfileId    int64     `gorm:"column:user_profile_id" json:"user_profile_id"`
	ShouldMigrate    string    `gorm:"column:should_migrate" json:"should_migrate"`
	UmsTokenStatus   string    `gorm:"column:ums_token_status" json:"ums_token_status"`
	MigrationStatus  string    `gorm:"column:migration_status" json:"migration_status"`
	MigratedPhone    string    `gorm:"column:migrated_phone" json:"migrated_phone"`
	KbToken          string    `gorm:"column:kb_token" json:"kb_token"`
	KbRefreshToken   string    `gorm:"column:kb_refresh_token" json:"kb_refresh_token"`
}
