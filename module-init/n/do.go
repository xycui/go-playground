package n

import "fmt"

var (
	str = initstr()
)

func initstr() string {
	fmt.Println("var init in 'n'")
	return ""
}

func init() {
	fmt.Println("module init in 'n'")
}

func DoNothing() {

}
