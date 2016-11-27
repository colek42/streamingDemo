#!/bin/bash

set -e
set -x

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )"
PREFIX_DIR="$DIR/shared"
BIN_DIR="$DIR/bin"

rm -rf $PREFIX_DIR

# Build libvpx
cd $DIR/lib/libvpx
make distclean || true
./configure --prefix="$PREFIX_DIR" --disable-examples
time make -j16
make install
make clean

# Build x264
cd $DIR/lib/x264
make distclean || true
./configure --prefix="$PREFIX_DIR" --bindir="$BIN_DIR" --enable-static
time make -j16
make install
make distclean

cd $DIR/lib/FFmpeg
make distclean || true
PKG_CONFIG_PATH="$PREFIX_DIR/lib/pkgconfig" ./configure \
    --prefix="$PREFIX_DIR" \
    --bindir="$BIN_DIR" \
    --enable-gpl \
    --enable-static \
    --disable-shared \
    --enable-libx264 \
    --enable-nonfree \
    --enable-libvpx
time make -j16
make install
make distclean