package account

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/guiflauzino18/economizze/internal/domain"
	"github.com/guiflauzino18/economizze/internal/ports"
)

type DeleteAccountInput struct {
	AccountID uuid.UUID
	UserID    uuid.UUID
}

type DeleteAccountUseCase struct {
	AccountRepo     ports.AccountRepository
	TransactionRepo ports.TransactionRepository
}

func NewDeleteAccountUseCase(AccountRepo ports.AccountRepository, TransactionRepo ports.TransactionRepository) *DeleteAccountUseCase {
	return &DeleteAccountUseCase{AccountRepo, TransactionRepo}
}

func (uc *DeleteAccountUseCase) Execute(ctx context.Context, in DeleteAccountInput) error {
	account, err := uc.AccountRepo.FindByID(ctx, in.AccountID)
	if err != nil {
		return fmt.Errorf("DeleteAccount.FindByID: %w", err)
	}

	if account.UserID() != in.UserID {
		return fmt.Errorf("DeleteAccount: %w", domain.ErrForbidden)
	}

	//Regra de negócio: Não permite delete conta com saldo
	if account.Balance().IsPositive() {
		return fmt.Errorf("DeleteAccount: account has balance, transfer funds before deleting: %w",
			domain.ErrInvalidOperation)
	}

	// Desativa via método do dominio para tratamentos
	if err := account.Deactivate(); err != nil {
		return fmt.Errorf("DeleteAccount:Deactivate: %w", err)
	}

	// Softdelete, somente alterou status para inativo (uc.Account.Delete também somente desativa a conta)
	return uc.AccountRepo.Save(ctx, account)
}
