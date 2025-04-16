package entity

type ScanResult struct {
	Action       Action `json:"action"`
	ModifiedURL  string `json:"modified_url,omitempty"`
	ModifiedBody string `json:"modified_body,omitempty"`
	Reason       string `json:"reason"`
}
