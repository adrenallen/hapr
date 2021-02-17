package featureflags

type FeatureFlag struct {
	ID          int    `table:"feature_flags" json:"id" column:"id"`
	Description string `column:"description" json:"description"`
	Enabled     bool   `column:"enabled" json:"enabled"`
}
