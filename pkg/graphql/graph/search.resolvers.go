package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	graph_models "github.com/tensoremr/server/pkg/graphql/graph/model"
	"github.com/tensoremr/server/pkg/repository"
)

func (r *queryResolver) Search(ctx context.Context, searchTerm string) (*graph_models.SearchResult, error) {
	var patient repository.Patient

	patients, err := patient.Search(searchTerm)
	if err != nil {
		return nil, err
	}

	var user repository.User
	providers, err := user.SearchPhysicians(searchTerm)
	if err != nil {
		return nil, err
	}

	searchResult := &model.SearchResult{
		Patients:  patients,
		Providers: providers,
	}

	return searchResult, nil
}
