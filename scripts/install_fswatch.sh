#!/bin/bash
echo "this script will download/configure/compile/install fswatch Release 1.16.0 on your Linux Box"
cd ~
wget https://github.com/emcrisostomo/fswatch/releases/download/1.16.0/fswatch-1.16.0.tar.gz
tar xvfz fswatch-1.16.0.tar.gz
cd fswatch-1.16.0/
./configure
make
sudo make install
echo "if you are here without error you can start using fswatch "
