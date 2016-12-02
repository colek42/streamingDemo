#!/bin/bash
# FFMPEG multicast video generation script

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )"
source $DIR/dev/dev.env

OUT_URI=${1:-"udp://@234.5.5.5:8209"}

echo "If you don't have libx264 support, please re-run /build/makeffmpeg.sh"
echo "Generating H.264 video on multicast video on ${OUT_URI}..."

ffmpeg \
    -v info \
    -re \
    -f lavfi \
    -i "testsrc=size=1280x720:rate=ntsc" \
    -f lavfi \
    -i "anoisesrc=c=pink" \
    -c:v libx264 \
    -pix_fmt yuv420p \
    -c:a aac \
    -f mpegts \
    "${OUT_URI}"
