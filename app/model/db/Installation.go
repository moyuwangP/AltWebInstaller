package db

import "time"

type Installation struct {
	UDID          string    `json:"udid" gorm:"column:udid;index:udid_md5,unique;index:bundle,unique"`
	MD5           string    `json:"md5" gorm:"index:udid_md5,unique"`
	BundleID      string    `json:"bundle_id" gorm:"index:bundle,unique"`
	BundleVersion string    `json:"bundle_version" gorm:"index:bundle,unique"`
	RemovePlugIns bool      `json:"remove_plug_ins"`
	RefreshedAt   time.Time `json:"refreshed_at"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
