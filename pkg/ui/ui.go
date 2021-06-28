// Package ui provides a set of components to organize an application layout.
package ui

import (
	"strconv"
)

var (
	// The padding of block-like components in px.
	BlockPadding = 30

	// The padding of block-like components in px when app width is <= 480px.
	BlockMobilePadding = 18

	// The horizontal padding of base-like components in px.
	BaseHPadding = 36

	// The horizontal padding of base-like components in px when app width is <= 480px.
	BaseMobileHPadding = 12

	// The vertical padding of base-like components in px.
	BaseVPadding = 12
)

func pxToString(px int) string {
	return strconv.Itoa(px) + "px"
}

type alignment int

const (
	top alignment = iota
	right
	bottom
	left
	middle
)

type style struct {
	key   string
	value string
}
