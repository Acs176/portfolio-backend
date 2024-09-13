package connect

import (
	"context"
	"localbe/experience"
	v1 "localbe/gen/experience/v1"

	"connectrpc.com/connect"
	"github.com/bufbuild/protovalidate-go"
)

type ExperienceService struct {
	validator                 *protovalidate.Validator
	createExperienceEntryFunc experience.CreateExperienceEntryFunc
}

func NewExperienceService(
	createExperienceEntryFunc experience.CreateExperienceEntryFunc,
) (*ExperienceService, error) {
	v, err := protovalidate.New()
	if err != nil {
		return nil, err
	}
	es := &ExperienceService{
		validator:                 v,
		createExperienceEntryFunc: createExperienceEntryFunc,
	}
	return es, nil
}

func (es *ExperienceService) CreateExperienceEntry(
	parentCtx context.Context,
	req *connect.Request[v1.CreateExperienceEntryRequest],
) (*connect.Response[v1.CreateExperienceEntryResponse], error) {
	err := es.validator.Validate(req.Msg)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	// convert req.Msg to domain object
	// call the repo function with the object
	// convert the object that the function returns to protobuf
	// construct a response and return it
}
