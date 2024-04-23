#!/bin/sh
if [ ! -d linux.git ]
then
git clone --depth 1 --branch v6.8 https://git.kernel.org/pub/scm/linux/kernel/git/torvalds/linux.git linux.git
fi

pushd linux.git

cp ../kernel.config .config

make vmlinux -j `nproc`

popd

cp linux.git/vmlinux .

echo "Kernel build complete, vmlinux is ready"
 