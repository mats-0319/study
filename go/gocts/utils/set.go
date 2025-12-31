package utils

type Set struct {
	Data []string
}

func (s *Set) Add(value string) *Set {
	if s == nil || s.Data == nil { // no initialize
		return &Set{
			Data: []string{value},
		}
	}

	isExist := false
	for i := range s.Data {
		if s.Data[i] == value {
			isExist = true
			break
		}
	}

	if !isExist {
		s.Data = append(s.Data, value)
	}

	return s
}
