package usecase

import (
	"context"

	"github.com/dinorain/kalobranded/config"
	"github.com/dinorain/kalobranded/internal/models"
	"github.com/dinorain/kalobranded/internal/session"
)

// Session use case
type sessionUC struct {
	sessionRepo session.SessRepository
	cfg         *config.Config
}

var _ session.SessUseCase = (*sessionUC)(nil)

// New session use case constructor
func NewSessionUseCase(sessionRepo session.SessRepository, cfg *config.Config) session.SessUseCase {
	return &sessionUC{sessionRepo: sessionRepo, cfg: cfg}
}

// Create new session
func (u *sessionUC) CreateSession(ctx context.Context, session *models.Session, expire int) (string, error) {
	return u.sessionRepo.CreateSession(ctx, session, expire)
}

// Delete session by id
func (u *sessionUC) DeleteById(ctx context.Context, sessionID string) error {
	return u.sessionRepo.DeleteById(ctx, sessionID)
}

// get session by id
func (u *sessionUC) GetSessionById(ctx context.Context, sessionID string) (*models.Session, error) {
	return u.sessionRepo.GetSessionById(ctx, sessionID)
}
