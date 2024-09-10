package configuration

import "errors"

var (
	ErrorLoadingConfiguration     = errors.New("error loading configuration")
	ErrorValidatingPostgresConfig = errors.New("error validating postgres configuration")
)
