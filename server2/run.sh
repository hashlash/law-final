#!/usr/bin/env bash

export $(grep -v '^#' ../.env | xargs -d '\n')

go1.14.2 run .
