package model

import "github.com/jinzhu/copier"

func ToDomain[Model, D any](model *Model) (*D, error) {
	var result D
	if err := copier.Copy(&result, model); err != nil {
		return nil, err
	}
	return &result, nil
}

func ToModel[Model, Domain any](model *Model, domain *Domain) error {
	if domain == nil {
		return nil
	}
	return copier.Copy(model, domain)
}

func ToDomainSlice[Model, Domain any](models []Model) ([]Domain, error) {
	result := make([]Domain, len(models))
	for i, model := range models {
		domain, err := ToDomain[Model, Domain](&model)
		if err != nil {
			return nil, err
		}
		result[i] = *domain
	}
	return result, nil
}
