# syntax=docker/dockerfile:1

FROM node:22-alpine AS deps

WORKDIR /app

RUN corepack enable

COPY package.json pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile

FROM node:22-alpine AS build

WORKDIR /app

RUN corepack enable

COPY --from=deps /app/node_modules ./node_modules
COPY package.json pnpm-lock.yaml ./
COPY . .
RUN pnpm build

FROM node:22-alpine AS ssr

WORKDIR /app

ENV NODE_ENV=production

COPY --from=build /app/package.json ./package.json
COPY --from=build /app/node_modules ./node_modules
COPY --from=build /app/dist ./dist

EXPOSE 3000

CMD ["node", "/app/node_modules/.pnpm/node_modules/srvx/bin/srvx.mjs", "serve", "--prod", "--entry", "/app/dist/server/server.js", "--static", "/app/dist/client", "--port", "3000", "--host", "0.0.0.0"]

FROM node:22-alpine AS runtime

WORKDIR /app

ENV NODE_ENV=production

RUN apk add --no-cache nginx supervisor

COPY --from=build /app/package.json ./package.json
COPY --from=build /app/node_modules ./node_modules
COPY --from=build /app/dist ./dist

RUN cat > /etc/nginx/http.d/default.conf <<'EOF'
server {
    listen 3000;
    server_name _;

    access_log /dev/stdout;
    error_log /dev/stderr warn;
    resolver 127.0.0.11 valid=30s ipv6=off;

    location /api/ {
        set $backend_origin http://backend:8080;
        proxy_pass $backend_origin;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    location /assets/ {
        alias /app/dist/client/assets/;
        expires 1y;
        add_header Cache-Control "public, immutable";
        try_files $uri =404;
    }

    location / {
        proxy_pass http://127.0.0.1:3001;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
EOF

RUN cat > /etc/supervisord.conf <<'EOF'
[supervisord]
nodaemon=true
logfile=/dev/null
pidfile=/tmp/supervisord.pid

[program:frontend-ssr]
command=node /app/node_modules/.pnpm/node_modules/srvx/bin/srvx.mjs serve --prod --entry /app/dist/server/server.js --static /app/dist/client --port 3001 --host 127.0.0.1
directory=/app
autorestart=true
stdout_logfile=/dev/fd/1
stdout_logfile_maxbytes=0
stderr_logfile=/dev/fd/2
stderr_logfile_maxbytes=0

[program:nginx]
command=nginx -g "daemon off;"
autorestart=true
stdout_logfile=/dev/fd/1
stdout_logfile_maxbytes=0
stderr_logfile=/dev/fd/2
stderr_logfile_maxbytes=0
EOF

EXPOSE 3000

CMD ["supervisord", "-c", "/etc/supervisord.conf"]
