package annkesdk

import "fmt"

type AnnkeRestError struct {
	Status  int
	Message string
	Path    string
}

type AnnkeInitError struct {
	Parameter string
}

func NewAnnkeError(status int, message string, path string) AnnkeRestError {
	return AnnkeRestError{
		Status:  status,
		Message: message,
		Path:    path,
	}
}

func (ae AnnkeRestError) Error() string {
	return fmt.Sprintf("received unexpected response from: %s status: %d payload: %s", ae.Path, ae.Status, ae.Message)
}

func NewAnnkeInitError(parameter string) AnnkeInitError {
	return AnnkeInitError{
		Parameter: parameter,
	}
}

func (ae AnnkeInitError) Error() string {
	return fmt.Sprintf("error initializing Annke connection, missing parameter: %s", ae.Parameter)
}
