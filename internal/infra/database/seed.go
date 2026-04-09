package database

import (
	"fmt"

	"gorm.io/gorm"
)

// ==========================================================
// seedCategoryModel cria categorias pré definidas do sistema.
// ==========================================================

type seedCategoryModel struct {
	Name        string
	DefaultType string
}

// seedCategories cria categorias pré definidas do sistema
func SeedCategories(db *gorm.DB) error {
	categories := []seedCategoryModel{

		// Receitas
		{Name: "Salário", DefaultType: "income"},
		{Name: "Freelance", DefaultType: "income"},
		{Name: "Investimentos", DefaultType: "income"},
		{Name: "Presente", DefaultType: "income"},
		{Name: "Outras receitas", DefaultType: "income"},

		// Despesas
		{Name: "Alimentação", DefaultType: "expense"},
		{Name: "Moradia", DefaultType: "expense"},
		{Name: "Transporte", DefaultType: "expense"},
		{Name: "Saúde", DefaultType: "expense"},
		{Name: "Educação", DefaultType: "expense"},
		{Name: "Lazer", DefaultType: "expense"},
		{Name: "Vestuário", DefaultType: "expense"},
		{Name: "Assinaturas", DefaultType: "expense"},
		{Name: "Pet", DefaultType: "expense"},
		{Name: "Outras despesas", DefaultType: "expense"},
	}

	for _, cat := range categories {
		result := db.Exec(`
			INSERT INTO categories (id, users_id, name, default_type, active, created_at, updated_at)
			VALUES (gen_random_uuid(), NULL, ?, NULLIF(?, '')::transaction_type, true, NOW(), NOW())
			ON CONFLICT (user_id, name) DO NOTHING
		`, cat.Name, cat.DefaultType)

		if result.Error != nil {
			return fmt.Errorf("SeedCategories %q: %w", cat.Name, result.Error)
		}
	}

	return nil
}
