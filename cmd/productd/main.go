package main

import (
	_ "github.com/lib/pq"
	"wiliam.dev/product/cmd/productd/cmd"
)

func main() {
	cmd.Execute()
}
