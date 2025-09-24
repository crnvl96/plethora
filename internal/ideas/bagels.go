package ideas

import (
	"fmt"
)

func bagels() {
	fmt.Println("Bagels!")
}

func init() {
	Ideas["bagels"] = bagels
}
