package dto

const (
	AntiAbuseActionReject = "reject"
	AntiAbuseActionReview = "review"
)

// PublishRateLimit defines the maximum number of one content type that may be
// created in a rolling time window. A zero value disables that individual rule.
type PublishRateLimit struct {
	DurationMinutes int `json:"durationMinutes"`
	MaxCount        int `json:"maxCount"`
}

// PublishFrequencyConfig holds the three independently configured publishing
// limits for one principal (a user or an IP address).
type PublishFrequencyConfig struct {
	Enabled bool             `json:"enabled"`
	Action  string           `json:"action"`
	Topic   PublishRateLimit `json:"topic"`
	Article PublishRateLimit `json:"article"`
	Comment PublishRateLimit `json:"comment"`
}

// AntiAbuseConfig is intentionally small: it only configures user and IP
// publishing frequency. More detection strategies are not part of this
// setting model.
type AntiAbuseConfig struct {
	User PublishFrequencyConfig `json:"user"`
	IP   PublishFrequencyConfig `json:"ip"`
}

func DefaultAntiAbuseConfig() AntiAbuseConfig {
	return AntiAbuseConfig{
		User: PublishFrequencyConfig{
			Action:  AntiAbuseActionReject,
			Topic:   PublishRateLimit{DurationMinutes: 10, MaxCount: 1},
			Article: PublishRateLimit{DurationMinutes: 10, MaxCount: 1},
			Comment: PublishRateLimit{DurationMinutes: 10, MaxCount: 3},
		},
		IP: PublishFrequencyConfig{
			Action:  AntiAbuseActionReject,
			Topic:   PublishRateLimit{DurationMinutes: 10, MaxCount: 5},
			Article: PublishRateLimit{DurationMinutes: 10, MaxCount: 5},
			Comment: PublishRateLimit{DurationMinutes: 10, MaxCount: 12},
		},
	}
}
