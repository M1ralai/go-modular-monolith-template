package service

import (
	"context"
	"errors"
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/infrastructure/logger"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/people/domain"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/people/dto"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/people/repository"
)

type personService struct {
	repo   repository.PersonRepository
	logger *logger.ZapLogger
}

func NewPersonService(repo repository.PersonRepository, logger *logger.ZapLogger) PersonService {
	return &personService{repo: repo, logger: logger}
}

func (s *personService) Create(ctx context.Context, req *dto.CreatePersonRequest, userID int) (*dto.PersonResponse, error) {
	s.logger.Info("Creating person", map[string]interface{}{"user_id": userID, "name": req.Name, "action": "CREATE_PERSON"})
	now := time.Now()
	person := &domain.Person{UserID: userID, Name: req.Name, Email: req.Email, Phone: req.Phone, Company: req.Company, Relationship: req.Relationship, Tags: req.Tags, Notes: req.Notes, CreatedAt: now, UpdatedAt: now}
	created, err := s.repo.Create(ctx, person)
	if err != nil {
		s.logger.Error("Failed to create person", err, map[string]interface{}{"user_id": userID, "action": "CREATE_PERSON_FAILED"})
		return nil, err
	}
	s.logger.Info("Person created", map[string]interface{}{"user_id": userID, "person_id": created.ID, "action": "CREATE_PERSON_SUCCESS"})
	return dto.ToPersonResponse(created), nil
}

func (s *personService) GetByID(ctx context.Context, id, userID int) (*dto.PersonResponse, error) {
	person, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if person == nil {
		return nil, errors.New("person not found")
	}
	if person.UserID != userID {
		return nil, errors.New("unauthorized")
	}
	return dto.ToPersonResponse(person), nil
}

func (s *personService) GetAll(ctx context.Context, userID int) ([]*dto.PersonResponse, error) {
	people, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return dto.ToPersonResponseList(people), nil
}

func (s *personService) SearchByTag(ctx context.Context, userID int, tag string) ([]*dto.PersonResponse, error) {
	people, err := s.repo.SearchByTag(ctx, userID, tag)
	if err != nil {
		return nil, err
	}
	return dto.ToPersonResponseList(people), nil
}

func (s *personService) Search(ctx context.Context, userID int, query string) ([]*dto.PersonResponse, error) {
	people, err := s.repo.Search(ctx, userID, query)
	if err != nil {
		return nil, err
	}
	return dto.ToPersonResponseList(people), nil
}

func (s *personService) Update(ctx context.Context, id int, req *dto.UpdatePersonRequest, userID int) (*dto.PersonResponse, error) {
	s.logger.Info("Updating person", map[string]interface{}{"user_id": userID, "person_id": id, "action": "UPDATE_PERSON"})
	person, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if person == nil {
		return nil, errors.New("person not found")
	}
	if person.UserID != userID {
		return nil, errors.New("unauthorized")
	}
	if req.Name != nil {
		person.Name = *req.Name
	}
	if req.Email != nil {
		person.Email = *req.Email
	}
	if req.Phone != nil {
		person.Phone = *req.Phone
	}
	if req.Company != nil {
		person.Company = *req.Company
	}
	if req.Relationship != nil {
		person.Relationship = *req.Relationship
	}
	if req.Tags != nil {
		person.Tags = req.Tags
	}
	if req.Notes != nil {
		person.Notes = *req.Notes
	}
	person.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, person); err != nil {
		s.logger.Error("Failed to update person", err, map[string]interface{}{"user_id": userID, "person_id": id, "action": "UPDATE_PERSON_FAILED"})
		return nil, err
	}
	s.logger.Info("Person updated", map[string]interface{}{"user_id": userID, "person_id": id, "action": "UPDATE_PERSON_SUCCESS"})
	return dto.ToPersonResponse(person), nil
}

func (s *personService) Delete(ctx context.Context, id, userID int) error {
	s.logger.Info("Deleting person", map[string]interface{}{"user_id": userID, "person_id": id, "action": "DELETE_PERSON"})
	person, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if person == nil {
		return errors.New("person not found")
	}
	if person.UserID != userID {
		return errors.New("unauthorized")
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete person", err, map[string]interface{}{"user_id": userID, "person_id": id, "action": "DELETE_PERSON_FAILED"})
		return err
	}
	s.logger.Info("Person deleted", map[string]interface{}{"user_id": userID, "person_id": id, "action": "DELETE_PERSON_SUCCESS"})
	return nil
}
