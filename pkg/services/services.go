package services

type LinkServiceInterface interface {
	CreateLinkId() (string, error)
	StoreLink(path, link string) (string, error)
	GetLinkFromId(path string) (string, error)
}
