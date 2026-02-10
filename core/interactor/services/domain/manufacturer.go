package domain


type Manufacturer struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (m *Manufacturer) ToLogger() []string {
	return []string{
		"id:" + m.ID,
		"name:" + m.Name,
	}
}
