package main

import (
	"fmt"
	hw09structvalidator "github.com/Grog2903/hw/hw09_struct_validator"
	"runtime/debug"
)

type App struct {
	Version string `validate:"len:5"`
}

func main() {
	app := App{
		Version: "0.0.12",
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Panic occurred: %v\n", r)
			debug.PrintStack() // Вывод стека вызовов
		}
	}()

	validate, err := hw09structvalidator.Validate(app)
	fmt.Println(validate, err)
}
