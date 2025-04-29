package delivery

import "errors"

func errMissingFields() error {
	return errors.New("required fields are missing")
}

func errMissingID() error {
	return errors.New("id must be provided")
}
