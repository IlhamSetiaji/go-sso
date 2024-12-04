#!/bin/sh

envsubst < /app/config.template.json > /app/config.json

./main