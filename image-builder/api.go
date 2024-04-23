package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	AgentRepoFolder string
	DockerClient    *Docker
	StorageService  *StorageService
	DB              *sql.DB
}

func startApi(agentRepoFolder string, dockerClient *Docker, storageService *StorageService, db *sql.DB) {

	controller := Controller{AgentRepoFolder: agentRepoFolder, DockerClient: dockerClient, StorageService: storageService, DB: db}

	router := gin.Default()

	router.POST("/build", controller.build)

	err := router.Run()

	if err != nil {
		panic(err)
	}

}

type buildBody struct {
	Id      string            `json:"id"`
	Variant string            `json:"variant"`
	Files   map[string]string `json:"files"`
}

func (controller *Controller) build(c *gin.Context) {

	var body buildBody

	err := c.BindJSON(&body)

	if err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid body",
		})
		return
	}

	variant := body.Variant

	// check if the variant is valid

	if _, err := os.Stat(path.Join(controller.AgentRepoFolder, variant)); os.IsNotExist(err) {
		c.JSON(400, gin.H{
			"error": "Invalid variant",
		})
		return
	}

	userCodeFolder := path.Join(controller.AgentRepoFolder, "user-code")

	err = os.RemoveAll(userCodeFolder)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": "Filesystem error",
		})
		return
	}

	err = os.Mkdir(userCodeFolder, 0755)

	if err != nil {
		log.Println(err)

		c.JSON(500, gin.H{
			"error": "Filesystem error",
		})
		return
	}

	go buildForked(body, controller.AgentRepoFolder, controller.DockerClient, controller.StorageService, controller.DB)

	c.JSON(200, gin.H{
		"message": "Image build started",
	})

}

func buildForked(body buildBody, agentRepoFolder string, dockerClient *Docker, storageService *StorageService, db *sql.DB) {
	userCodeFolder := path.Join(agentRepoFolder, "user-code")
	variant := body.Variant

	var err error

	// Insert the files

	for key, value := range body.Files {

		folder := path.Dir(key)

		err = os.MkdirAll(path.Join(userCodeFolder, folder), 0755)

		if err != nil {
			log.Println("err : ", err)
			return
		}

		filePath := path.Join(userCodeFolder, key)

		file, err := os.Create(filePath)

		if err != nil {
			log.Println("err : ", err)
			return
		}

		_, err = file.WriteString(value)

		if err != nil {
			log.Println("err : ", err)
			return
		}
	}

	// build the image

	var logs string
	logs, err = dockerClient.buildImage(agentRepoFolder, []string{"user-code", variant}, variant+"/Dockerfile", []string{"grobuzin/nodejs-agent:latest"})

	if err != nil {
		log.Println("err:", err, logs)
		return
	}

	rootfLocation := "/tmp/rootfs.ext4"

	// Create a 1GB ext4 filesystem

	err = createRootfs(rootfLocation, 1024*1024*1024)

	if err != nil {
		log.Println("err : ", err)
		return
	}

	err = dockerClient.copyToRootfs("grobuzin/nodejs-agent:latest", rootfLocation, agentRepoFolder)

	if err != nil {
		log.Println("err : ", err)
		return
	}

	objectLocation := fmt.Sprintf("function/%s/rootfs.ext4", body.Id)

	println("Uploading file to ", objectLocation)

	err = storageService.UploadFile(objectLocation, rootfLocation)

	if err != nil {
		log.Println("upload err : ", err)
		return
	}

	// set as ready in the database
	err = SetFunctionReady(db, body.Id)

	if err != nil {
		log.Println("err : ", err)
		return
	}

	log.Println("Image build complete")
}
