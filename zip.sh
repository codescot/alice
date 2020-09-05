#!/bin/sh
GOOS=linux go build -o gameinfo
zip function.zip gameinfo
