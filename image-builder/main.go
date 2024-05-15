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
		log.Println("Error loading .env file")
	}

	apiClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	defer apiClient.Close()

	docker := Docker{client: apiClient}

	agentRepoFolder := os.Getenv("AGENT_REPO_FOLDER")

	if agentRepoFolder == "" {
		agentRepoFolder = "../"
	}

	args := os.Args

	// if the are arguments, we are doing a custom run
	if len(args) > 1 {

		variant := args[1]

		err = manualBuild(docker, agentRepoFolder, variant)
		if err != nil {
			log.Println("err : ", err)
		}
		return
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

func manualBuild(docker Docker, agentRepoFolder string, variant string) (err error) {

	// Build the main-agent image when starting

	var logs string
	logs, err = docker.buildImage(agentRepoFolder, []string{"main-agent"}, "main-agent/Dockerfile", []string{"grobuzin/main-agent:latest"})

	if err != nil {
		fmt.Println(err, logs)
		os.Exit(1)
	}

	image_name := fmt.Sprintf("grobuzin/%s-%s:latest", variant, "aaaaa")

	logs, err = docker.buildImage(agentRepoFolder, []string{"user-code", variant}, variant+"/Dockerfile", []string{image_name})

	if err != nil {
		log.Println("err:", err, logs)
		return
	}

	rootfLocation := "./rootfs.ext4"

	// Create a 1GB ext4 filesystem

	err = createRootfs(rootfLocation, 1024*1024*1024)

	if err != nil {
		log.Println("err : ", err)
		return
	}

	err = docker.copyToRootfs(image_name, rootfLocation, agentRepoFolder)

	if err != nil {
		log.Println("err : ", err)
		return
	}
	return
}
