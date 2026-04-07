package repository

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ============================================================
// UserModel — representa a linha na tabela users
//
// ============================================================
type UserModel struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey"`
	Name         string         `gorm:"not null;size:100"`
	Email        string         `gorm:"not null;uniqueIndex;type:citext"`
	PasswordHash string         `gorm:"not null;column:password_hash"`
	Role         string         `gorm:"not null;default:user;size:20"`
	Active       bool           `gorm:"not null;default:true"`
	CreatedAt    time.Time      `gorm:"not null"`
	UpdatedAt    time.Time      `gorm:"not null"`
	DeletedAt    gorm.DeletedAt `gorm:"index"` // Soft delete nativo. Quando preenchido gorm filtra automaticamente
}

// TableName sobrescreve a convenção de nome plural do GORM. Sem isso, GORM usaria "user_models" em vez de "users"
func (UserModel) TableName() string { return "users" }

// ============================================================
// AccountModel — representa a linha na tabela accounts
// ============================================================
type AccountModel struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID       uuid.UUID `gorm:"type:uuid;not null;index"` //index: cria index automaticamente nesse campo
	Name         string    `gorm:"not null;size:100"`
	AccountType  string    `gorm:"not null, column:account_type;size:20"`
	BalanceCents int64     `gorm:"not null; default:0;column:balance_cents"`
	Currency     string    `gorm:"not null;size:3;default:BRL"`
	IsDefault    bool      `gorm:"not null;default:false;column:is_default"`
	Active       bool      `gorm:"not null;default:true"`
	CreatedAt    time.Time `gorm:"not null"`
	UpdatedAt    time.Time `gorm:"not null"`
}

func (AccountModel) TableName() string { return "accounts" }

// ============================================================
// CategoryModel — representa a linha na tabela categories
// ============================================================

type CategoryModel struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey"`
	UserID      *uuid.UUID `gorm:"type:uuid;index"`
	Name        string     `gorm:"not null;size:100"`
	DefaultType *string    `gorm:"column:default_type;size:20"`
	Active      bool       `gorm:"not null;default:true"`
	CreatedAt   time.Time  `gorm:"not null"`
	UpdatedAt   time.Time  `gorm:"not null"`
}

func (CategoryModel) Tablename() string { return "categories" }

// ============================================================
// TransactionModel — representa a linha na tabela transactions
// ============================================================

type TransactionModel struct {
	ID             uuid.UUID  `gorm:"type:uuid;primaryKey"`
	AccountID      uuid.UUID  `gorm:"not null;type:uuid;index"`
	CategoryID     *uuid.UUID `gorm:"type:uuid;index"`
	TransferPeerID *uuid.UUID `gorm:"type:uuid;index;column:transfer_peer_id"`
	Cents          int64      `gorm:"not null;"`
	Currency       string     `gorm:"not null;size:3;"`
	Type           string     `gorm:"not null;size:20"`
	Description    string     `gorm:"not null;size:255"`
	Notes          *string    `gorm:"type:text"`
	OccurredOn     time.Time  `gorm:"not null;type:date;column:occurred_on"`
	RecurringID    *uuid.UUID `gorm:"type:uuid;column:recurring_id"`
	CreatedAt      time.Time  `gorm:"not null"`
	UpdatedAt      time.Time  `gorm:"not null"`

	// Preload: GORM carrega a categoria quando solicitado
	// constraint:false evita que GORM crie FK automática (já temos nas migrations)
	Category *CategoryModel `gorm:"foreignKey:CategoryID;constraint:false"`
}

func (TransactionModel) TableName() string { return "transactions" }

// ============================================================
// BudgetModel — representa a linha na tabela budgets
// ============================================================

type BudgetModel struct {
	ID                 uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID             uuid.UUID `gorm:"type:uuid;not null;index"`
	CategoryID         uuid.UUID `gorm:"type:uuid;not null"`
	PeriodStart        time.Time `gorm:"not null;type:date;column:period_start"`
	PeriodEnd          time.Time `gorm:"not null;type:date;column:period_end"`
	LimitCents         int64     `gorm:"not null;column:limit_cents"`
	SpentCents         int64     `gorm:"not null;default:0;column:spent_cents"`
	Currency           string    `gorm:"not null;size:3;default:BRL"`
	NotifyWhenExceeded bool      `gorm:"not null;defaut:true;column:notify_when_exceeded"`
	CreatedAt          time.Time `gorm:"not null"`
	UpdatedAt          time.Time `gorm:"not null"`

	// Preload da categoria para exibição no dashboard
	Category *CategoryModel `gorm:"foreignKey:CategoryID;constraint:false"`
}

func (BudgetModel) TableName() string { return "budgets" }

// ============================================================
// RecurringTransactionModel
// ============================================================
type RecurringTransactionModel struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	UserID     uuid.UUID  `gorm:"type:uuid;not null;index"`
	AccountID  uuid.UUID  `gorm:"type:uuid;not null"`
	CategoryID *uuid.UUID `gorm:"type:uuid"`

	Description string `gorm:"not null;size:255"`
	AmountCents int64  `gorm:"not null;column:amount_cents"`
	Currency    string `gorm:"not null;size:3;default:BRL"`
	Type        string `gorm:"not null;size:20"`
	Frequency   string `gorm:"not null;size:20"`

	// Ponteiro porque só é relevante para frequência mensal
	DayOfMonth *int       `gorm:"column:day_of_month"`
	StartsOn   time.Time  `gorm:"not null;type:date;column:starts_on"`
	EndsOn     *time.Time `gorm:"type:date;column:ends_on"`

	// Indexado para o worker diário
	NextOccurrence time.Time `gorm:"not null;type:date;column:next_occurrence;index"`

	Active    bool      `gorm:"not null;default:true"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

func (RecurringTransactionModel) TableName() string { return "recurring_transactions" }
