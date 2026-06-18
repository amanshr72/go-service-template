package user

import (
	"errors"
	"go-crud2/internal/notification"
	"log"
	"strings"
)

type service struct {
	repo     Repository
	notifier notification.Client
}

func NewService(repo Repository, notifier notification.Client) Service {
	return &service{repo: repo, notifier: notifier}
}

func (s *service) Create(input CreateUserInput) (*User, error) {
	if strings.TrimSpace(input.Name) == "" {
		return nil, errors.New("name is required")
	}
	if strings.TrimSpace(input.Email) == "" {
		return nil, errors.New("email is required")
	}
	u := &User{
		Name:     strings.TrimSpace(input.Name),
		Email:    strings.TrimSpace(input.Email),
		IsActive: true, // default active on creation
	}
	if err := s.repo.Create(u); err != nil {
		return nil, err
	}
	_, err := s.notifier.SendEmail(notification.SendEmailRequest{
		To:      u.Email,
		Subject: "Welcome!",
		Body:    "Your account has been created, " + u.Name,
	})
	if err != nil {
		log.Printf("warning: welcome email failed for user %d: %v", u.ID, err)
	}
	return u, nil
}

func (s *service) GetAll() ([]User, error) {
	return s.repo.FindAll()
}

func (s *service) GetByID(id int) (*User, error) {
	return s.repo.FindByID(id)
}

func (s *service) GetByActive(active bool) ([]User, error) {
	return s.repo.FindByActive(active)
}

func (s *service) GetCount() (int, error) {
	return s.repo.Count()
}

func (s *service) Update(id int, input UpdateUserInput) (*User, error) {
	u, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if strings.TrimSpace(input.Name) != "" {
		u.Name = strings.TrimSpace(input.Name)
	}
	if strings.TrimSpace(input.Email) != "" {
		u.Email = strings.TrimSpace(input.Email)
	}
	// pointer check: nil means "not sent", false means "explicitly set to false"
	if input.IsActive != nil {
		u.IsActive = *input.IsActive
	}
	if err := s.repo.Update(u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *service) Delete(id int) error {
	if _, err := s.repo.FindByID(id); err != nil {
		return errors.New("user not found")
	}
	return s.repo.Delete(id)
}
