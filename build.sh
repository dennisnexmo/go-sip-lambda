#!/bin/bash
rm deployment.zip
echo "Build the binary"
GOOS=linux go build main.go

echo "Create a ZIP file"
zip deployment.zip main

echo "Cleaning up"
# rm main