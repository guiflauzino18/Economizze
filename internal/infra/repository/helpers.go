package repository

func normalizePaginator(page, size int) (int, int) {
	// Página míminma é 1
	if page < 1 {
		page = 1
	}

	// Tamanho padrão é 20 e máximo é 100
	if size < 0 {
		size = 20
	}

	if size > 100 {
		size = 100
	}

	return page, size
}
