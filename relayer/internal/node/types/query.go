package types

type QueryNodesParams struct {
	Chain  string `json:"chain"`
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
}
