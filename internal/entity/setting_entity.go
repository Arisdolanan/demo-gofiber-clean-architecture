package entity

// Setting represents the settings database table
type Setting struct {
	ID           int64   `json:"id" db:"id"`
	SchoolID     int64   `json:"school_id" db:"school_id"`
	SettingKey   string  `json:"setting_key" db:"setting_key"`
	SettingValue *string `json:"setting_value" db:"setting_value"`
	GroupName    string  `json:"group_name" db:"group_name"`
	Description  *string `json:"description" db:"description"`
	BaseEntity
}

// SettingUpdateRequest represents the payload for updating settings
type SettingUpdateRequest struct {
	Settings []SettingItem `json:"settings" validate:"required,min=1"`
}

// SettingItem represents a single setting to update
type SettingItem struct {
	SchoolID     int64   `json:"school_id"`
	SettingKey   string  `json:"setting_key" validate:"required"`
	SettingValue *string `json:"setting_value"`
	GroupName    string  `json:"group_name" validate:"required"`
	Description  *string `json:"description"`
}

// SettingResponse represents the response payload for a setting
type SettingResponse struct {
	ID           int64   `json:"id"`
	SchoolID     int64   `json:"school_id"`
	SettingKey   string  `json:"setting_key"`
	SettingValue *string `json:"setting_value"`
	GroupName    string  `json:"group_name"`
	Description  *string `json:"description"`
}
