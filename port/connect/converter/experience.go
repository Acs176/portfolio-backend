package converter

import (
	"errors"
	"localbe/experience"
	v1 "localbe/gen/experience/v1"
)

func ExperiencePbToDomain(pb *v1.CreateExperienceEntryRequest) (*experience.Experience, error) {
	pbObject := pb.GetExperience()
	if pbObject == nil {
		return nil, errors.New("pb experience is nil")
	}
	exp := &experience.Experience{
		CompanyName:     pbObject.CompanyName,
		Position:        pbObject.Position,
		PeriodStart:     pbObject.PeriodStart,
		PeriodEnd:       *pbObject.PeriodEnd,
		RoleDescription: pbObject.RoleDescription,
	}
	return exp, nil
}

func ExperiencesToPb(expList []experience.Experience) ([]*v1.Experience, error) {
	pbExpList := make([]*v1.Experience, len(expList))
	for i, e := range expList {
		pbExp, err := ExperienceToPb(&e)
		if err != nil {
			return nil, err
		}
		pbExpList[i] = pbExp
	}
	return pbExpList, nil
}

func ExperienceToPb(exp *experience.Experience) (*v1.Experience, error) {
	pbExp := &v1.Experience{
		CompanyName:     exp.CompanyName,
		Position:        exp.Position,
		PeriodStart:     exp.PeriodStart,
		PeriodEnd:       &exp.PeriodEnd,
		RoleDescription: exp.RoleDescription,
	}
	return pbExp, nil
}
