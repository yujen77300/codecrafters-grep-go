package models

type patternItemType int

const (
	LiteralType patternItemType = iota
	DigitType
	WordCharType
	NegatedCharSetType
	CharSetType
)

type PatternItem struct {
	ItemType patternItemType
	Value    byte
	CharSet  string
}
