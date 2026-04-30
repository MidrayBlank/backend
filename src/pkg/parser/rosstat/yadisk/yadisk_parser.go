package yadisk

import (
	"context"
	"fmt"

	"backend/src/pkg/parser/rosstat/domain"
	"backend/src/pkg/parser/rosstat/yadisk/config"
	"backend/src/pkg/parser/rosstat/yadisk/storage"
	"backend/src/pkg/parser/rosstat/yadisk/subparser"

	"golang.org/x/sync/errgroup"
)

type YadiskRosstatParser struct {
	population          subparser.PopulationSubparser
	birth               subparser.BirthSubparser
	death               subparser.DeathSubparser
	arrival             subparser.ArrivalSubparser
	departure           subparser.DepartureSubparser
	maleFemaleAge       subparser.MaleFemaleAgeSubparser
	landArea            subparser.LandAreaSubparser
	averageSalary       subparser.AverageSalarySubparser
	medicalFacilities   subparser.MedicialFacilitiesSubparser
	schools             subparser.SchoolsSubparser
	housingCommissioned subparser.HousingCommissionedSubparser

	storage *storage.Storage
	config  *config.Config
}

func NewYadiskRosstatParser() *YadiskRosstatParser {
	return &YadiskRosstatParser{
		storage: storage.NewStorage(),
		config:  config.NewConfig(),
	}
}

func (p *YadiskRosstatParser) Parse(ctx context.Context) (domain.RosstatParsedSlice, error) {
	if err := p.population.Parse(ctx, p.storage, p.config.PopulationURL); err != nil {
		return nil, fmt.Errorf("population load failed: %w", err)
	}

	if err := p.birth.Parse(ctx, p.storage, p.config.BirthURL); err != nil {
		return nil, fmt.Errorf("births load failed: %w", err)
	}

	if err := p.death.Parse(ctx, p.storage, p.config.DeathURL); err != nil {
		return nil, fmt.Errorf("deaths load failed: %w", err)
	}

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return p.landArea.Parse(ctx, p.storage, p.config.LandAreaURL)
	})

	g.Go(func() error {
		return p.medicalFacilities.Parse(ctx, p.storage, p.config.MedicialFacilitiesURL)
	})

	g.Go(func() error {
		return p.schools.Parse(ctx, p.storage, p.config.SchoolsURL)
	})

	g.Go(func() error {
		return p.housingCommissioned.Parse(ctx, p.storage, p.config.HousingCommissionedURL)
	})

	g.Go(func() error {
		return p.averageSalary.Parse(ctx, p.storage, p.config.AverageSalaryURL)
	})

	g.Go(func() error {
		return p.maleFemaleAge.Parse(ctx, p.storage, p.config.MaleFemaleAgeURL)
	})

	g.Go(func() error {
		return p.arrival.Parse(ctx, p.storage, p.config.ArrivalURL)
	})

	g.Go(func() error {
		return p.departure.Parse(ctx, p.storage, p.config.DepartureURL)
	})

	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("parallel loading failed: %w", err)
	}

	return p.storage.Result(), nil
}
