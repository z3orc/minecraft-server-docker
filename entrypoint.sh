#!/bin/sh
/app/runner -version=$VERSION -dir=/data -jar=$SERVER_JAR -memory=$MEMORY -timeout=$TIMEOUT -sigkill=$USE_SIGKILL
