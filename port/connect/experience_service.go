package connect

import (
	"context"
	"localbe/experience"
	v1 "localbe/gen/experience/v1"
	"localbe/port/connect/converter"

	"connectrpc.com/connect"
	"github.com/bufbuild/protovalidate-go"
	"github.com/google/uuid"
)

type ExperienceService struct {
	validator                 *protovalidate.Validator
	createExperienceEntryFunc experience.CreateExperienceEntryFunc
	getExperienceEntryFunc    experience.GetExperienceEntryFunc
	getExperienceFunc         experience.GetExperienceFunc
}

func NewExperienceService(
	createExperienceEntryFunc experience.CreateExperienceEntryFunc,
	getExperienceEntryFunc experience.GetExperienceEntryFunc,
	getExperienceFunc experience.GetExperienceFunc,
) (*ExperienceService, error) {
	v, err := protovalidate.New()
	if err != nil {
		return nil, err
	}
	es := &ExperienceService{
		validator:                 v,
		createExperienceEntryFunc: createExperienceEntryFunc,
		getExperienceEntryFunc:    getExperienceEntryFunc,
		getExperienceFunc:         getExperienceFunc,
	}
	return es, nil
}

func (es *ExperienceService) CreateExperienceEntry(
	parentCtx context.Context,
	req *connect.Request[v1.CreateExperienceEntryRequest],
) (*connect.Response[v1.CreateExperienceEntryResponse], error) {
	err := es.validator.Validate(req.Msg.Experience)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	// convert req.Msg to domain object
	domExp, err := converter.ExperiencePbToDomain(req.Msg)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	// call the repo function with the object
	createdExp, err := es.createExperienceEntryFunc(
		parentCtx,
		domExp.CompanyName,
		domExp.Position,
		domExp.PeriodStart,
		domExp.PeriodEnd,
		domExp.RoleDescription,
	)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	// convert the object that the function returns to protobuf
	createdPb, err := converter.ExperienceToPb(createdExp)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	// construct a response and return it
	res := &v1.CreateExperienceEntryResponse{Experience: createdPb}
	return connect.NewResponse(res), nil
}

func (es *ExperienceService) GetExperienceEntry(
	parentCtx context.Context,
	req *connect.Request[v1.GetExperienceEntryRequest],
) (*connect.Response[v1.GetExperienceEntryResponse], error) {
	err := es.validator.Validate(req.Msg)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	convertedUuid, err := uuid.Parse(req.Msg.GetId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	domExp, err := es.getExperienceEntryFunc(parentCtx, convertedUuid)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	//convert to pb
	pbExp, err := converter.ExperienceToPb(domExp)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	res := &v1.GetExperienceEntryResponse{
		Experience: pbExp,
	}
	return connect.NewResponse(res), nil
}

func (es *ExperienceService) GetExperience(
	parentCtx context.Context,
	req *connect.Request[v1.GetExperienceRequest],
) (*connect.Response[v1.GetExperienceResponse], error) {

	domList, err := es.getExperienceFunc(parentCtx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	//convert to pb
	pbExp, err := converter.ExperiencesToPb(domList)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	res := &v1.GetExperienceResponse{
		ExperienceList: pbExp,
	}
	return connect.NewResponse(res), nil
}
