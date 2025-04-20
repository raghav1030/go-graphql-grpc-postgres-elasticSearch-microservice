package account

import (
	"context"

	"github.com/segmentio/ksuid"
)

type Service interface {
	PostAccount(ctx context.Context, name string) (*Account, error)
	GetAccount(ctx context.Context, id string) (*Account, error)
	GetAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error)
}

type Account struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type accountService struct {
	respository Repository
}

func NewService(r Repository) Service {
	return &accountService{r}
}

func (s *accountService) PostAccount(ctx context.Context, name string) (*Account, error) {
	a := &Account{
		Id:   ksuid.New().String(),
		Name: name,
	}
	err := s.respository.PutAccount(ctx, *a)
	if err != nil {
		return nil, err
	}
	return a, nil
}
func (s *accountService) GetAccount(ctx context.Context, id string) (*Account, error) {
	return s.respository.GetAccountById(ctx, id)

}
func (s *accountService) GetAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error) {
	if take > 100 || (skip == 0 && take == 0){
		take = 100
	}
	return s.respository.ListAccounts(ctx, skip, take)
}
