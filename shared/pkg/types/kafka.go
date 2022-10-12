package types

type KafkaNodeUptimeStatus struct {
	EventId     string `json:"eventId"`
	Timestamp   string `json:"timestamp"`
	SuccessTime int64  `json:"successTime"`
	Round       int64  `json:"round"`
	ViteAddress string `json:"viteAddress"`
	NodeName    string `json:"nodeName"`
	Ip          string `json:"ip"`
}
