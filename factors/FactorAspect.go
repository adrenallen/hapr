package factors

type FactorAspect struct {
	ID           *int   `table:"factor_aspects" json:"id" column:"id"`
	FactorID     int    `json:"factorID" column:"factor_id"`
	FactorAspect string `json:"factorAspect" column:"factor_aspect" encrypted:"true"`
	Archived     bool   `json:"archived" column:"archived"`
}
