# syntax=docker/dockerfile:1
# -------------------------------------------------------------
# Production image for Crime Surveillance Frontend (static SPA)
# -------------------------------------------------------------
# 1. Use tiny, production-ready NGINX image
FROM nginx:1.25-alpine as runtime

# 2. Clean default NGINX webroot
RUN rm -rf /usr/share/nginx/html/*

# 3. Copy compiled / static frontend assets
COPY . /usr/share/nginx/html

# Copy entrypoint that injects env vars
COPY docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh

# 4. Expose port 80 (default for NGINX)
EXPOSE 80


ENTRYPOINT ["/docker-entrypoint.sh"]
CMD ["nginx", "-g", "daemon off;"]
