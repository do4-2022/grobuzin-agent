package main

import (
	"fmt"
	"log"
	"os"

	"github.com/docker/docker/client"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	agentRepoFolder := os.Getenv("AGENT_REPO_FOLDER")
	if agentRepoFolder == "" {
		agentRepoFolder = "../"
	}
	storageEndpoint := os.Getenv("MINIO_ENDPOINT")
	storageAccessKey := os.Getenv("MINIO_ACCESS_KEY")
	storageSecretKey := os.Getenv("MINIO_SECRET_KEY")
	storageSecure := os.Getenv("MINIO_SECURE")
	postgresConnectionStr := os.Getenv("POSTGRES_CONNECTION_STR")

	db, err := StartDBConnection(postgresConnectionStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	apiClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	defer apiClient.Close()

	docker := Docker{client: apiClient}

	storageService, err := StartStorageConnection(storageEndpoint, storageAccessKey, storageSecretKey, storageSecure)

	if err != nil {
		panic(err)
	}

	// Build the main-agent image when starting

	var logs string
	logs, err = docker.buildImage(agentRepoFolder, []string{"main-agent"}, "main-agent/Dockerfile", []string{"grobuzin/main-agent:latest"})

	if err != nil {
		fmt.Println(err, logs)
		os.Exit(1)
	}

	startApi(agentRepoFolder, &docker, &storageService, db)

}
