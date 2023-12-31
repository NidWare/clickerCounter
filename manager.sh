#!/bin/bash

CONTAINER_NAME="golang_server"
VOLUME_PATH="$(pwd)"

case $1 in
  start)
    # Build and start the container detached with restart policy
    docker build -t golang_server .
    docker run -d --restart=unless-stopped --name $CONTAINER_NAME -p 8080:8080 -v $VOLUME_PATH/clicks.db:/app/clicks.db golang_server
    ;;
  stop)
    # Stop and remove the container
    docker stop $CONTAINER_NAME
    docker rm $CONTAINER_NAME
    ;;
  restart)
    # Restart the container
    docker restart $CONTAINER_NAME
    ;;
  status)
    # Check the status of the container
    docker ps -a | grep $CONTAINER_NAME
    ;;
  *)
    echo "Usage: ./manager.sh {start|stop|restart|status}"
    exit 1
    ;;
esac
