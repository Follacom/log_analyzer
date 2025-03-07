package model

// ApacheLogHeaders represents HTTP headers stored as JSON.
type ApacheLogHeaders struct {
	Host      string `log:"Host" json:"host"`             // %{Host}i
	Referer   string `log:"Referer" json:"referer"`       // %{Referer}i
	UserAgent string `log:"User-Agent" json:"user_agent"` // %{User-Agent}i
}
