#!/usr/bin/env bash
docker run --rm -p 5672:5672 -p 15672:15672 -p 61613:61613 -p 15670:15670 -p 15674:15674 -e RABBITMQ_DEFAULT_VHOST=law -v ${PWD}/enabled_plugins:/etc/rabbitmq/enabled_plugins rabbitmq:3.5.7-management
