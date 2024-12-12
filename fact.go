package inference

type Fact struct {
	ID           string      `json:"id"`
	Description  string      `json:"description"`
	Value        interface{} `json:"value"`
	DerivedFrom  []string    `json:"derived_from"`
	Accumulative bool        `json:"accumulative"`
}

func (f *Fact) Equal(other *Fact) bool {
	return f.ID == other.ID && f.Value == other.Value
}
