#!/bin/sh
GOPATH=/Users/zzm/go
GOROOT=/Users/zzm/Documents/tools/go
APP_PATH=/Users/zzm/go/awesomeProject
PROJECT=tag_engine

start() {
  echo "start tag engine.........."
  /usr/bin/nohup ${APP_PATH}/bin/${PROJECT} -c ${APP_PATH}/etc/app.toml >/dev/null 2>&1 &
  return
}

stop() {
  echo "stop tag engine........."
  ${APP_PATH}/bin/${PROJECT} -s shutdown
}

if [[ $1 = "start" ]]; then
  start
elif [[ $1 = "stop" ]]; then
  stop
else
  echo "usage: start [start|stop]"
fi

