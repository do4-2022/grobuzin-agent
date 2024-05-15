#!/bin/sh

# run firecracker with 
# sudo firecracker --api-sock /tmp/firecracker.socket

# need to run this script as root

ROOTFS_PATH=/home/sautax/git/grobuzin-agent/image-builder/rootfs.ext4
KERNEL_PATH=/home/sautax/git/grobuzin-agent/kernel/vmlinux

curl -X PUT --unix-socket /tmp/firecracker.socket -i \
   --data  '{
          "vcpu_count": 2,
          "mem_size_mib": 1024
          }' "http://localhost/machine-config"

curl --unix-socket /tmp/firecracker.socket -i \
    -X PUT "http://localhost/drives/rootfs" \
    -H "Accept: application/json" \
    -H "Content-Type: application/json" \
    -d '{
          "drive_id": "rootfs",
          "path_on_host": "'$ROOTFS_PATH'",
          "is_root_device": true,
          "is_read_only": false
        }'

 curl --unix-socket /tmp/firecracker.socket -i \
    -X PUT "http://localhost/boot-source" \
    -H "Accept: application/json" \
    -H "Content-Type: application/json" \
    -d '{
          "kernel_image_path": "'$KERNEL_PATH'",
          "boot_args": "console=ttyS0 reboot=k panic=1 pci=off"
        }'

# network

curl --unix-socket /tmp/firecracker.socket -i \
    -X PUT "http://localhost/network-interfaces/eth0" \
    -H "accept: application/json" \
    -H "Content-Type: application/json" \
    -d '{
        "iface_id": "eth0",
        "guest_mac": "AA:FC:00:00:00:01",
        "host_dev_name": "tap0"
    }'

# start the microVM
curl --unix-socket /tmp/firecracker.socket -i \
    -X PUT "http://localhost/actions" \
    -H "Accept: application/json" \
    -H "Content-Type: application/json" \
    -d '{
          "action_type": "InstanceStart"
        }'


