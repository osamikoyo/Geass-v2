package main

import (
	"fmt"
	"github.com/osamikoyo/geass-v2/internal/app"
)

func main()  {
	if err := app.App();err != nil{
		fmt.Println(err)
	}
}