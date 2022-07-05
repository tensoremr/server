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

func (r *mutationResolver) CreatePharmacy(ctx context.Context, input graph_models.PharmacyInput) (*models.Pharmacy, error) {
	var entity models.Pharmacy
	deepCopy.Copy(&input).To(&entity)

	var repository repository.PharmacyRepository
	if err := repository.Save(&entity); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdatePharmacy(ctx context.Context, input graph_models.PharmacyUpdateInput) (*models.Pharmacy, error) {
	var entity models.Pharmacy
	deepCopy.Copy(&input).To(&entity)

	var repository repository.PharmacyRepository
	if err := repository.Update(&entity); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) DeletePharmacy(ctx context.Context, id int) (bool, error) {
	var repository repository.PharmacyRepository
	if err := repository.Delete(id); err != nil {
		return false, err
	}
	return true, nil
}

func (r *queryResolver) Pharmacy(ctx context.Context, id int) (*models.Pharmacy, error) {
	var repository repository.PharmacyRepository
	var entity models.Pharmacy

	if err := repository.Get(&entity, id); err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *queryResolver) Pharmacies(ctx context.Context, page models.PaginationInput) (*graph_models.PharmacyConnection, error) {
	var repository repository.PharmacyRepository
	result, count, err := repository.GetAll(page, nil)

	if err != nil {
		return nil, err
	}

	edges := make([]*graph_models.PharmacyEdge, len(result))

	for i, entity := range result {
		e := entity

		edges[i] = &graph_models.PharmacyEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(result, count, page)
	return &graph_models.PharmacyConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}
