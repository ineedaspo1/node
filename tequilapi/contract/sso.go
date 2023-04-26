package contract

// MystnodesSSOLinkResponse contains a link to initiate auth via mystnodes
// swagger:model MystnodesSSOLinkResponse
type MystnodesSSOLinkResponse struct {
	Link string `json:"link"`
}
