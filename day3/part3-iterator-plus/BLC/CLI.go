package BLC

import (
	"fmt"
	"os"
	"flag"
	"log"
)

type CLI struct {

}

func Create() {
	fmt.Println("create")
}

func (cli CLI) Run() {
	fs := flag.NewFlagSet("hahaha", flag.ExitOnError)

	args := os.Args

	switch args[1] {
		case "hahaha":
		err := fs.Parse(args[2:])
		if err != nil {
			log.Panic(err)
		}
	}
}