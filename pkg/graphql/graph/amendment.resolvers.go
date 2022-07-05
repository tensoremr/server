package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	graph_models "github.com/tensoremr/server/pkg/graphql/graph/model"
	"github.com/tensoremr/server/pkg/models"
	"github.com/tensoremr/server/pkg/repository"
	deepCopy "github.com/ulule/deepcopier"
)

func (r *mutationResolver) CreateAmendment(ctx context.Context, input graph_models.AmendmentInput) (*models.Amendment, error) {
	var entity models.Amendment
	deepCopy.Copy(&input).To(&entity)

	var repository repository.AmendmentRepository
	if err := repository.Create(&entity); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdateAmendment(ctx context.Context, input graph_models.AmendmentUpdateInput) (*models.Amendment, error) {
	var entity models.Amendment
	deepCopy.Copy(&input).To(&entity)

	var repository repository.AmendmentRepository
	if err := repository.Update(&entity); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) DeleteAmendment(ctx context.Context, id int) (bool, error) {
	var repository repository.AmendmentRepository
	if err := repository.Delete(id); err != nil {
		return false, err
	}

	return true, nil
}

func (r *queryResolver) Amendment(ctx context.Context, id int) (*models.Amendment, error) {
	var entity models.Amendment
	var repository repository.AmendmentRepository
	if err := repository.Get(&entity, id); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *queryResolver) Amendments(ctx context.Context, filter *graph_models.AmendmentFilter) ([]*models.Amendment, error) {
	var f models.Amendment
	if filter != nil {
		deepCopy.Copy(filter).To(&f)
	}
	var repository repository.AmendmentRepository
	result, err := repository.GetAll(&f)

	if err != nil {
		return nil, err
	}

	return result, nil
}
