#!/bin/sh
/app/runner -dir=/data -jar=$SERVER_JAR -timeout=$TIMEOUT -sigkill=$USE_SIGKILL
