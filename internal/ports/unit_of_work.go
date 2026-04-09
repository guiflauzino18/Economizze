package ports

import "context"

// TxRepositories agrupa todos os repositório dentro de uma transação
type TxRepositories struct {
	Accounts     AccountRepository
	Transactions TransactionRepository
	Budgets      BudgetRepository
	Categories   CategoryRepository
}

type UnitOfWork interface {
	Execute(ctx context.Context, fn func(repos TxRepositories) error) error
}
