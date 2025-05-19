package utilis

import "github.com/codecrafters-io/grep-starter-go/app/models"

func MatchLength(pattern []models.PatternItem) int {
	length := 0
	for _, item := range pattern {
		switch item.ItemType {
		case models.LiteralType, models.DigitType, models.WordCharType, models.CharSetType, models.NegatedCharSetType:
			length++
		}
	}
	return length
}
