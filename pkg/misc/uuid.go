package misc

import (
	"github.com/google/uuid"
)

// UUID генерирует новый идентификатор в виде строки
func UUID() string {
	return uuid.New().String()
}
