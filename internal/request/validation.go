package request

import (
	"strings"

	"github.com/jaeyoung0509/titanbay-funds-api/internal/domain/vo"
)

func IsUUID(raw string) bool {
	_, err := vo.ParseID(raw)
	return err == nil
}

func IsEmail(raw string) bool {
	_, err := vo.NewEmail(strings.TrimSpace(raw))
	return err == nil
}
