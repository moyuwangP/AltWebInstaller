package db

import "time"

type Package struct {
	CFBundleName               string    `json:"cf_bundle_name"`
	CFBundleShortVersionString string    `json:"cf_bundle_short_version_string"`
	CFBundleIdentifier         string    `json:"cf_bundle_identifier"`
	ContainsPlugIn             bool      `json:"contains_plug_in"`
	MD5                        string    `json:"md5" gorm:"unique"`
	CreatedAt                  time.Time `json:"created_at"`
	UpdatedAt                  time.Time `json:"updated_at"`
}
