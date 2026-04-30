package domain

// Config key constants used in the key-value config store.
const (
	ConfigKeyMasterFolder = "master_folder"
)

// AppConfig holds application-level settings.
type AppConfig struct {
	MasterFolder string `json:"masterFolder"`
}
