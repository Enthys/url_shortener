package services

import (
	"github.com/Enthys/url_shortener/pkg/repository"
	"github.com/xyproto/randomstring"
)

type linkService struct {
	repository repository.LinkRepository
}

func NewLinkService(repository repository.LinkRepository) LinkService {
	return &linkService{
		repository: repository,
	}
}

func (l *linkService) CreateLinkId() string {
	return randomstring.EnglishFrequencyString(16)
}

func (l *linkService) StoreLink(id, link string) (string, error) {
	existingId, err := l.repository.GetLinkId(link)
	if err != nil {
		switch err.(type) {
		case repository.ErrorPathNotFound:
			break
		}
	}

	// The link is already in the storage
	if existingId != "" {
		return existingId, nil
	}

	return id, l.repository.StoreLink(id, link)
}

func (l *linkService) GetLinkPath(link string) (string, error) {
	path, err := l.repository.GetLinkId(link)
	if err != nil {
		return "", err
	}

	return path, nil
}

func (l *linkService) GetLinkFromId(id string) (string, error) {
	link, err := l.repository.GetById(id)
	if err != nil {
		return "", err
	}

	return link, nil
}
