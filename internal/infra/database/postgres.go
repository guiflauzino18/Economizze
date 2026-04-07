package database

import (
	"embed"
	"errors"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq"
	gormpg "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// Config centraliza todas as configurações do banco
type Config struct {
	DSN string // DSN completo: "postgres://user:pass@host:5432/dbname?sslmode=disable"

	// Pool de conexões — valores padrão para APIs de médio porte
	MaxOpenConns    int           // máximo de conexões abertas simultâneas
	MaxIdleConns    int           // conexões mantidas "prontas" para uso
	ConnMaxLifetime time.Duration // recria conexão após este tempo
	ConnMaxIdleTime time.Duration // fecha conexão idle após este tempo
}

// DefaultConfig retorna uma configuração segura para produção
func DefaultConfig(dsn string) Config {
	return Config{
		DSN:             dsn,
		MaxOpenConns:    25,              // suficiente para ~500 req/s
		MaxIdleConns:    10,              // evita abrir/fechar a todo momento
		ConnMaxLifetime: 5 * time.Minute, // renova conexões periodicamente
		ConnMaxIdleTime: 2 * time.Minute, // libera conexões não usadas
	}
}

// NewGorm abre a conexão GORM com o pool configurado corretamente.
// Retorna *gorm.DB que será injetado nos repositórios via DI.
func NewGorm(cfg Config) (*gorm.DB, error) {
	// Abre conexão GORM com o driver postgres
	db, err := gorm.Open(gormpg.Open(cfg.DSN), &gorm.Config{
		// PrepareStmt: reutiliza query plans para queries repetidas
		// Melhoria de ~20% em performance para queries frequentes
		PrepareStmt: true,

		// SkipDefaultTransaction: não envolve cada operação em transação automaticamente
		// Melhoria de performance para operações simples de leitura
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, fmt.Errorf("gorm.Open: %w", err)
	}

	// Obtém a *sql.DB subjacente para configurar o pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("db.DB(): %w", err)
	}

	// Configura o pool de conexões
	// IMPORTANTE: sem isso, o Go abre uma conexão nova a cada operação
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	// Verifica conectividade no startup
	// Falha rápida: melhor falhar ao iniciar do que falhar na primeira request
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("postgres ping: %w", err)
	}

	return db, nil
}

// RunMigrations executa todas as migrations pendentes.
func RunMigrations(db *gorm.DB) error {
	// Obtém *sql.DB para usar com golang-migrate
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("RunMigrations.db.DB(): %w", err)
	}

	// Cria o source de migrations a partir do filesystem embutido
	// "migrations" é o diretório dentro do embed.FS
	source, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("RunMigrations.source: %w", err)
	}

	// Cria o driver postgres para o golang-migrate
	// schema_migrations é a tabela de controle de versões
	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{
		MigrationsTable: "schema_migrations",
	})
	if err != nil {
		return fmt.Errorf("RunMigrations.driver: %w", err)
	}

	// Inicializa o migrator
	m, err := migrate.NewWithInstance("iofs", source, "postgres", driver)
	if err != nil {
		return fmt.Errorf("RunMigrations.new: %w", err)
	}

	// Verifica estado atual antes de migrar
	version, dirty, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		return fmt.Errorf("RunMigrations.version: %w", err)
	}

	// dirty=true significa que a última migration falhou no meio
	// Não é seguro continuar — exige intervenção manual
	if dirty {
		return fmt.Errorf(
			"database is in dirty state at version %d — "+
				"fix manually with: migrate force <version>", version)
	}

	// Aplica todas as migrations pendentes
	if err := m.Up(); err != nil {
		// ErrNoChange não é um erro — significa que já está atualizado
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}
		return fmt.Errorf("RunMigrations.up: %w", err)
	}

	newVersion, _, _ := m.Version()
	fmt.Printf("✓ migrations: %d → %d\n", version, newVersion)

	return nil
}
