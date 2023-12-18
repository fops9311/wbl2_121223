package unpack

import "errors"

var ErrInvalidPattern error = errors.New("invalid pattern")
var ErrEmptyString error = errors.New("empty string")

func Unpack(s string) (string, error) {
	head := AsLinkedListOfRunes(s)
	if head == nil {
		return "", ErrEmptyString
	}
	result := make([]rune, 0, len(s))
	multiplier := 0
	addToResult := func(r rune, excape bool) error {
		if excape {
			if multiplier == 0 {
				multiplier = 1
			}
			for i := 0; i < multiplier; i++ {
				result = append(result, r)
			}
			multiplier = 0
			return nil
		}
		switch r {
		case '0':
			if multiplier == 0 {
				return ErrInvalidPattern
			}
			multiplier = multiplier * 10
			return nil
		case '1':
			multiplier = multiplier*10 + 1
			return nil
		case '2':
			multiplier = multiplier*10 + 2
			return nil
		case '3':
			multiplier = multiplier*10 + 3
			return nil
		case '4':
			multiplier = multiplier*10 + 4
			return nil
		case '5':
			multiplier = multiplier*10 + 5
			return nil
		case '6':
			multiplier = multiplier*10 + 6
			return nil
		case '7':
			multiplier = multiplier*10 + 7
			return nil
		case '8':
			multiplier = multiplier*10 + 8
			return nil
		case '9':
			multiplier = multiplier*10 + 9
			return nil
		default:
			if multiplier == 0 {
				multiplier = 1
			}
			for i := 0; i < multiplier; i++ {
				result = append(result, r)
			}
			multiplier = 0
			return nil
		}
	}
	escapeflag := false
	excapesimbol := '\\'
	head.Iterate(func(s *Simbol) error {
		if escapeflag {
			err := addToResult(s.Val, escapeflag)
			escapeflag = false
			return err
		}
		if (s.Val) == rune(excapesimbol) {
			escapeflag = true
			return nil
		}
		err := addToResult(s.Val, escapeflag)
		escapeflag = false
		return err
	})
	if len(result) == 0 {
		return "", ErrInvalidPattern
	}
	return string(result), nil
}

type Simbol struct {
	Next *Simbol
	Val  rune
}

func (sim *Simbol) Iterate(f func(s *Simbol) error) error {
	err := f(sim)
	if err != nil {
		return err
	}
	if sim.Next != nil {
		return sim.Next.Iterate(f)
	}
	return nil
}

func AsLinkedListOfRunes(s string) *Simbol {
	runes := []rune(s)
	if len(runes) < 1 {
		return nil
	}
	head := &Simbol{
		Val: rune(runes[0]),
	}
	var prev *Simbol = head
	for i := 1; i < len(runes); i++ {
		var curr *Simbol = &Simbol{
			Val: runes[i],
		}
		prev.Next = curr
		prev = curr
	}
	return head
}
