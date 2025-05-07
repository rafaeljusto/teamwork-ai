package actions

import (
	"cmp"
	"slices"
	"sort"

	"github.com/rafaeljusto/teamwork-ai/internal/teamwork"
)

// intersection returns the intersection of two slices. It will sort the slices
// before performing the intersection. The slices must contain elements that
// implement the comparable and cmp.Ordered interfaces.
func intersection[T interface {
	comparable
	cmp.Ordered
}](a []T, b []T) []T {
	slices.Sort(a)
	slices.Sort(b)

	set := make([]T, 0)
	for _, v := range a {
		idx := sort.Search(len(b), func(i int) bool {
			return b[i] == v
		})
		if idx < len(b) && b[idx] == v {
			set = append(set, v)
		}
	}
	return set
}

// extractMappedIDs ensure that only IDs from relationships mapped to the source
// are returned.
func extractMappedIDs[T any](relationships []teamwork.Relationship, sourceMap map[int64]T) []int64 {
	result := make([]int64, 0, len(relationships))
	for _, relationship := range relationships {
		_, ok := sourceMap[relationship.ID]
		if !ok {
			continue
		}
		result = append(result, relationship.ID)
	}
	return result
}
