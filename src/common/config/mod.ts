import z from "@zod/zod";

const configSchema = z.object({
  apps: z.object({
    rssFeed: z.object({
      host: z.hostname().trim(),
      port: z.union([z.number(), z.string().trim()]).pipe(z.coerce.number()),
    }),
    admin: z.object({
      host: z.hostname().trim(),
      port: z.union([z.number(), z.string().trim()]).pipe(z.coerce.number()),
    }),
  }),
  database: z.object({
    path: z.string().min(1),
  }),
});

export type ProjectConfig = z.infer<typeof configSchema>;

let _config: ProjectConfig;

export const initConfig = () => {
  _config = configSchema.parse({
    apps: {
      rssFeed: {
        host: Deno.env.get("APPS_RSS_FEED_HOST") ?? "127.0.0.1",
        port: Deno.env.get("APPS_RSS_FEED_PORT") ?? 4321,
      },
      admin: {
        host: Deno.env.get("APPS_ADMIN_HOST") ?? "127.0.0.1",
        port: Deno.env.get("APPS_ADMIN_PORT") ?? 4322,
      },
    },
    database: {
      path: Deno.env.get("DATABASE_PATH") ?? "./data/app.db",
    },
  });
};

export const config = () => {
  if (!_config) {
    throw new Error("Config was not initiatlized, please call `initConfig`");
  }

  return _config;
};
