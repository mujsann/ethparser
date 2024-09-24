package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	parser "mujsann.com/ethparser/pkg"
)

type App struct {
	svc parser.Parser
}

func main() {
	loadEnv()

	Txlimit, _ := strconv.Atoi(os.Getenv("TRANSACTION_DAYS_LIMIT"))
	parserService := parser.NewParser(int64(Txlimit))

	app := &App{
		svc: &parserService,
	}

	PORT, err := strconv.Atoi(os.Getenv("PORT"))
	if PORT == 0 || err != nil {
		log.Println("Port is not set in the env file; using 8080")
		PORT = 8080
	}

	defineRoutes(app)

	fmt.Printf("Serving ethparser on port %d...", PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil))
}

// load env variables
func loadEnv() {
	file, err := os.Open(".env")
	if err != nil {
		log.Printf("could not open load environmental variables")
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// skip empty lines and comments
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			os.Setenv(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])) // Set the environment variable
		}
	}
}
