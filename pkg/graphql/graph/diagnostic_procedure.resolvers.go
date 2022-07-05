package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"
	"time"

	graph_models "github.com/tensoremr/server/pkg/graphql/graph/model"
	"github.com/tensoremr/server/pkg/middleware"
	"github.com/tensoremr/server/pkg/models"
	"github.com/tensoremr/server/pkg/repository"
	deepCopy "github.com/ulule/deepcopier"
)

func (r *mutationResolver) OrderDiagnosticProcedure(ctx context.Context, input graph_models.OrderDiagnosticProcedureInput) (*models.DiagnosticProcedureOrder, error) {
	// Get current user
	gc, err := middleware.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}

	email := gc.GetString("email")
	if len(email) == 0 {
		return nil, errors.New("Cannot find user")
	}

	var userRepository repository.UserRepository
	var user models.User

	if err := userRepository.GetByEmail(&user, email); err != nil {
		return nil, err
	}

	var repository repository.DiagnosticProcedureOrderRepository
	// Save diagnostic procedure
	var diagnosticProcedureOrder models.DiagnosticProcedureOrder
	if err := repository.Save(&diagnosticProcedureOrder, input.DiagnosticProcedureTypeID, input.PatientChartID, input.PatientID, input.BillingID, user, input.OrderNote, input.ReceptionNote); err != nil {
		return nil, err
	}

	return &diagnosticProcedureOrder, nil
}

func (r *mutationResolver) OrderAndConfirmDiagnosticProcedure(ctx context.Context, input graph_models.OrderAndConfirmDiagnosticProcedureInput) (*models.DiagnosticProcedureOrder, error) {
	// Get current user
	gc, err := middleware.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}

	email := gc.GetString("email")
	if len(email) == 0 {
		return nil, errors.New("Cannot find user")
	}

	var userRepository repository.UserRepository
	var user models.User

	if err := userRepository.GetByEmail(&user, email); err != nil {
		return nil, err
	}

	var appointmentRepository repository.AppointmentRepository
	var appointment models.Appointment
	if err := appointmentRepository.Get(&appointment, input.AppointmentID); err != nil {
		return nil, err
	}

	var patientChartRepository repository.PatientChartRepository
	var patientChart models.PatientChart
	if err := patientChartRepository.GetByAppointmentID(&patientChart, appointment.ID); err != nil {
		return nil, err
	}

	var repository repository.DiagnosticProcedureOrderRepository
	var diagnosticProcedureOrder models.DiagnosticProcedureOrder
	if err := repository.Save(&diagnosticProcedureOrder, input.DiagnosticProcedureTypeID, patientChart.ID, appointment.PatientID, input.BillingID, user, input.OrderNote, ""); err != nil {
		return nil, err
	}

	if err := repository.Confirm(&diagnosticProcedureOrder, diagnosticProcedureOrder.ID, input.InvoiceNo); err != nil {
		return nil, err
	}

	return &diagnosticProcedureOrder, nil
}

func (r *mutationResolver) ConfirmDiagnosticProcedureOrder(ctx context.Context, id int, invoiceNo string) (*models.DiagnosticProcedureOrder, error) {
	var entity models.DiagnosticProcedureOrder
	var repository repository.DiagnosticProcedureOrderRepository

	if err := repository.Confirm(&entity, id, invoiceNo); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdateDiagnosticProcedureOrder(ctx context.Context, input graph_models.DiagnosticProcedureOrderUpdateInput) (*models.DiagnosticProcedureOrder, error) {
	var entity models.DiagnosticProcedureOrder
	deepCopy.Copy(&input).To(&entity)

	var repository repository.DiagnosticProcedureOrderRepository
	if input.Status != nil {
		entity.Status = models.DiagnosticProcedureOrderStatus(*input.Status)
	}

	if err := repository.Update(&entity); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) SaveDiagnosticProcedure(ctx context.Context, input graph_models.DiagnosticProcedureInput) (*models.DiagnosticProcedure, error) {
	var entity models.DiagnosticProcedure
	deepCopy.Copy(&input).To(&entity)

	var repository repository.DiagnosticProcedureRepository
	if err := repository.Save(&entity); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdateDiagnosticProcedure(ctx context.Context, input graph_models.DiagnosticProcedureUpdateInput) (*models.DiagnosticProcedure, error) {
	var entity models.DiagnosticProcedure
	deepCopy.Copy(&input).To(&entity)

	// Images ...
	for _, fileUpload := range input.Images {
		fileName, hashedFileName, hash, ext := HashFileName(fileUpload.Name)

		err := WriteFile(fileUpload.File.File, hashedFileName+"."+ext)
		if err != nil {
			return nil, err
		}

		entity.Images = append(entity.Images, models.File{
			ContentType: fileUpload.File.ContentType,
			Size:        fileUpload.File.Size,
			FileName:    fileName,
			Extension:   ext,
			Hash:        hash,
		})
	}

	// Documents
	for _, fileUpload := range input.Documents {
		fileName, hashedFileName, hash, ext := HashFileName(fileUpload.Name)
		err := WriteFile(fileUpload.File.File, hashedFileName+"."+ext)
		if err != nil {
			return nil, err
		}

		entity.Documents = append(entity.Documents, models.File{
			ContentType: fileUpload.File.ContentType,
			Size:        fileUpload.File.Size,
			FileName:    fileName,
			Extension:   ext,
			Hash:        hash,
		})
	}

	if input.Status != nil {
		entity.Status = models.DiagnosticProcedureStatus(*input.Status)
	}

	var repository repository.DiagnosticProcedureRepository
	if err := repository.Update(&entity); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) DeleteDiagnosticProcedure(ctx context.Context, id int) (bool, error) {
	var repository repository.DiagnosticProcedureRepository

	if err := repository.Delete(id); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) SaveDiagnosticProcedureType(ctx context.Context, input graph_models.DiagnosticProcedureTypeInput) (*models.DiagnosticProcedureType, error) {
	var entity models.DiagnosticProcedureType
	deepCopy.Copy(&input).To(&entity)

	var billingRepository repository.BillingRepository
	billings, err := billingRepository.GetByIds(input.BillingIds)
	if err != nil {
		return nil, err
	}

	entity.Billings = billings

	var diagnosticProcedureTypeRepository repository.DiagnosticProcedureTypeRepository
	if err := diagnosticProcedureTypeRepository.Save(&entity); err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *mutationResolver) UpdateDiagnosticProcedureType(ctx context.Context, input graph_models.DiagnosticProcedureTypeUpdateInput) (*models.DiagnosticProcedureType, error) {
	var entity models.DiagnosticProcedureType
	deepCopy.Copy(&input).To(&entity)

	var billingRepository repository.BillingRepository
	billings, err := billingRepository.GetByIds(input.BillingIds)
	if err != nil {
		return nil, err
	}

	entity.Billings = billings

	var diagnosticProcedureTypeRepository repository.DiagnosticProcedureTypeRepository

	if err := diagnosticProcedureTypeRepository.Update(&entity); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) DeleteDiagnosticProcedureType(ctx context.Context, id int) (bool, error) {
	var repository repository.DiagnosticProcedureTypeRepository
	if err := repository.Delete(id); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) DeleteDiagnosticImage(ctx context.Context, input graph_models.DiagnosticProcedureDeleteFileInput) (bool, error) {
	var repository repository.DiagnosticProcedureRepository
	if err := repository.DeleteFile("Images", input.DiagnosticProcedureID, input.FileID); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) DeleteDiagnosticRightEyeImage(ctx context.Context, input graph_models.DiagnosticProcedureDeleteFileInput) (bool, error) {
	var repository repository.DiagnosticProcedureRepository

	if err := repository.DeleteFile("RightEyeImages", input.DiagnosticProcedureID, input.FileID); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) DeleteDiagnosticLeftEyeImage(ctx context.Context, input graph_models.DiagnosticProcedureDeleteFileInput) (bool, error) {
	var repository repository.DiagnosticProcedureRepository

	if err := repository.DeleteFile("LeftEyeImages", input.DiagnosticProcedureID, input.FileID); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) DeleteDiagnosticRightEyeSketch(ctx context.Context, input graph_models.DiagnosticProcedureDeleteFileInput) (bool, error) {
	var repository repository.DiagnosticProcedureRepository

	if err := repository.DeleteFile("RightEyeSketches", input.DiagnosticProcedureID, input.FileID); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) DeleteDiagnosticLeftEyeSketch(ctx context.Context, input graph_models.DiagnosticProcedureDeleteFileInput) (bool, error) {
	var repository repository.DiagnosticProcedureRepository

	if err := repository.DeleteFile("LeftEyeSketches", input.DiagnosticProcedureID, input.FileID); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) DeleteDiagnosticDocument(ctx context.Context, input graph_models.DiagnosticProcedureDeleteFileInput) (bool, error) {
	var repository repository.DiagnosticProcedureRepository

	if err := repository.DeleteFile("Documents", input.DiagnosticProcedureID, input.FileID); err != nil {
		return false, err
	}

	return true, nil
}

func (r *queryResolver) DiagnosticProcedure(ctx context.Context, filter graph_models.DiagnosticProcedureFilter) (*models.DiagnosticProcedure, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) DiagnosticProcedures(ctx context.Context, page models.PaginationInput, filter *graph_models.DiagnosticProcedureFilter) (*graph_models.DiagnosticProcedureConnection, error) {
	var f models.DiagnosticProcedure
	if filter != nil {
		deepCopy.Copy(filter).To(&f)
	}

	var repository repository.DiagnosticProcedureRepository
	procedures, count, err := repository.GetAll(page, &f)

	if err != nil {
		return nil, err
	}

	edges := make([]*graph_models.DiagnosticProcedureEdge, len(procedures))

	for i, entity := range procedures {
		e := entity

		edges[i] = &graph_models.DiagnosticProcedureEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(procedures, count, page)
	return &graph_models.DiagnosticProcedureConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}

func (r *queryResolver) DiagnosticProcedureOrder(ctx context.Context, patientChartID int) (*models.DiagnosticProcedureOrder, error) {
	var entity models.DiagnosticProcedureOrder
	var repository repository.DiagnosticProcedureOrderRepository
	if err := repository.GetByPatientChartID(&entity, patientChartID); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *queryResolver) SearchDiagnosticProcedureOrders(ctx context.Context, page models.PaginationInput, filter *graph_models.DiagnosticProcedureOrderFilter, date *time.Time, searchTerm *string) (*graph_models.DiagnosticProcedureOrderConnection, error) {
	var f models.DiagnosticProcedureOrder
	if filter != nil {
		deepCopy.Copy(filter).To(&f)
	}

	if filter.Status != nil {
		f.Status = models.DiagnosticProcedureOrderStatus(*filter.Status)
	}

	var repository repository.DiagnosticProcedureOrderRepository
	result, count, err := repository.Search(page, &f, date, searchTerm, false)

	if err != nil {
		return nil, err
	}

	edges := make([]*graph_models.DiagnosticProcedureOrderEdge, len(result))

	for i, entity := range result {
		e := entity

		edges[i] = &graph_models.DiagnosticProcedureOrderEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(result, count, page)
	return &graph_models.DiagnosticProcedureOrderConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}

func (r *queryResolver) DiagnosticProcedureTypes(ctx context.Context, page models.PaginationInput, searchTerm *string) (*graph_models.DiagnosticProcedureTypeConnection, error) {
	var repository repository.DiagnosticProcedureTypeRepository
	result, count, err := repository.GetAll(page, searchTerm)
	if err != nil {
		return nil, err
	}

	edges := make([]*graph_models.DiagnosticProcedureTypeEdge, len(result))

	for i, entity := range result {
		e := entity

		edges[i] = &graph_models.DiagnosticProcedureTypeEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(result, count, page)
	return &graph_models.DiagnosticProcedureTypeConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}

func (r *queryResolver) Refraction(ctx context.Context, patientChartID int) (*models.DiagnosticProcedure, error) {
	var repository repository.DiagnosticProcedureRepository
	var entity models.DiagnosticProcedure

	if err := repository.GetRefraction(&entity, patientChartID); err != nil {
		return nil, nil
	}

	return &entity, nil
}
