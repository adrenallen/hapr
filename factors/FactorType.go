package factors

type FactorType struct {
	ID         int    `table:"factor_types" json:"id"`
	FactorType string `json:"factorType"`
}
