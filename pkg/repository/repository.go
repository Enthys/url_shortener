package repository

type LinkRepository interface {
	StoreLink(id, link string) error
	GetById(id string) (string, error)
	GetLinkId(link string) (string, error)
}
