package usecase

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"taskManager/user/internal/model"
	"time"
)

type Client struct {
	Repo     ClientRepo
	producer ClientEventStorage
}

func NewClient(repo ClientRepo, producer ClientEventStorage) *Client {
	return &Client{
		Repo:     repo,
		producer: producer,
	}
}

func (c *Client) Create(ctx context.Context, user model.User) (model.User, error) {
	user.Password, _ = c.hashNewPassword(user.Password)

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	newUser, err := c.Repo.Create(ctx, &user)
	if err != nil {
		return model.User{}, err
	}

	if err = c.producer.Push(ctx, newUser); err != nil {
		return model.User{}, err
	}

	return newUser, nil
}

func (c *Client) Login(ctx context.Context, user model.User) (int64, error) {
	id, err := c.Repo.Login(ctx, &user)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (c *Client) Get(ctx context.Context, id int64) (model.User, error) {
	user, err := c.Repo.GetByID(ctx, id)
	if err != nil {
		return model.User{}, err
	}

	return *user, nil
}

func (c *Client) Delete(ctx context.Context, id int64) error {
	err := c.Repo.DeleteByID(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) hashNewPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (c *Client) checkCurrentPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
