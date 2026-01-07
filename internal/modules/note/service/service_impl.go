package service

import (
	"context"
	"errors"
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/infrastructure/logger"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/note/domain"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/note/dto"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/note/repository"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/notification"
	notifService "github.com/M1ralai/go-modular-monolith-template/internal/modules/notification/service"
)

type noteService struct {
	repo        repository.NoteRepository
	logger      *logger.ZapLogger
	broadcaster *notifService.Broadcaster
}

func NewNoteService(repo repository.NoteRepository, logger *logger.ZapLogger, broadcaster *notifService.Broadcaster) NoteService {
	return &noteService{
		repo:        repo,
		logger:      logger,
		broadcaster: broadcaster,
	}
}

func (s *noteService) Create(ctx context.Context, req *dto.CreateNoteRequest, userID int) (*dto.NoteResponse, error) {
	s.logger.Info("Creating note", map[string]interface{}{
		"user_id": userID,
		"title":   req.Title,
		"action":  "CREATE_NOTE",
	})

	now := time.Now()
	note := &domain.Note{
		UserID:      userID,
		CourseID:    req.CourseID,
		ComponentID: req.ComponentID,
		LifeAreaID:  req.LifeAreaID,
		Title:       req.Title,
		Content:     req.Content,
		IsFavorite:  false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	created, err := s.repo.Create(ctx, note)
	if err != nil {
		s.logger.Error("Failed to create note", err, map[string]interface{}{
			"user_id": userID,
			"title":   req.Title,
			"action":  "CREATE_NOTE_FAILED",
		})
		return nil, err
	}

	s.logger.Info("Note created successfully", map[string]interface{}{
		"user_id": userID,
		"note_id": created.ID,
		"action":  "CREATE_NOTE_SUCCESS",
	})

	response := dto.ToNoteResponse(created)
	if s.broadcaster != nil {
		s.broadcaster.Publish(userID, notification.EventNoteCreated, map[string]interface{}{
			"note_id": created.ID,
			"note":    response,
		})
		s.logger.Info("WebSocket event published", map[string]interface{}{
			"event_type": notification.EventNoteCreated,
			"user_id":    userID,
			"entity_id":  created.ID,
			"action":     "WS_EVENT_PUBLISHED",
		})
	}

	return response, nil
}

func (s *noteService) GetByID(ctx context.Context, id, userID int) (*dto.NoteResponse, error) {
	note, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if note == nil {
		return nil, errors.New("note not found")
	}
	if note.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	response := dto.ToNoteResponse(note)

	backlinks, _ := s.repo.GetBacklinks(ctx, id)
	response.BacklinkCount = len(backlinks)

	backlinkInfos := make([]*dto.NoteLinkInfo, 0, len(backlinks))
	for _, link := range backlinks {
		sourceNote, _ := s.repo.GetByID(ctx, link.SourceNoteID)
		if sourceNote != nil {
			backlinkInfos = append(backlinkInfos, &dto.NoteLinkInfo{
				ID:       link.ID,
				NoteID:   sourceNote.ID,
				Title:    sourceNote.Title,
				LinkText: link.LinkText,
			})
		}
	}
	response.Backlinks = backlinkInfos

	outgoingLinks, _ := s.repo.GetOutgoingLinks(ctx, id)
	outgoingInfos := make([]*dto.NoteLinkInfo, 0, len(outgoingLinks))
	for _, link := range outgoingLinks {
		targetNote, _ := s.repo.GetByID(ctx, link.TargetNoteID)
		if targetNote != nil {
			outgoingInfos = append(outgoingInfos, &dto.NoteLinkInfo{
				ID:       link.ID,
				NoteID:   targetNote.ID,
				Title:    targetNote.Title,
				LinkText: link.LinkText,
			})
		}
	}
	response.OutgoingLinks = outgoingInfos

	return response, nil
}

func (s *noteService) GetAll(ctx context.Context, userID int) ([]*dto.NoteResponse, error) {
	notes, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return dto.ToNoteResponseList(notes), nil
}

func (s *noteService) GetByCourse(ctx context.Context, courseID, userID int) ([]*dto.NoteResponse, error) {
	notes, err := s.repo.GetByCourseID(ctx, courseID)
	if err != nil {
		return nil, err
	}

	filtered := make([]*dto.NoteResponse, 0)
	for _, note := range notes {
		if note.UserID == userID {
			filtered = append(filtered, dto.ToNoteResponse(note))
		}
	}

	return filtered, nil
}

func (s *noteService) GetByLifeArea(ctx context.Context, lifeAreaID, userID int) ([]*dto.NoteResponse, error) {
	notes, err := s.repo.GetByLifeAreaID(ctx, lifeAreaID)
	if err != nil {
		return nil, err
	}

	filtered := make([]*dto.NoteResponse, 0)
	for _, note := range notes {
		if note.UserID == userID {
			filtered = append(filtered, dto.ToNoteResponse(note))
		}
	}

	return filtered, nil
}

func (s *noteService) GetFavorites(ctx context.Context, userID int) ([]*dto.NoteResponse, error) {
	notes, err := s.repo.GetFavorites(ctx, userID)
	if err != nil {
		return nil, err
	}

	return dto.ToNoteResponseList(notes), nil
}

func (s *noteService) Search(ctx context.Context, userID int, query string) ([]*dto.NoteResponse, error) {
	notes, err := s.repo.Search(ctx, userID, query)
	if err != nil {
		return nil, err
	}

	return dto.ToNoteResponseList(notes), nil
}

func (s *noteService) Update(ctx context.Context, id int, req *dto.UpdateNoteRequest, userID int) (*dto.NoteResponse, error) {
	s.logger.Info("Updating note", map[string]interface{}{
		"user_id": userID,
		"note_id": id,
		"action":  "UPDATE_NOTE",
	})

	note, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if note == nil {
		return nil, errors.New("note not found")
	}
	if note.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	if req.CourseID != nil {
		note.CourseID = req.CourseID
	}
	if req.ComponentID != nil {
		note.ComponentID = req.ComponentID
	}
	if req.LifeAreaID != nil {
		note.LifeAreaID = req.LifeAreaID
	}
	if req.Title != nil {
		note.Title = *req.Title
	}
	if req.Content != nil {
		note.Content = *req.Content
	}
	if req.IsFavorite != nil {
		note.IsFavorite = *req.IsFavorite
	}

	note.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, note); err != nil {
		s.logger.Error("Failed to update note", err, map[string]interface{}{
			"user_id": userID,
			"note_id": id,
			"action":  "UPDATE_NOTE_FAILED",
		})
		return nil, err
	}

	s.logger.Info("Note updated successfully", map[string]interface{}{
		"user_id": userID,
		"note_id": id,
		"action":  "UPDATE_NOTE_SUCCESS",
	})

	response := dto.ToNoteResponse(note)
	if s.broadcaster != nil {
		s.broadcaster.Publish(userID, notification.EventNoteUpdated, map[string]interface{}{
			"note_id": id,
			"note":    response,
		})
		s.logger.Info("WebSocket event published", map[string]interface{}{
			"event_type": notification.EventNoteUpdated,
			"user_id":    userID,
			"entity_id":  id,
			"action":     "WS_EVENT_PUBLISHED",
		})
	}

	return response, nil
}

func (s *noteService) Delete(ctx context.Context, id, userID int) error {
	s.logger.Info("Deleting note", map[string]interface{}{
		"user_id": userID,
		"note_id": id,
		"action":  "DELETE_NOTE",
	})

	note, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if note == nil {
		return errors.New("note not found")
	}
	if note.UserID != userID {
		return errors.New("unauthorized")
	}

	if err := s.repo.DeleteLinksByNote(ctx, id); err != nil {
		s.logger.Error("Failed to delete note links", err, map[string]interface{}{
			"note_id": id,
		})
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete note", err, map[string]interface{}{
			"user_id": userID,
			"note_id": id,
			"action":  "DELETE_NOTE_FAILED",
		})
		return err
	}

	s.logger.Info("Note deleted successfully", map[string]interface{}{
		"user_id": userID,
		"note_id": id,
		"action":  "DELETE_NOTE_SUCCESS",
	})

	if s.broadcaster != nil {
		s.broadcaster.Publish(userID, notification.EventNoteDeleted, map[string]interface{}{
			"note_id": id,
			"title":   note.Title,
		})
		s.logger.Info("WebSocket event published", map[string]interface{}{
			"event_type": notification.EventNoteDeleted,
			"user_id":    userID,
			"entity_id":  id,
			"action":     "WS_EVENT_PUBLISHED",
		})
	}

	return nil
}

func (s *noteService) CreateLink(ctx context.Context, sourceNoteID int, req *dto.CreateNoteLinkRequest, userID int) (*dto.NoteLinkResponse, error) {
	s.logger.Info("Creating note link", map[string]interface{}{
		"user_id":        userID,
		"source_note_id": sourceNoteID,
		"target_note_id": req.TargetNoteID,
		"action":         "CREATE_NOTE_LINK",
	})

	sourceNote, err := s.repo.GetByID(ctx, sourceNoteID)
	if err != nil || sourceNote == nil {
		return nil, errors.New("source note not found")
	}
	if sourceNote.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	targetNote, err := s.repo.GetByID(ctx, req.TargetNoteID)
	if err != nil || targetNote == nil {
		return nil, errors.New("target note not found")
	}
	if targetNote.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	link := &domain.NoteLink{
		SourceNoteID: sourceNoteID,
		TargetNoteID: req.TargetNoteID,
		LinkText:     req.LinkText,
		CreatedAt:    time.Now(),
	}

	created, err := s.repo.CreateLink(ctx, link)
	if err != nil {
		s.logger.Error("Failed to create note link", err, map[string]interface{}{
			"source_note_id": sourceNoteID,
			"target_note_id": req.TargetNoteID,
			"action":         "CREATE_NOTE_LINK_FAILED",
		})
		return nil, err
	}

	s.logger.Info("Note link created", map[string]interface{}{
		"link_id":        created.ID,
		"source_note_id": sourceNoteID,
		"target_note_id": req.TargetNoteID,
		"action":         "CREATE_NOTE_LINK_SUCCESS",
	})

	return dto.ToNoteLinkResponse(created), nil
}

func (s *noteService) GetBacklinks(ctx context.Context, noteID, userID int) ([]*dto.NoteLinkInfo, error) {
	note, err := s.repo.GetByID(ctx, noteID)
	if err != nil {
		return nil, err
	}
	if note == nil {
		return nil, errors.New("note not found")
	}
	if note.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	backlinks, err := s.repo.GetBacklinks(ctx, noteID)
	if err != nil {
		return nil, err
	}

	result := make([]*dto.NoteLinkInfo, 0, len(backlinks))
	for _, link := range backlinks {
		sourceNote, _ := s.repo.GetByID(ctx, link.SourceNoteID)
		if sourceNote != nil {
			result = append(result, &dto.NoteLinkInfo{
				ID:       link.ID,
				NoteID:   sourceNote.ID,
				Title:    sourceNote.Title,
				LinkText: link.LinkText,
			})
		}
	}

	return result, nil
}

func (s *noteService) DeleteLink(ctx context.Context, linkID, userID int) error {
	s.logger.Info("Deleting note link", map[string]interface{}{
		"user_id": userID,
		"link_id": linkID,
		"action":  "DELETE_NOTE_LINK",
	})

	if err := s.repo.DeleteLink(ctx, linkID); err != nil {
		s.logger.Error("Failed to delete note link", err, map[string]interface{}{
			"link_id": linkID,
			"action":  "DELETE_NOTE_LINK_FAILED",
		})
		return err
	}

	s.logger.Info("Note link deleted", map[string]interface{}{
		"link_id": linkID,
		"action":  "DELETE_NOTE_LINK_SUCCESS",
	})

	return nil
}
