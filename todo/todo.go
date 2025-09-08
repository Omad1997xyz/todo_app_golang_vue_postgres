package todo

import (
	"errors"
	"fmt"
)

type Vazifa struct {
	Nomi   string
	Holati bool
}

func YangiVazifa(nom string) Vazifa {
	return Vazifa{
		Nomi:   nom,
		Holati: false,
	}
}

type Royxat struct {
	Vazifalar []Vazifa
}

func Yangiroyxat() *Royxat {
	return &Royxat{
		Vazifalar: []Vazifa{},
	}
}

func (l *Royxat) Qoshish(istalgan_ish string) {
	l.Vazifalar = append(l.Vazifalar, YangiVazifa(istalgan_ish))
}

func (l *Royxat) Bajarildi(tr int) error {
	if tr < 0 || tr >= len(l.Vazifalar) {
		return errors.New("Bunaqa tartib raqami mavjud emas !")
	}

	l.Vazifalar[tr].Holati = true
	return nil
}

func (l *Royxat) String() string {
	if len(l.Vazifalar) == 0 {
		errors.New("Hech qanday vazifa kiritilmagan !")
	}
	result := "Todo List:\n"

	for i, v := range l.Vazifalar {
		status := " "
		if v.Holati {
			status = "✓"
		}
		result += fmt.Sprintf("%d. [%s] %v\n", i+1, status, v)
	}
	return result
}
