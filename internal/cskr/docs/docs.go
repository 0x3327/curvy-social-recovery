package main

import (
	"os"

	"github.com/0x3327/curvy-social-recovery/internal/cskr"
	docs "github.com/urfave/cli-docs/v3"
)

func main() {
	app := cskr.NewCli()

	md, err := docs.ToMarkdown(app)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("./docs.md", []byte(md), 0666)
	if err != nil {
		panic(err)
	}
}
