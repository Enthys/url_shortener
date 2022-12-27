package services

import (
	"errors"
	"fmt"

	"github.com/Enthys/url_shortener/pkg/repository"
)

type LinkService struct {
	repository repository.LinkRepository
}

func NewLinkService(repository repository.LinkRepository) LinkServiceInterface {
	return &LinkService{
		repository: repository,
	}
}

// CreateLinkId calls the repository to generate a unique id
func (l *LinkService) CreateLinkId() (string, error) {
	return l.repository.GenerateId()
}

// StoreLink accepts a link and an id with which to associate the link with. It first checks if the link is already
// stored in the data. If the link is not present in the database then it will instruct the repository to save the link
// under the provided id. If the link does already exist in the database then it will not attemt a rewrite of the link
// but instead it will return the already existing id of the link in the database.
//
// If the repository fails to fetch the existing link id due to an error with the query or the database a
// `errors.ErrorFailedToSaveLink` error will be returned.
func (l *LinkService) StoreLink(id, link string) (string, error) {
	existingId, err := l.repository.GetLinkId(link)
	if err != nil && !errors.Is(err, repository.ErrorIDNotFound{}) {
		return "", ErrorFailedToSaveLink{Err: err}
	}

	// The link is already in the storage
	if existingId != "" {
		return existingId, nil
	}

	return id, l.repository.StoreLink(id, link)
}

// GetLinkFromId retrieves the link identified by the given ID.
//
// If the repository fails to retireve the link due to a query or database error a `errors.ErrorFailedRetrieval` error
// will be returned.
//
// If the link does not exist a `errors.ErrorLinkNotFound` error will be returned.
func (l *LinkService) GetLinkFromId(id string) (string, error) {
	link, err := l.repository.GetById(id)

	switch err.(type) {
	case nil:
		return link, nil
	case repository.ErrorFailedRetrieval:
		return "", ErrorFailedRetrieval{}
	case repository.ErrorLinkNotFound:
		return "", ErrorLinkNotFound{}
	default:
		return "", fmt.Errorf("unknown error received. Error: %w", err)
	}
}
