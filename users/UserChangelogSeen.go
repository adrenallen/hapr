package users

type UserChangelogSeen struct {
	ID            int    `table:"user_changelog_seen" json:"id"`
	UserID        int    `json:"userID"`
	VersionString string `json:"versionString"`
}
