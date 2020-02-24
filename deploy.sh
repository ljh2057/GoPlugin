#!/bin/sh
toolName=CompassChecker
SHELL_FOLDER=$(cd "$(dirname "$0")";pwd)
if [ ! -d "$toolName" ]; then
  mkdir $toolName
fi
curl -sL -o $toolName/config.json https://github.com/ljh2057/GoPlugin/releases/download/v0.0.2/config.json
curl -LJ https://github.com/ljh2057/GoPlugin/releases/download/v0.0.2/GoPlugin_0.0.2_Linux_x86_64.tar.gz | tar -zx -C $toolName
chmod +x $toolName/GoPlugin
cd $SHELL_FOLDER/$toolName
./GoPlugin
