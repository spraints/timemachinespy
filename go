#!/bin/bash

set -ex

# dockerfile based on info in https://superuser.com/questions/306497/can-linux-mount-a-normal-time-machine-sparse-bundle-disk-image-directory
docker build --rm -t tmfs .

# see https://stackoverflow.com/questions/48402218/fuse-inside-docker for more flags i may need to set.
exec docker run -it --rm \
  -v /data/timemachine:/tm:ro \
  --device /dev/fuse \
  --device /dev/loop-control \
  --device /dev/loop0 \
  --device-cgroup-rule="b 7:* rmw" \
  --cap-add SYS_ADMIN \
  tmfs

# see https://markpith.wordpress.com/2016/06/15/restoring-files-from-timemachine-backup-on-synology/
#
# example:
#
# $ sparsebundlefs /tm/Slick.sparsebundle /mnt/sparse
#
# $ parted /mnt/sparse/sparsebundle.dmg unit B print
# ...
# 
# Number  Start       End             Size            File system  Name       Flags
#  1      20480B      209735679B      209715200B      fat32        EFI System Partition  boot, esp
#  2      209735680B  6815609761791B  6815400026112B  hfsx         disk image
# 
# $ losetup -f /mnt/sparse/sparsebundle.dmg --offset 209735680 --sizelimit 6815400026112 --show
#
# $ mount -t hfsplus /dev/loop0 /mnt/hfs
#
# $ tmfs /mnt/hfs /mnt/tmfs
