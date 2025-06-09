package dto

type HealthResponse struct {
	Success   bool   `json:"success"`
	Status    string `json:"status"`
	Timestamp int64  `json:"timestamp"`
	Version   string `json:"version"`
}
