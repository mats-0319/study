package utils

import "slices"

type Set struct {
	Data []string
}

func (s *Set) Add(value string) *Set {
	if s == nil || s.Data == nil { // no initialize
		return &Set{
			Data: []string{value},
		}
	}

	if !slices.Contains(s.Data, value) {
		s.Data = append(s.Data, value)
	}

	return s
}
