package account

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/guiflauzino18/economizze/internal/domain"
	"github.com/guiflauzino18/economizze/internal/ports"
)

type UpdateAccountInput struct {
	AccountID uuid.UUID // para verificar ownership. Usuário só altera suas contas
	userID    uuid.UUID
	Name      *string //ponteiro = campo opcional (nil = não atualizar)
	IsDefault *bool
}

type UpdateAccountUseCase struct {
	accountRepo ports.AccountRepository
}

func NewupdateAccountUseCase(accountRepo ports.AccountRepository) *UpdateAccountUseCase {
	return &UpdateAccountUseCase{accountRepo}
}

func (uc *UpdateAccountUseCase) Execute(ctx context.Context, in UpdateAccountInput) (*domain.Account, error) {

	account, err := uc.accountRepo.FindByID(ctx, in.AccountID)
	if err != nil {
		return nil, fmt.Errorf("UpdateAccount.FindByID: %w", err)
	}

	// Verifica se a conta pertence ao usuário
	if account.UserID() != in.userID {
		return nil, fmt.Errorf("UpdateAccount: account %s: %w", in.AccountID, domain.ErrForbidden)
	}

	// Partial Update = aplica apenas os campos fornecidos

	if in.Name != nil {
		if err := account.Rename(*in.Name); err != nil {
			return nil, fmt.Errorf("UpdateAccount.Rename: %w", err)
		}
	}

	if in.IsDefault != nil && *in.IsDefault && !account.IsDefault() {
		if err := uc.clearDefaultAccount(ctx, in.userID); err != nil {
			return nil, err
		}

		account.SetDefaultAccount()
	}

	if err := uc.accountRepo.Save(ctx, account); err != nil {
		return nil, fmt.Errorf("UpdateAccount.Save: %w", err)
	}

	return account, nil
}

func (uc *UpdateAccountUseCase) clearDefaultAccount(ctx context.Context, userID uuid.UUID) error {
	accounts, err := uc.accountRepo.FindByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("clearDefaultAccount: %w", err)
	}
	for _, a := range accounts {
		if a.IsDefault() {
			a.UnsetDefault()
			if err := uc.accountRepo.Save(ctx, a); err != nil {
				return fmt.Errorf("clearDefaultAccount.Save: %w", err)
			}
		}
	}
	return nil
}
