#!/usr/bin/env bash

export $(grep -v '^#' ../.env | xargs -d '\n')
export SERVER_HOST=$SERVER1_HOST

go1.14.2 run .
