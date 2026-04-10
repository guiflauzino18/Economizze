package account

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/guiflauzino18/economizze/internal/domain"
	"github.com/guiflauzino18/economizze/internal/ports"
)

type ListAccountOutput struct {
	Accounts     []*domain.Account
	TotalBalance domain.Money
}

type ListAccountUseCase struct {
	accountRepo ports.AccountRepository
}

func NewListAccountUseCase(accountRepo ports.AccountRepository) *ListAccountUseCase {
	return &ListAccountUseCase{accountRepo}
}

func (uc *ListAccountUseCase) Execute(ctx context.Context, UserID uuid.UUID) (*ListAccountOutput, error) {
	accounts, err := uc.accountRepo.FindByUserID(ctx, UserID)
	if err != nil {
		return nil, fmt.Errorf("ListAccount: %w", err)
	}

	total, _ := domain.NewMoney(0, accounts[0].Currency())

	for _, a := range accounts {
		// Soma valor de mesma moeda da moedada primeira conta. Multimoeda poderá ser implementada posteriormente
		if a.Currency() == total.Currency() {
			total, _ = total.Add(a.Balance())
		}
	}

	return &ListAccountOutput{
		Accounts:     accounts,
		TotalBalance: total,
	}, nil
}
