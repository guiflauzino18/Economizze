package repository_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/guiflauzino18/economizze/internal/domain"
	"github.com/guiflauzino18/economizze/internal/infra/database"
	"github.com/guiflauzino18/economizze/internal/infra/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/gorm"
)

func TestAccountRepository_Integration(t *testing.T) {

	// Pula se estiver rodando com -Short (CI rápido)
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db := setupTestDB(t)
	repo := repository.NewAccountRepository(db)

	// userID fixo para os testes
	userID := uuid.New()

	t.Run("salva e recupera conta", func(t *testing.T) {

		// Cria balance
		balance, err := domain.NewMoney(150000, "BRL")
		require.NoError(t, err)

		// Cria account
		account, err := domain.NewAccount(userID, "Conta Corrente Itaú", domain.AccountTypeChecking, balance)
		require.NoError(t, err)

		// SAlva no banco
		err = repo.Save(context.Background(), account)
		require.NoError(t, err)

		// REcupera pelo ID
		found, err := repo.FindByID(context.Background(), account.ID())
		require.NoError(t, err)

		// Compara os campos que são importantes
		assert.Equal(t, account.ID(), found.ID)
		assert.Equal(t, account.Name(), found.Name())
		assert.Equal(t, account.Balance(), found.Balance())
		assert.Equal(t, account.AccountType(), found.AccountType())
		assert.True(t, found.IsActive())

	})
}

func setupTestDB(t *testing.T) *gorm.DB {

	t.Helper()
	ctx := context.Background()

	pgContainer, err := postgres.Run(
		ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("economizze-test"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
		),
	)
	require.NoError(t, err)
	t.Cleanup(func() { pgContainer.Terminate(ctx) })

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	db, err := database.NewGorm(database.Config{DSN: connStr})
	require.NoError(t, err)

	// Roda as migrations reais - teste schema e migration juntos
	wd, err := os.Getwd() // Pega caminho absoluto da pasta atual
	require.NoError(t, err)

	absPath := filepath.Join(wd, "../../../") // volta té a pasta raiz do projeto onde está migrations
	migrationsDir := os.DirFS(absPath)

	err = database.RunMigrations(db, migrationsDir)
	require.NoError(t, err)

	// Seed de categorias padrão
	sqlDB, _ := db.DB()
	err = database.SeedCategories(db)
	require.NoError(t, err)
	_ = sqlDB

	return db

}
