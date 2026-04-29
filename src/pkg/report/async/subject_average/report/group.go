package report

import (
	"backend/src/pkg/report/async/subject_average/domain"
	"context"
)

func groupBySubjectCode(ctx context.Context, rosstatSlice domain.RosstatSlice) (map[int]domain.RosstatSlice, error) {
	result := make(map[int]domain.RosstatSlice)

	for _, item := range rosstatSlice {

		if err := ctx.Err(); err != nil {
			return nil, ctx.Err()
		}

		key := item.SubjectCode
		result[key] = append(result[key], item)
	}

	return result, nil
}
