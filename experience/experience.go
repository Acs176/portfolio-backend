package experience

import "github.com/google/uuid"

type Experience struct {
	Id          uuid.UUID
	CompanyName string
	Position    string
	Start       string
	End         string
	Description string
}
