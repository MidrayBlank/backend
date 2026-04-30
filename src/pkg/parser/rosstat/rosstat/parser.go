package rosstat

import (
	"context"

	"backend/src/pkg/parser/rosstat/domain"
	"backend/src/pkg/parser/rosstat/rosstat/storage"
	"backend/src/pkg/parser/rosstat/rosstat/subparser"
)

type RosstatParser struct {
	code subparser.CodeSubparser

	population subparser.PopulationSubparser
	birth      subparser.BirthSubparser
	death      subparser.DeathSubparser
	// arrival             subparser.ArrivalSubparser
	// departure           subparser.DepartureSubparser
	// maleFemaleAge       subparser.MaleFemaleAgeSubparser
	// landArea            subparser.LandAreaSubparser
	// averageSalary       subparser.AverageSalarySubparser
	// medicalFacilities   subparser.MedicialFacilitiesSubparser
	// schools             subparser.SchoolsSubparser
	// housingCommissioned subparser.HousingCommissionedSubparser

	storage *storage.Storage
}

func NewRosstatParser() *RosstatParser {
	return &RosstatParser{
		storage: storage.NewStorage(),
	}
}

func (p *RosstatParser) Parse(ctx context.Context) (domain.RosstatParsedSlice, error) {
	err := p.code.Parse(ctx, p.storage)

	if err != nil {
		return nil, err
	}

	err = p.population.Parse(ctx, p.storage)

	if err != nil {
		return nil, err
	}

	err = p.birth.Parse(ctx, p.storage)
	if err != nil {
		return nil, err
	}

	err = p.death.Parse(ctx, p.storage)
	if err != nil {
		return nil, err
	}

	return p.storage.Result(), nil
}
