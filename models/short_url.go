package models

import (
	"fit/config"
	"time"
)

// ShortUrl Short URL
type ShortUrl struct {
	ID          uint        `json:"id" gorm:"primarykey"`
	PremiumLink PremiumLink `json:"premium_link" gorm:"foreignKey:ShortUrlId;references:ID;constraint:OnDelete:CASCADE"`
	ShortKey    string      `json:"short_key" gorm:"unique;not null;type:varchar(50)"`
	Url         string      `json:"url" gorm:"not null;type:text"`
	HashedUrl   string      `json:"hashed_url" gorm:"unique;not null;type:text"`
	CreatedAt   time.Time   `json:"created_at"`
}

// TableName 테이블 이름
func (s *ShortUrl) TableName() string {
	return "short_urls"
}

// ShortId 임시 데이터
type ShortId struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 테이블 이름
func (s *ShortId) TableName() string {
	return "short_ids"
}

// PremiumLink 프리미엄 링크
type PremiumLink struct {
	ID              uint      `json:"id" gorm:"primarykey"`
	ShortUrlId      uint      `json:"short_url_id" gorm:"unique"`
	Email           string    `json:"email" gorm:"not null;type:varchar(100)"`
	AndroidDeeplink string    `json:"android_deeplink" gorm:"not null;type:varchar(1000)"`
	AndroidFallback string    `json:"android_fallback" gorm:"not null;type:varchar(2000)"`
	IOSDeeplink     string    `json:"ios_deeplink" gorm:"not null;type:varchar(1000)"`
	IOSFallback     string    `json:"ios_fallback" gorm:"not null;type:varchar(2000)"`
	WebHookUrl      string    `json:"web_hook_url" gorm:"not null;type:varchar(2000)"`
	OGTags          string    `json:"og_tags" gorm:"not null;type:text"`
	IsActive        bool      `json:"is_active" gorm:"not null;default:true"`
	IsVerified      bool      `json:"is_verified" gorm:"not null;default:false"`
	CreatedAt       time.Time `json:"created_at"`
}

func init() {
	writeDb := config.GetWriteDB()
	_ = writeDb.AutoMigrate(&ShortUrl{})
	_ = writeDb.AutoMigrate(&ShortId{})
	_ = writeDb.AutoMigrate(&PremiumLink{})
}
