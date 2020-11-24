package main

import (
	"fmt"

	"github.com/xycui/go-playground/module-init/module-init/m"
	"github.com/xycui/go-playground/module-init/module-init/n"
)

var (
	str = initstr()
)

func initstr() string {
	fmt.Println("var init in main")
	return ""
}

func init() {
	fmt.Println("module init in main")
}

func main() {
	n.DoNothing()
	m.DoNothing()
}
