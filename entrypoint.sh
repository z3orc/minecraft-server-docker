#!/bin/sh
/app/runner -dir=/data -jar=$SERVER_JAR -memory=$MEMORY -timeout=$TIMEOUT -sigkill=$USE_SIGKILL
