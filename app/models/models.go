package models

type patternItemType int

const (
	LiteralType patternItemType = iota
	DigitType
	WordCharType
	NegatedCharSetType
	CharSetType
	StartOfLineType
	EndOfLineType
)

type PatternItem struct {
	ItemType patternItemType
	Value    byte
	CharSet  string
}
