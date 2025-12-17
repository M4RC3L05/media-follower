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
  session: z.object({
    secret: z.string().min(1),
  }),
  crypto: z.object({
    pbkdf2: z.object({
      iterations: z.union([
        z.number().min(1),
        z.string().pipe(z.coerce.number()).pipe(z.number().min(1)),
      ]),
      hashFunction: z.enum(["SHA-256", "SHA-384", "SHA-512"]),
      saltLength: z.union([
        z.number().min(1),
        z.string().pipe(z.coerce.number()).pipe(z.number().min(1)),
      ]),
    }),
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
    session: {
      secret: Deno.env.get("SESSION_SECRET") ?? "foobar",
    },
    crypto: {
      pbkdf2: {
        iterations: Deno.env.get("CRYPTO_PBKDF2_ITERATIONS") ?? 800_000,
        hashFunction: Deno.env.get("CRYPTO_PBKDF2_HASH_FUNCTION") ?? "SHA-512",
        saltLength: Deno.env.get("CRYPTO_PBKDF2_SALT_LENGTH") ?? 16,
      },
    },
  });
};

export const config = () => {
  if (!_config) {
    throw new Error("Config was not initiatlized, please call `initConfig`");
  }

  return _config;
};
