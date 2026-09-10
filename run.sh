#!/usr/bin/env sh
set -eu
docker compose up -d --build --wait
printf 'Talknet is ready at http://localhost:%s\n' "${TALKNET_PORT:-8088}"
