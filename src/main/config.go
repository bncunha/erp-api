package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DB_PASS string
	DB_HOST string
	DB_PORT string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	listFiles()
	if err != nil {
		return nil, err
	}

	return &Config{
		DB_PASS: os.Getenv("DB_PASS"),
		DB_HOST: os.Getenv("DB_HOST"),
		DB_PORT: os.Getenv("DB_PORT"),
	}, nil
}

func listFiles() {
	// Listar arquivos em um diretório específico
	diretorio := "." // Diretório atual
	arquivos, err := ioutil.ReadDir(diretorio)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Arquivos no diretório:", diretorio)
	for _, arquivo := range arquivos {
		fmt.Println(arquivo.Name(), arquivo.IsDir())
	}

	// Listar arquivos recursivamente
	fmt.Println("\nArquivos recursivamente (usando filepath.Walk):")
	filepath.Walk(diretorio, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !strings.Contains(path, ".git") {
		fmt.Println(path)
		}
		return nil
	})
}