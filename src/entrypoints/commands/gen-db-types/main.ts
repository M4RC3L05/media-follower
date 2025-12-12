import { toPascalCase } from "@std/text";
import { config, initConfig } from "#src/common/config/mod.ts";
import { CustomDatabase } from "../../../common/database/mod.ts";

initConfig();

using db = new CustomDatabase(config().database.path);

type TypeDef = {
  type: string;
  namedImportPath?: string;
};

const sqliteToTsType: Record<string, TypeDef> = {
  INT: { type: "number" },
  INTEGER: { type: "number" },
  REAL: { type: "number" },
  TEXT: { type: "string" },
  BLOB: { type: "Uint8Array" },
  ANY: { type: "any" },
};

const tableOverrides: Record<string, Record<string, TypeDef>> = {
  release_sources: {
    provider: {
      namedImportPath: "#src/common/database/enums/release-source-provider.ts",
      type: "ReleaseSourceProvider",
    },
    raw: {
      type: "string",
    },
  },
  releases: {
    provider: {
      namedImportPath: "#src/common/database/enums/release-source-provider.ts",
      type: "ReleaseSourceProvider",
    },
    type: {
      namedImportPath: "#src/common/database/enums/release-type.ts",
      type: "ReleaseType",
    },
    raw: {
      type: "string",
    },
  },
};

const mapSqliteTypeToTs = (table: string, column: string, type: string) => {
  const normalized = type.trim().toUpperCase() as keyof typeof sqliteToTsType;

  if (tableOverrides[table]) {
    const override = tableOverrides[table][column];

    if (override) return override;
  }

  if (!sqliteToTsType[normalized]) {
    throw new Error(`SQLite type "${type}" is not mappable`);
  }

  return sqliteToTsType[normalized];
};

const imports: Required<TypeDef>[] = [];
const types: string[] = [];

const rg = /^create table\s*\"?(\S*)\"?\s*\(/gim;
const text = Deno.readTextFileSync("src/common/database/schema.sql");
let match;

// TODO: Replace with matchAll and for of.
while ((match = rg.exec(text)) !== null) {
  const table = match[1];

  const columns =
    (db.prepare(`pragma table_info('${table}')`).all()) as unknown as {
      name: string;
      type: string;
      notnull: boolean;
    }[];

  const tmpType = [`export type ${toPascalCase(`db_${table}_table`)} = {`];

  for (const column of columns) {
    const typedef = mapSqliteTypeToTs(table!, column.name, column.type);

    if (typedef.namedImportPath) {
      imports.push(typedef as Required<TypeDef>);
    }

    tmpType.push(
      `  ${column.name}${column.notnull ? "" : "?"}: ${typedef.type};`,
    );
  }

  tmpType.push("};");

  types.push(tmpType.join("\n"));
}

Deno.writeTextFileSync(
  "src/common/database/types.ts",
  `/* GENERATED FILE CONTENT DO NOT EDIT */\n\n${
    Object.entries(imports.reduce((acc, curr) => {
      if (!acc[curr.namedImportPath]) {
        acc[curr.namedImportPath] = [curr.type];
      } else {
        if (!acc[curr.namedImportPath]!.includes(curr.type)) {
          acc[curr.namedImportPath]!.push(curr.type);
        }
      }

      return acc;
    }, {} as Record<string, string[]>)).map(([imp, types]) =>
      `import { ${types.join(",")} } from "${imp}";`
    ).join("\n")
  }\n\n${types.join("\n\n")}\n`,
);
