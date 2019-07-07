package kideungeo

import (
	"github.com/lucasb-eyer/go-colorful"
)

func HexToInteger(hex string) int {
	color, _ := colorful.Hex("#0099ff")
	var rgb int = int(color.R)
	rgb = (rgb << 8) + int(color.G)
	rgb = (rgb << 8) + int(color.B)
	return rgb
}
