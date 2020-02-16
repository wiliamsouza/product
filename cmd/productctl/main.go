package main

import (
	_ "github.com/lib/pq"
	"wiliam.dev/product/cmd/productctl/cmd"
)

func main() {
	cmd.Execute()
}
