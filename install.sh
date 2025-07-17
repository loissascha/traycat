#!/bin/bash

mkdir -p ~/.local/bin

chmod +x traycat

killall -9 traycat
cp traycat ~/.local/bin/traycat

echo "traycat has been installed"
