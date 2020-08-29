#!/bin/bash

if [[ $1 = "build" ]];
then
    shift # remove the first argument
    command "$(dirname "$0")/omm" $PWD "$@" "--compile"
else #otherwise, just run it

    #if it looks like `oat run` remove the `run`
    if [[ $1 = "run" ]];
    then
        shift
    fi

    command "$(dirname "$0")/omm" $PWD "$@" "--run"
fi