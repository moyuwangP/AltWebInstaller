package model

type InstallConfig struct {
	RemoveExtension bool   `json:"remove_extension"`
	Account         string `json:"account"`
}

type Config struct {
	Port          string `json:"port" validate:"required"`
	AppleId       string `json:"apple_id" validate:"required"`
	Password      string `json:"password" validate:"required"`
	AltserverPath string `json:"altserver_path" validate:"required"`
	AnisetteUrl   string `json:"anisette_url"`
	AutoRefresh   string `json:"auto_refresh"`
}
