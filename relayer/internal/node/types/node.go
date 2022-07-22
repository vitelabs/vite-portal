package types

type Node struct {
	Id            string `json:"id"`
	Chain         string `json:"chain"`
	IpAddress     string `json:"ipAddress"`
	RewardAddress string `json:"rewardAddress"`
}

func (n *Node) IsValid() bool {
	return n != nil && n.Id != "" && n.Chain != "" && n.IpAddress != ""
}

type GetNodesParams struct {
	Chain  string `json:"chain"`
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
}
