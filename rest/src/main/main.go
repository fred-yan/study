package main

import "rest"

func main() {
	a := rest.App{}
	a.Initialize(
		"root",
		"redis",
		"testdb")

	a.Run(":8080")
}
