package domain


type Engine struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}


func (e *Engine) ToLogger() []string {
	return []string{
		"id:" + e.ID,
		"name:" + e.Name,
	}
}
