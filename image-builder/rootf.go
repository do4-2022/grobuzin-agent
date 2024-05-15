package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
)

type Docker struct {
	client *client.Client
}

func (d *Docker) buildImage(buildContextFolder string, include_files []string, dockerfilePath string, tags []string) (logs string, err error) {

	buildOptions := types.ImageBuildOptions{
		Tags:       tags,
		Dockerfile: dockerfilePath,
		Remove:     true,
	}

	archive, err := archive.TarWithOptions(buildContextFolder, &archive.TarOptions{
		IncludeFiles: include_files,
	})
	if err != nil {
		return
	}

	resp, err := d.client.ImageBuild(context.Background(), archive, buildOptions)

	if err != nil {
		return
	}

	defer resp.Body.Close()

	buffer := bytes.NewBuffer(nil)

	io.Copy(buffer, resp.Body)

	logs = buffer.String()

	log.Println(logs)

	return
}

func createRootfs(filename string, size int) (err error) {

	// Create the file
	file, err := os.Create(filename)
	if err != nil {
		return
	}
	defer file.Close()

	// Set the file size by writing a single zero byte at the last position
	_, err = file.WriteAt([]byte{0}, int64(size-1))
	if err != nil {
		return
	}

	// Close the file as mkfs.ext4 needs to be run on a closed file
	file.Close()

	// Run mkfs.ext4 to create an ext4 filesystem in the file
	cmd := exec.Command("mkfs.ext4", "-F", filename)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
	}

	return
}

func (d *Docker) copyToRootfs(image string, rootfs string, agent_repo_folder string) (err error) {

	mountPath := "/tmp/rootfs"

	_, err = os.Stat(mountPath)

	if err != nil && !os.IsNotExist(err) {
		return
	}
	if !os.IsNotExist(err) {
		os.RemoveAll(mountPath)
	}

	// mount the image file into the folder

	err = os.MkdirAll(mountPath, 0755)

	if err != nil {
		return
	}

	out, err := exec.Command("mount", rootfs, mountPath).CombinedOutput()

	if err != nil {
		fmt.Println(string(out))
		return
	}

	defer exec.Command("umount", mountPath).Run()

	// mount the folder into the container
	hostConfig := &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: mountPath,
				Target: "/mnt/my-rootfs",
			},
		},
	}

	containerConfig := &container.Config{
		Image: image,
		Cmd:   []string{"sh", "-c", "/image-builder/copy.sh"},
	}

	containerResp, err := d.client.ContainerCreate(context.Background(), containerConfig, hostConfig, nil, nil, "")

	if err != nil {
		return
	}

	containerID := containerResp.ID

	// copy ./copy.sh script into the container

	copyScript, err := archive.TarWithOptions(agent_repo_folder, &archive.TarOptions{
		IncludeFiles: []string{"image-builder/copy.sh", "main-agent/openrc/agent"},
	})

	if err != nil {
		return
	}

	err = d.client.CopyToContainer(context.Background(), containerID, "/", copyScript, types.CopyToContainerOptions{})

	if err != nil {
		return
	}

	err = d.client.ContainerStart(context.Background(), containerID, container.StartOptions{})

	log.Println("Container started", containerID)

	if err != nil {
		return
	}

	// wait for the container to be started
	statusCh, errCh := d.client.ContainerWait(context.Background(), containerID, container.WaitConditionNotRunning)

	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
	case a := <-statusCh:

		log.Println("Container finished", a.StatusCode)
	}

	return

}
