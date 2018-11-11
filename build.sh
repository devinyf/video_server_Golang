#!/bin/bash

# Build web UI
cd /Users/squall/go/src/myproject/video_server/web
go install
cp ~/go/bin/web ~/go/bin/video_server_ui/web
cp -R ~/go/src/myproject/video_server/templates ~/go/bin/video_server_ui/

