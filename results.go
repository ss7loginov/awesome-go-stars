package main

import (
	"sort"
	"strconv"
	"sync"
)

// Results struct that contains info on fetched repos
type Results struct {
	sync.RWMutex
	m       map[int][]string
	fetched int
	total   int
}

// groupedRepos returns repos grouped by the number of stars and its keys sorted in a descending order
func (r *Results) groupedRepos() (groupedRepos map[string][]string, sortedKeys []string) {
	groupedRepos = make(map[string][]string)

	r.Lock()
	defer r.Unlock()

	var keys []int
	for k := range r.m {
		keys = append(keys, k)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(keys)))
	lastKey := ""
	for _, k := range keys {
		group := starsGroup(k)
		if group != lastKey {
			sortedKeys = append(sortedKeys, group)
			lastKey = group
		}
		groupedRepos[group] = append(groupedRepos[group], r.m[k]...)
	}
	return
}

// starsGroup returns a string for a group
func starsGroup(stars int) string {
	switch {
	case stars == 0:
		return "0"
	case stars < 5:
		return "1+"
	case stars < 10:
		return "5+"
	case stars < 100:
		return strconv.Itoa(stars/10*10) + "+"
	case stars < 1000:
		return strconv.Itoa(stars/50*50) + "+"
	case stars < 5000:
		return strconv.Itoa(stars/500*500) + "+"
	case stars < 10000:
		return strconv.Itoa(stars/1000*1000) + "+"
	default:
		return strconv.Itoa(stars/5000*5000) + "+"
	}
}
