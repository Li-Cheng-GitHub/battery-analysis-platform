#!/bin/bash

ABSOLUTE_PROJECT_PATH="$(cd "$(dirname $0)/..";pwd)"
APP_PATH="${ABSOLUTE_PROJECT_PATH}/app"
BUILD_PATH="${ABSOLUTE_PROJECT_PATH}/build"

for dir in `ls $APP_PATH`; do
  # 如果文件夹不存在，创建文件夹
  path="${BUILD_PATH}/${dir}"
  echo "$path"
  if [ ! -d "$path" ]; then
    mkdir "${path}"
  fi
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o "${path}/app" "${APP_PATH}/${dir}/cmd/main.go"
done