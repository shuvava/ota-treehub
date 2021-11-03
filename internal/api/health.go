package api

// Status of service
type Status string

const (
	// LivenessPath endpoint default path
	LivenessPath = "/healthz"

	// ReadinessPath endpoint default path
	ReadinessPath = "/readyz"

	// StatusHealthy means that service in health state
	StatusHealthy = Status("StatusHealthy")

	// StatusUnhealthy means that service is in degraded state
	StatusUnhealthy = Status("StatusUnhealthy")
)

// HealthEntryStatus is status of external dependency like db or queue
type HealthEntryStatus struct {
	Status   Status      `json:"status"`
	Data     interface{} `json:"data"`
	Resource string      `json:"resource"`
}

// HealthStatusResponse response of health endpoints
type HealthStatusResponse struct {
	Status  Status              `json:"status"`
	Entries []HealthEntryStatus `json:"entries,omitempty"`
}
