#!/bin/bash

commands=(compile run dev build)

invalidUsage () {
    echo "Valid Usage: ./cmd <command>"
    echo "commands:"
    for command in ${commands[@]}; do
      printf "\t$command\n"
    done
    exit 1
}

if [ -z $1 ]; then
  invalidUsage
fi

if [ $1 == compile ]; then
  templ generate
elif [ $1 == run ]; then
  go run .
elif [ $1 == dev ]; then
  templ generate --watch --cmd "air"
elif [ $1 == build ]; then
  go build
else
  invalidUsage
fi
