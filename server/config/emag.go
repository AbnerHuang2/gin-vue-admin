package config

type Emag struct {
	ApiUrl          string  `mapstructure:"api-url" json:"api-url" yaml:"api-url"`                            // Emag API URL
	Cookie          string  `mapstructure:"cookie" json:"cookie" yaml:"cookie"`                               // Emag Cookie（用于认证）
	RequestInterval float64 `mapstructure:"request-interval" json:"request-interval" yaml:"request-interval"` // 请求间隔（秒），默认 1.5
	RetryCount      int     `mapstructure:"retry-count" json:"retry-count" yaml:"retry-count"`                // 单次 API 调用重试次数，默认 3
	BatchSize       int     `mapstructure:"batch-size" json:"batch-size" yaml:"batch-size"`                   // 每批处理数量，默认 50
	SnapshotDayGap  int     `mapstructure:"snapshot-day-gap" json:"snapshot-day-gap" yaml:"snapshot-day-gap"` // 快照日期间隔天数，默认 15
	MaxFailCount    int     `mapstructure:"max-fail-count" json:"max-fail-count" yaml:"max-fail-count"`       // 最大连续失败次数，超过则标记为 bad_request，默认 3
}
