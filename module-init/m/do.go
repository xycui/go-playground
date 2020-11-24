package m

import "fmt"

var (
	str = initstr()
)

func initstr() string {
	fmt.Println("var init in 'm'")
	return ""
}

func init() {
	fmt.Println("module init in 'm'")
}

func DoNothing() {

}
