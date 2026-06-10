FROM docker.io/denoland/deno:alpine-2.8.2

RUN mkdir /app
RUN chown -R deno:deno /app

USER deno

WORKDIR /app

COPY --chown=deno:deno . .

RUN deno install --skip-types --prod --entrypoint \
  src/entrypoints/apps/admin/main.ts \
  src/entrypoints/apps/rss-feed/main.ts \
  src/entrypoints/jobs/sync-inputs/main.ts \
  src/entrypoints/jobs/fetch-outputs/main.ts

RUN BUILD_DRY_RUN=true DATABASE_PATH=":memory:" timeout 20s deno run -A --cached-only src/entrypoints/apps/admin/main.ts
RUN BUILD_DRY_RUN=true DATABASE_PATH=":memory:" timeout 20s deno run -A --cached-only src/entrypoints/apps/rss-feed/main.ts
RUN BUILD_DRY_RUN=true DATABASE_PATH=":memory:" timeout 20s deno run -A --cached-only src/entrypoints/jobs/sync-inputs/main.ts
RUN BUILD_DRY_RUN=true DATABASE_PATH=":memory:" timeout 20s deno run -A --cached-only src/entrypoints/jobs/fetch-outputs/main.ts

RUN mkdir /app/data

VOLUME [ "/app/data" ]

EXPOSE 4321
