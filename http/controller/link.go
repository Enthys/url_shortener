package controller

import (
	"fmt"
	"github.com/Enthys/url_shortener/pkg/repository"
	"github.com/Enthys/url_shortener/pkg/services"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"strings"
)

type LinkController struct {
	linkService services.LinkService
}

func NewLinkController(linkService services.LinkService) *LinkController {
	return &LinkController{
		linkService: linkService,
	}
}

// SetupRoutes adds the routes which the controller will use to the provided `echo.Echo` object
func (l *LinkController) SetupRoutes(router *echo.Echo) {
	router.GET("/", l.HomePage)
	router.GET("/:id", l.RedirectToLink)
	router.POST("/", l.CreateLink)
}

// HomePage displays the home page of the application. Given that the application will only be serving a single page
// the decision was made to manually retrieve the HTML from the index file and provide it instead of using a renderer
func (l *LinkController) HomePage(c echo.Context) error {
	workdir, _ := os.Getwd()
	html, err := os.ReadFile(fmt.Sprintf("%s/http/templates/index.html", workdir))
	if err != nil {
		c.Logger().Warnf("failed to retrieve template for home page. Error: %s", err)
		return err
	}

	if err = c.HTML(http.StatusOK, string(html)); err != nil {
		c.Logger().Warnf("could not send template. Error: %s", err)
	}

	return err
}

// RedirectToLink accepts redirect the user to the appropriate link if the id matches a link stored in the database
func (l *LinkController) RedirectToLink(c echo.Context) error {
	linkId := c.Param("id")
	link, err := l.linkService.GetLinkFromId(linkId)

	if err != nil {
		c.Logger().Warnf("Retrieval of link failed. Error: %s", err)
		switch err.(type) {
		case repository.ErrorLinkNotFound:
			if err = c.String(http.StatusNotFound, "No such link found"); err != nil {
				c.Logger().Warnf("Failed to return response for error. Error: %s", err.Error())
			}
			return err
		default:
			if err = c.String(http.StatusInternalServerError, "Issue is on our end"); err != nil {
				c.Logger().Warnf("Failed to return service issue. Error: %s", err)
			}
			return err
		}
	}

	if err = c.Redirect(http.StatusFound, link); err != nil {
		c.Logger().Warnf("Failed to redirect client to found link. Error: %s", err)
		return nil
	}

	return nil
}

type CreateLinkDTO struct {
	Link string `json:"link"`
}

type LinkCreatedDTO struct {
	ID string `json:"id"`
}

type BadRequestDTO struct {
	Error string `json:"error"`
}

// CreateLink accepts requests for creating a new shortened link for given link. The provided link has to begin with
// either `http://` or `https://` due to the way echo handles its redirects.
func (l *LinkController) CreateLink(c echo.Context) error {
	var dto CreateLinkDTO
	if err := c.Bind(&dto); err != nil {
		return err
	}

	if !strings.HasPrefix(dto.Link, "http://") && !strings.HasPrefix(dto.Link, "https://") {
		return c.JSON(
			http.StatusBadRequest,
			BadRequestDTO{Error: "the provided link should begin with either 'http://' or 'https://'"},
		)
	}

	newId, err := l.linkService.StoreLink(l.linkService.CreateLinkId(), dto.Link)
	if err != nil {
		return err
	}

	if err = c.JSON(http.StatusOK, LinkCreatedDTO{ID: newId}); err != nil {
		return err
	}

	return nil
}
