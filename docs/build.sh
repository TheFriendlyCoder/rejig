#!/usr/bin/env bash
set -e

if [ "${1:-serve}" != "serve" ]; then
    echo "Usage: build.sh [serve]"
    exit 1
fi

which virtualenv 1>/dev/null || (echo "virtualenv not installed"; exit 1)
if [ ! -d "venv" ]; then
    virtualenv -p python3 venv
    source venv/bin/activate
    pip install --upgrade pip
    pip install -r requirements.txt
fi

if [ ! -e "VIRTUAL_ENV" ]; then
    source venv/bin/activate
fi

which mkdocs 1>/dev/null || (echo "mkdocs executable is not found"; exit 1)

mkdocs build --strict

if [ "$1" = "serve" ]; then
    mkdocs serve
else
    open public/index.html
fi
