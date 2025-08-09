package repo

import (
	"github.com/inlovewithgo/transit-backend/main/models"
)

type WaitlistRepository interface {
	Create(waitlist *models.Waitlist) error
	GetByEmail(email string) (*models.Waitlist, error)
	Update(waitlist *models.Waitlist) error
	GetAll() ([]models.Waitlist, error)
	GetCount() (int64, error)
}
