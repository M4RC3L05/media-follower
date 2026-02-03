import type { IDatabase } from "#src/common/database/database.ts";

export const runMigrations = async (db: IDatabase) => {
  const migrationsPathRelative = "../database/migrations";
  const migartionsDir = new URL(migrationsPathRelative, import.meta.url);

  const migrationFiles = (await Array.fromAsync(Deno.readDir(migartionsDir)))
    .filter((item) => item.isFile && item.name.endsWith(".sql")).map((item) =>
      item.name
    ).sort();

  for (const file of migrationFiles) {
    const filePath = new URL(
      `${migrationsPathRelative}/${file}`,
      import.meta.url,
    );

    db.sql.run`${await Deno.readTextFile(filePath)}`;
  }
};
