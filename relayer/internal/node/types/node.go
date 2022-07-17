package types

type Node struct {
	Id            string `json:"id"`
	Chain         string `json:"chain"`
	IpAddress     string `json:"ipAddress"`
	RewardAddress string `json:"rewardAddress"`
}
