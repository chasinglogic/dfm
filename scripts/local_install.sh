#!/usr/bin/env bash

INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"

REPO=$(git rev-parse --show-toplevel)

cd $REPO
go build .

mv ./dfm "$INSTALL_DIR/"
