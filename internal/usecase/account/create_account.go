// Responsabilidade: criar uma nova conta financeira para o usuário.
// O use case valida as regras de negócio que envolvem múltiplos aggregates (ex: limite de contas por usuário) e delega a criação da entidade ao domínio.

package account

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/guiflauzino18/economizze/internal/domain/aggregates"
	"github.com/guiflauzino18/economizze/internal/domain/vos"
	"github.com/guiflauzino18/economizze/internal/ports"
)

type CreateAccountInput struct {
	userID         uuid.UUID // extraído do token JWT
	Name           string
	AccountType    aggregates.AccountType
	InitialBalance int64
	Currency       string
	IsDefault      bool
}

type CreateAccountOutput struct {
	Account aggregates.Account
}

type CreateAccountUseCase struct {
	accountRepo ports.AccountRepository
	publisher   ports.EventPublisher
}

func NewCreateAccountUseCase(accountRepo ports.AccountRepository, publisher ports.EventPublisher) *CreateAccountUseCase {

	return &CreateAccountUseCase{
		accountRepo: accountRepo,
		publisher:   publisher,
	}
}

func (uc *CreateAccountUseCase) Execute(ctx context.Context, in CreateAccountInput) (*CreateAccountOutput, error) {

	// Regra de negócio: limite de 10 contas por usuário
	existing, err := uc.accountRepo.FindByUserID(ctx, in.userID)
	if err != nil {
		return nil, fmt.Errorf("CreateAccount.FindByUserID: %w", err)
	}

	if len(existing) >= 10 {
		return nil, fmt.Errorf("CreateAccount: maximum 10 accounts reached: %w", err)
	}

	// Cria initialBalance via NewMoeny
	initialBalance, err := vos.NewMoney(in.InitialBalance, in.Currency)
	if err != nil {
		return nil, fmt.Errorf("CreateAccount.NewMoney: %w", err)
	}

	account, err := aggregates.NewAccount(in.userID, in.Name, in.AccountType, initialBalance)
	if err != nil {
		return nil, fmt.Errorf("CreateAccount.NewAccount: %e", err)
	}

	// Regra: se marcado como padrão desativa padrão anterior
	if in.IsDefault {
		if err := uc.ClearDefaultAccount(ctx, in.userID); err != nil {
			return nil, err
		}
	}

	// Persiste
	if err := uc.accountRepo.Save(ctx, account); err != nil {
		return nil, fmt.Errorf("CreateAccount.Save: %w", err)
	}

	// Publica domainEvents
	if len(account.Events()) > 0 {
		if err := uc.publisher.Publish(ctx, account.Events()...); err != nil {
			fmt.Printf("failed to publisher events: %s", err) // Futuramente será implementado em log
		}

		account.ClearEvents()
	}

	return &CreateAccountOutput{Account: *account}, nil

}

// ClearDefautAccount remove flag de default account de todas as contas do usuário antes de definir uma nova conta padrão
func (uc *CreateAccountUseCase) ClearDefaultAccount(ctx context.Context, userID uuid.UUID) error {

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
