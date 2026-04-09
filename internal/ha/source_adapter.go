package ha

type SourceResult struct {
	Name     string                 `json:"source"`
	Type     string                 `json:"type"`
	EntityID string                 `json:"entity_id,omitempty"`
	Data     map[string]interface{} `json:"data"`
	Status   Status                 `json:"status"`
}

type Status struct {
	OK       bool     `json:"ok"`
	Warnings []string `json:"warnings"`
}
