# Grobuzin agent

This repo contains all the code that runs in the microVM and the image builder for firecracker.

## Building an image 

Run image builder, you need to set "AGENT_REPO_FOLDER" to point to the root of this repository, it will build the `main-agent` docker image and `nodejs` agent then put everything in `rootfs.ext4` in the current working directory. 