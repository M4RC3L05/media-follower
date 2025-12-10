FROM docker.io/denoland/deno:alpine-2.5.6

RUN mkdir /app
RUN chown -R deno:deno /app

USER deno

WORKDIR /app

COPY --chown=deno:deno . .

RUN deno install --frozen --unstable-npm-lazy-caching --entrypoint \
  src/entrypoints/apps/admin/main.ts \
  src/entrypoints/apps/rss-feed/main.ts \
  src/entrypoints/jobs/sync-releases/main.ts \
  src/entrypoints/jobs/sync-release-sources/main.ts

RUN BUILD_DRY_RUN=true DATABASE_PATH=":memory:" timeout 30s deno run -A --cached-only --frozen --unstable-npm-lazy-caching src/entrypoints/apps/admin/main.ts || true
RUN BUILD_DRY_RUN=true DATABASE_PATH=":memory:" timeout 30s deno run -A --cached-only --frozen --unstable-npm-lazy-caching src/entrypoints/apps/rss-feed/main.ts || true
# RUN BUILD_DRY_RUN=true DATABASE_PATH=":memory:" timeout 30s deno run -A --cached-only --frozen --unstable-npm-lazy-caching src/entrypoints/jobs/sync-releases/main.ts || true
# RUN BUILD_DRY_RUN=true DATABASE_PATH=":memory:" timeout 30s deno run -A --cached-only --frozen --unstable-npm-lazy-caching src/entrypoints/jobs/sync-release-sources/main.ts || true

RUN mkdir /app/data

VOLUME [ "/app/data" ]

EXPOSE 4321
