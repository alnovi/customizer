package main

import "alnovi/customizer/internal/customizer"

func main() {
	config := initConfig()

	app := customizer.NewApp(config)
	app.Run()
}
