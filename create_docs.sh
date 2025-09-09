#!/bin/bash

# build swagger docs.
# Don't forget to install swaggo/swag:
# go install github.com/swaggo/swag/cmd/swag@latest

swag init --dir ./cmd,./internal/entity,./internal/controller --output ./internal/docs
