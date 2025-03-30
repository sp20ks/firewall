package entity

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type IPList struct {
	ID        string    `json:"id"`
	IP        net.IPNet `json:"ip"`
	ListType  string    `json:"list_type"`
	CreatorID string    `json:"creator_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (i IPList) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID        string    `json:"id"`
		IP        string    `json:"ip"`
		ListType  string    `json:"list_type"`
		CreatorID string    `json:"creator_id"`
		CreatedAt time.Time `json:"created_at"`
	}{
		ID:        i.ID,
		IP:        i.IP.String(),
		ListType:  i.ListType,
		CreatorID: i.CreatorID,
		CreatedAt: i.CreatedAt,
	})
}

func (i *IPList) UnmarshalJSON(data []byte) error {
	var alias struct {
		ID        string    `json:"id"`
		IP        string    `json:"ip"`
		ListType  string    `json:"list_type"`
		CreatorID string    `json:"creator_id"`
		CreatedAt time.Time `json:"created_at"`
	}
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}

	_, ipnet, err := net.ParseCIDR(alias.IP)
	if err != nil {
		return fmt.Errorf("invalid CIDR format: %s", alias.IP)
	}

	i.ID = alias.ID
	i.IP = *ipnet
	i.ListType = alias.ListType
	i.CreatorID = alias.CreatorID
	i.CreatedAt = alias.CreatedAt

	return nil
}
