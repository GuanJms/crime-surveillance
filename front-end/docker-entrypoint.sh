#!/bin/sh
# Generate env.js from container environment variables so the browser can read them
set -eu

# List of variable names we want to expose
VARS="BROKER_SERVICE_URL"

OUT=/usr/share/nginx/html/env.js
{
  echo "window.__ENV__ = {";
  for NAME in $VARS; do
    VALUE=$(printenv "$NAME" || true)
    if [ -n "$VALUE" ]; then
      # Escape double quotes for JS string
      ESC=$(printf '%s' "$VALUE" | sed 's/"/\\"/g')
      echo "  $NAME: \"$ESC\",";
    fi
  done
  echo "};";
} > "$OUT"

exec "$@"
