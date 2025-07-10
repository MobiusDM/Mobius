package contract

import "github.com/notawar/mobius/backend/server/mobius"

type ScimDetailsResponse struct {
	mobius.ScimDetails
	Err error `json:"-"`
}

func (r ScimDetailsResponse) Error() error {
	return r.Err
}
