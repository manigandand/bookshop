package user

import (
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrRepoUserNotFound        = errors.New("user: not found")
	ErrRepoUserInvalidPassword = errors.New("user: invalid password")
)

type Repo interface {
	// Get returns a single user matching scopes.
	// Returns ErrRepoUserNotFound if scopes doesn't match any.
	Get(db *gorm.DB, ID int) (*User, error)

	// GetByEmail retuns a single user matching the email.
	// Return ErrRepoUserNotFound if email doesn't match any.
	GetByEmail(db *gorm.DB, email string) (*User, error)

	// Authenticate validates the email and password and
	// returns User instance if validation is success
	Authenticate(db *gorm.DB, email, password string) (*User, error)

	// Find returns slice of users matching the scopes.
	// Find(db *gorm.DB, scopes ...Scope) ([]User, error)

	// // Save either creates/update the user matching the scope.
	// Save(db *gorm.DB, u *User, scopes ...Scope) error

	// // Delete remove the users matching the scopes.
	// Delete(db *gorm.DB, scopes ...Scope) error
}

// repo implements simple Repo.
type repo struct {
	logger log.Logger
}

func NewRepo(logger log.Logger) Repo {
	return &repo{logger: logger}
}

func (r *repo) Get(db *gorm.DB, ID int) (*User, error) {
	var u User
	if err := db.First(&u, "id=?", ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrRepoUserNotFound
		}
		return nil, errors.Wrapf(err, "user.repo.Get")
	}

	return &u, nil
}

func (r *repo) GetByEmail(db *gorm.DB, email string) (*User, error) {
	var u User
	if err := db.First(&u, "email=?", email).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrRepoUserNotFound
		}
		return nil, errors.Wrapf(err, "user with email=%s not found", email)
	}
	return &u, nil
}

// Authenticate validates email and password.
// Returns a valid user if validation is success.
func (r *repo) Authenticate(db *gorm.DB, email, password string) (*User, error) {
	user, err := r.GetByEmail(db, email)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%+v\n", user)
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrRepoUserInvalidPassword
	}
	return user, nil
}