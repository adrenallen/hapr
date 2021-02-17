package users

type FirebaseMessagingDevice struct {
	ID       *int   `table:"firebase_messaging_devices" json:"id" column:"id"`
	DeviceID string `json:"deviceID" column:"device_id"`
	UserID   int    `json:"userID" column:"user_id"`
}
