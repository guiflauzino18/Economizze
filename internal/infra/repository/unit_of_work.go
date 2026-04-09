// UnitOfWork permite que um use case execute múltiplas operações de repositório dentro de uma única transação de banco de dados.
//
// Sem UoW, cada repositório tem sua própria conexão e operações separadas não são atômicas — se o segundo Save falhar, o primeiro já foi commitado e o banco fica inconsistente.

package repository

import (
	"context"
	"fmt"

	"github.com/guiflauzino18/economizze/internal/ports"
	"gorm.io/gorm"
)

type unitOfWork struct {
	db *gorm.DB
}

func NewUnitOfWork(db *gorm.DB) ports.UnitOfWork {
	return unitOfWork{db}
}

// Execute implements [ports.UnitOfWork].
func (u unitOfWork) Execute(ctx context.Context, fn func(repos ports.TxRepositories) error) error {

	return u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		repos := ports.TxRepositories{
			Accounts:     NewAccountRepository(tx),
			Transactions: NewTransactionRepository(tx),
			Budgets:      NewBudgetRepository(tx),
			Categories:   NewCategoryRepository(tx),
		}

		if err := fn(repos); err != nil {
			// Retornar Erro faz rollback automático
			return fmt.Errorf("unitOfWork.Execute: %w", err)
		}

		// Retornar Nil faz commit automático
		return nil
	})
}
