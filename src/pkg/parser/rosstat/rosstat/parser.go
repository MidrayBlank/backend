package rosstat

import (
	"context"
	"fmt"

	"backend/src/pkg/parser/rosstat/domain"
	"backend/src/pkg/parser/rosstat/rosstat/storage"
	"backend/src/pkg/parser/rosstat/rosstat/subparser"

	"golang.org/x/sync/errgroup"
)

type RosstatParser struct {
	population subparser.PopulationSubparser
	// birth               subparser.BirthSubparser
	// death               subparser.DeathSubparser
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
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return p.population.Parse(ctx, p.storage)
	})

	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("RosstatParser: async parsing failed: %w", err)
	}

	return p.storage.Result(), nil
}
