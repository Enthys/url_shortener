package services

type LinkService interface {
	CreateLinkId() string
	StoreLink(path, link string) (string, error)
	GetLinkPath(link string) (string, error)
	GetLinkFromId(path string) (string, error)
}
