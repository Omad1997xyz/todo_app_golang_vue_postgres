package main

import (
	"fmt"
	todo "github.com/Omad1997xyz/todo/internal/todo/interface"
)

func main() {
	yangiRoyxat := todo.Yangiroyxat()

	yangiRoyxat.Qoshish("Orudiya filmini ko'rish")
	yangiRoyxat.Qoshish("Nimadurlar qilib tashlash")
	yangiRoyxat.Bajarildi(1)

	for i, v := range yangiRoyxat.Vazifalar {
		fmt.Printf("%d. %v\n", i+1, v)
	}
}
