package health

type HealthStatus struct {
	Status string `json:"status"`
	Timestamp string `json:"timestamp"`
	Version string `json:"version"`
}