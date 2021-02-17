package notifications

type UserReminderSetting struct {
	ID           *int `table:"user_reminder_settings" json:"id" column:"id"`
	UserID       int  `json:"userID" column:"user_id"`
	ReminderTime int  `json:"reminderTime" column:"reminder_time"`
	Sunday       bool `json:"sunday" column:"sunday"`
	Monday       bool `json:"monday" column:"monday"`
	Tuesday      bool `json:"tuesday" column:"tuesday"`
	Wednesday    bool `json:"wednesday" column:"wednesday"`
	Thursday     bool `json:"thursday" column:"thursday"`
	Friday       bool `json:"friday" column:"friday"`
	Saturday     bool `json:"saturday" column:"saturday"`
	Email        bool `json:"email" column:"email"`
	Push         bool `json:"push" column:"push"`
	SMS          bool `json:"sms" column:"sms"`
}
