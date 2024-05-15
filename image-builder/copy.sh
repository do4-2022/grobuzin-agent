#!/bin/sh

# Need to mount the rootfs image to /my-rootfs

apk add openrc
apk add util-linux

# Copy the agent startup script in openrc folder
cp /main-agent/openrc/agent /etc/init.d/agent

# Set up a login terminal on the serial console (ttyS0):
ln -s agetty /etc/init.d/agetty.ttyS0
echo ttyS0 > /etc/securetty

echo auto lo >> /etc/network/interfaces
echo auto eth0 >> /etc/network/interfaces

mkdir -p 

rc-update add agetty.ttyS0 default

# Make sure special file systems are mounted on boot:
rc-update add devfs boot
rc-update add procfs boot
rc-update add sysfs boot

# start the agent at boot 
rc-update add agent default

# Then, copy the newly configured system to the rootfs image:
for d in bin etc lib root sbin usr app; do tar c "/$d" | tar x -C /mnt/my-rootfs; done

# The above command may trigger the following message:
# tar: Removing leading "/" from member names
# However, this is just a warning, so you should be able to
# proceed with the setup process.

for dir in dev proc run sys var "lib/modules" "var/run" ; do mkdir -p /mnt/my-rootfs/${dir}; done

# All done, exit docker shell.
exit