package repository

import models "github.com/ivzb/semaphore_server/user"

type UserRepository interface {
	Fetch(cursor string, num int64) ([]*models.User, error)
	GetByID(id int64) (*models.User, error)
	GetByTitle(title string) (*models.User, error)
	Update(article *models.User) (*models.User, error)
	Store(a *models.User) (int64, error)
	Delete(id int64) (bool, error)
}
