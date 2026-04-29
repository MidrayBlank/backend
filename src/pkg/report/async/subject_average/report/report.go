package report

import (
	"backend/src/pkg/report/async/subject_average/domain"
	"context"
)

type SubjectAverageReportCompiler struct{}

func NewSubjectAverageReportCompiler() *SubjectAverageReportCompiler {
	return &SubjectAverageReportCompiler{}
}

func (r *SubjectAverageReportCompiler) Compile(ctx context.Context, data domain.RosstatSlice) (domain.RosstatSlice, error) {
	if len(data) == 0 {
		return domain.RosstatSlice{}, nil
	}

	groups, err := groupBySubjectCode(ctx, data)

	if err != nil {
		return nil, err
	}

	result := make(domain.RosstatSlice, 0, len(groups))

	for subcode, municip := range groups {
		if err := ctx.Err(); err != nil {
			return nil, ctx.Err()
		}
		year := municip[0].Year
		regionRecord := buildRegionRecord(subcode, year, municip)
		result = append(result, regionRecord)
	}

	if err := ctx.Err(); err != nil {
		return nil, ctx.Err()
	}

	return result, nil
}
