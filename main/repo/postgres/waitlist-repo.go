package postgres

import (
	"github.com/inlovewithgo/transit-backend/main/models"
	"gorm.io/gorm"
)

type waitlistRepository struct {
	db *gorm.DB
}

func NewWaitlistRepository(db *gorm.DB) *waitlistRepository {
	return &waitlistRepository{db: db}
}

func (w *waitlistRepository) Create(waitlist *models.Waitlist) error {
	return w.db.Create(waitlist).Error
}

func (w *waitlistRepository) GetByEmail(email string) (*models.Waitlist, error) {
	var waitlist models.Waitlist
	err := w.db.Where("email = ?", email).First(&waitlist).Error
	if err != nil {
		return nil, err
	}
	return &waitlist, nil
}

func (w *waitlistRepository) Update(waitlist *models.Waitlist) error {
	return w.db.Save(waitlist).Error
}

func (w *waitlistRepository) GetAll() ([]models.Waitlist, error) {
	var waitlists []models.Waitlist
	err := w.db.Find(&waitlists).Error
	return waitlists, err
}

func (w *waitlistRepository) GetCount() (int64, error) {
	var count int64
	err := w.db.Model(&models.Waitlist{}).Count(&count).Error
	return count, err
}
