package pkg

import "fmt"

type ErrorMissingDatabaseType struct {
}

func (ErrorMissingDatabaseType) Error() string {
	return "the environment variable `DATABASE` is missing"
}

type ErrorUnkownDatabaseType struct {
	Type string
}

func (e ErrorUnkownDatabaseType) Error() string {
	return fmt.Sprintf("unknown database type '%s'", e.Type)
}
