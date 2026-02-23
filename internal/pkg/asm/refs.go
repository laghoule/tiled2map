package asm

import (
	"slices"
)

type ID map[int]int

// NewMap creates a new ID map from a slice of GIDs. It sorts the GIDs and maps each GID to its index in the sorted slice.
func NewMap(gids []int) *ID {
	slices.Sort(gids)
	id := make(ID, len(gids))

	for i, gid := range gids {
		id[gid] = i
	}

	return &id
}
