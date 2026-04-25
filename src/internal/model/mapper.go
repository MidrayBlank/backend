package model

import "github.com/jinzhu/copier"

func ToDomain[M, D any](m *M) (*D, error) {
	var result D
	if err := copier.Copy(&result, m); err != nil {
		return nil, err
	}
	return &result, nil
}

func ToModel[M, D any](m *M, d *D) error {
	if d == nil {
		return nil
	}
	return copier.Copy(m, d)
}

func ToDomainSlice[M, D any](models []M) ([]D, error) {
	result := make([]D, len(models))
	for i, m := range models {
		d, err := ToDomain[M, D](&m)
		if err != nil {
			return nil, err
		}
		result[i] = *d
	}
	return result, nil
}
