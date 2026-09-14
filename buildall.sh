#!/bin/bash

# set up project root
PROJECT_ROOT="$(cd "$(dirname "$0")" && pwd)"
cd "$PROJECT_ROOT"

# create bin dir
mkdir -p bin

# build all apps
go build -o bin/server-game ./apps/server-game
go build -o bin/client-game ./apps/client-game
go build -o bin/entity-creator ./apps/entity-creator
