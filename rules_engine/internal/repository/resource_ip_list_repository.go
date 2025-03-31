package repository

type ResourceIPListRepository interface {
	AttachIPList(resourceID, ipListID string) error
	DetachIPList(resourceID, ipListID string) error
}
