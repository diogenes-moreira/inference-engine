package inference

type Fact struct {
	ID          string
	Description string
	Values      map[string]interface{}
}
