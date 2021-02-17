package factors

type Factor struct {
	ID       int    `table:"factors" json:"id" column:"id"`
	UserID   int    `json:"userID" column:"user_id"`
	Factor   string `json:"factor" column:"factor" encrypted:"true"`
	Archived bool   `json:"archived" column:"archived"`
}
