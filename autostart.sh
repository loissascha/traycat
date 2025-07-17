#!/bin/bash

chmod +x install.sh
./install.sh

mkdir -p ~/.config/autostart

echo "[Desktop Entry]
Type=Application
Name=Cat
Exec=/home/$USER/.local/bin/traycat
Comment=Shows a beautiful cat" > ~/.config/autostart/traycat.desktop

echo "autostart added. Log out and back in to see it in action!"
