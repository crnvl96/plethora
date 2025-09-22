package ideas

import (
	"fmt"
)

func birthdayParadox() {
	fmt.Println("Birthday paradox!")
}

func init() {
	Ideas["birthdayParadox"] = birthdayParadox
}
