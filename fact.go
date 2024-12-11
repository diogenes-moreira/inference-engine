package inference

type Fact struct {
	ID          string      `json:"id"`
	Description string      `json:"description"`
	Value       interface{} `json:"values"`
	DerivedFrom []string    `json:"derived_from"`
}
