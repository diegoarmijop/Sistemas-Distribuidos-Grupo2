package utils

import (
	"fmt"
	"time"
)

// utils/time.go
func formatearTiempoRelativo(t time.Time) string {
	ahora := time.Now()
	diff := ahora.Sub(t)

	switch {
	case diff < time.Minute:
		return "Hace unos segundos"
	case diff < time.Hour:
		minutos := int(diff.Minutes())
		return fmt.Sprintf("Hace %d minutos", minutos)
	case diff < 24*time.Hour:
		horas := int(diff.Hours())
		return fmt.Sprintf("Hace %d horas", horas)
	default:
		dias := int(diff.Hours() / 24)
		return fmt.Sprintf("Hace %d días", dias)
	}
}
