FROM caddy:2.4.6-alpine

# Add curl for health checks
RUN apk add --no-cache curl

# Copy in your production Caddyfile
COPY Caddyfile.production /etc/caddy/Caddyfile
