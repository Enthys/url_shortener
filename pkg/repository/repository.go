package repository

type LinkRepository interface {
	// StoreLink adds the link to the database under the provided id.
	//
	// If the database fails to store the link under the database a `errors.ErrorLinkSaveFailure` should be returned.
	//
	// Other types of errors could be returned based on the selected database type.
	StoreLink(id, link string) error

	// GetById returns the link which corresponds to the provided id if such a link exists.
	//
	// If the database fails to retrieve the link from the database due to some error a `errors.ErrorFailedRetrieval`
	// should be returned.
	//
	// If no record with the given id is found a `errors.ErrorLinkNotFound` errors should be returned.
	GetById(id string) (string, error)

	// GetLinkId returns the id of the provided link if such exists. This method is provided in order to reduce the.
	// amount of records in the database by reusing the id if the link is already present in the database.
	//
	// If the database fails to retrieve the link from the database due to some error a `errors.ErrorFailedRetrieval`
	// should be returned.
	//
	// If no record with the given id is found a `errors.ErrorIDNotFound` errors should be returned.
	GetLinkId(link string) (string, error)
}
