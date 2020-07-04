#!/bin/bash -e

# This is useful so we can debug containers running inside of OpenShift that are
# failing to start properly.

echo "cert v0.0.6"

if [ "$PAUSE_ON_START" = "true" ] ; then
  echo
  echo "This container's startup has been paused indefinitely because PAUSE_ON_START has been set."
  echo
  while true; do
    sleep 10
  done
fi

echo "Start cert download server"
/bin/cert
