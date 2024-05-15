# Firecracker setup

Testing locally

## Build vmlinux

```sh
cd kernel
./build.sh
```

## Build rootfs

```sh
cd image-builder
go build .
sudo ./image-builder nodejs
```

By default there is no password set for root, to set it you need to mount the rootfs and change the password :

```sh
mkdir -p rootfs
sudo mount rootfs.ext4 rootfs
sudo chroot rootfs /bin/sh

# we're now in a shell in the image

passwd
# set your new password...

exit
sudo umount rootfs
```

## Network setup on the host

```sh
sudo ip tuntap add tap0 mode tap
sudo ip link set tap0 up
sudo ip addr add 192.168.0.1/24 dev tap0
```

Change `eth0` by the name of the interface on the host that is providing internet.

```sh
sudo iptables -t nat -A POSTROUTING -s 192.168.0.0/24 -o eth0 -j MASQUERADE
sudo iptables -A FORWARD -m conntrack --ctstate RELATED,ESTABLISHED -j ACCEPT
sudo iptables -A FORWARD -s 192.168.0.0/24 -j ACCEPT
```

## Start firecracker in a terminal

```sh
sudo rm -rf /tmp/firecracker.socket && sudo firecracker --api-sock /tmp/firecracker.socket
```

## Start the VM

Set ROOTFS_PATH and KERNEL_PATH in `setup-firecracker.sh` to point to the generated files and then run the script :

```sh
sudo ./setup-firecracker.sh
```

## Then setup the network on the guest

```sh
ip link set dev eth0 up
ip addr add 192.168.0.2/24 dev eth0
ip route add default via 192.168.0.1
```
