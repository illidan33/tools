package entity

import "time"

// @def primary Uuid
// @def index:photo Photo
// @def index:user_profile_id UserProfileId
// @def foreign_index:user_profiles_ibfk_1 Photo
type UserProfiles struct {
	Uuid             string    `gorm:"column:uuid" json:"uuid"`
	Name             string    `gorm:"column:name" json:"name"`
	Email            string    `gorm:"column:email" json:"email"`
	OfficeEmail      string    `gorm:"column:office_email" json:"officeEmail"`
	PhoneCountryCode string    `gorm:"column:phone_country_code" json:"phoneCountryCode"`
	Phone            string    `gorm:"column:phone" json:"phone"`
	AltCountryCode   string    `gorm:"column:alt_country_code" json:"altCountryCode"`
	OfficePhone      string    `gorm:"column:office_phone" json:"officePhone"`
	PhotoUrl         string    `gorm:"column:photo_url" json:"photoUrl"`
	CreatedAt        time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt        time.Time `gorm:"column:updated_at" json:"updatedAt"`
	Photo            string    `gorm:"column:photo" json:"photo"`
	DeletedAt        time.Time `gorm:"column:deleted_at" json:"deletedAt"`
	FirstLoginAt     time.Time `gorm:"column:first_login_at" json:"firstLoginAt"`
	OldPhone         string    `gorm:"column:old_phone" json:"oldPhone"`
	IsGroupAdmin     int8      `gorm:"column:is_group_admin" json:"isGroupAdmin"`
	IsDemoAccount    int8      `gorm:"column:is_demo_account" json:"isDemoAccount"`
	QrCode           string    `gorm:"column:qr_code" json:"qrCode"`
	QrcodeExpiredAt  time.Time `gorm:"column:qrcode_expired_at" json:"qrcodeExpiredAt"`
	QrcodeUpdatedAt  time.Time `gorm:"column:qrcode_updated_at" json:"qrcodeUpdatedAt"`
	UserProfileId    int64     `gorm:"column:user_profile_id" json:"userProfileId"`
	ShouldMigrate    string    `gorm:"column:should_migrate" json:"shouldMigrate"`
	UmsTokenStatus   string    `gorm:"column:ums_token_status" json:"umsTokenStatus"`
	MigrationStatus  string    `gorm:"column:migration_status" json:"migrationStatus"`
	MigratedPhone    string    `gorm:"column:migrated_phone" json:"migratedPhone"`
	KbToken          string    `gorm:"column:kb_token" json:"kbToken"`
	KbRefreshToken   string    `gorm:"column:kb_refresh_token" json:"kbRefreshToken"`
}
