package domain

// Manufacturer represents the manufacturer domain model
// Contains the manufacturer information for aircraft models
type Manufacturer struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ToLogger returns a slice of strings for logging manufacturer information
func (m *Manufacturer) ToLogger() []string {
	return []string{
		"id:" + m.ID,
		"name:" + m.Name,
	}
}
