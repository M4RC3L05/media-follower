import { DOMParser } from "@b-fuze/deno-dom";
import { EInputProvider } from "#src/common/database/enums/mod.ts";
import type {
  DbInputsTable,
  DbOutputsTable,
} from "#src/common/database/types.ts";
import type { IHttpFetch } from "#src/common/http/mod.ts";
import type {
  IProvider,
  IProviderFeedGetOutputsFeedProps,
  IProviderRepositoryQueryOutputsProps,
} from "../interfaces.ts";
import {
  bluRayComPhysicalReleaseInputSchema,
  bluRayComPhysicalReleaseOutputSchema,
  type BluRayComPhysicalReleaseOutputWithExtra,
  bluRayComPhysicalReleaseOutputWithExtraSchema,
  EBLuRayComPhysicalReleaseType,
  type Input,
  type Output,
} from "./types.ts";
import type { IDatabase } from "#src/common/database/database.ts";
import { inputListItem, outputListItem } from "./components/mod.tsx";
import type { VNode } from "preact";
import z from "@zod/zod";
import { Feed } from "feed";

type BluRayComPhysicalReleasesProviderProps = {
  httpClient: IHttpFetch;
  database: IDatabase;
};

export class BluRayComPhysicalReleasesProvider
  implements
    IProvider<EInputProvider.BLU_RAY_COM_PHYSICAL_RELEASE, Input, Output> {
  readonly provider = EInputProvider.BLU_RAY_COM_PHYSICAL_RELEASE;

  #props: BluRayComPhysicalReleasesProviderProps;

  constructor(props: BluRayComPhysicalReleasesProviderProps) {
    this.#props = props;
  }

  async lookupInput(
    term: string,
  ): Promise<Input | undefined> {
    const html = await this.#props.httpClient.fetchText(
      `https://www.blu-ray.com`,
      {
        headers: {
          "User-Agent":
            "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36",
          "Cookie": `country=all`,
        },
      },
    );
    const htmlDoc = new DOMParser().parseFromString(html, "text/html");

    return Array.from(
      htmlDoc.querySelectorAll("#search_locale #search_locale_list>ul>li"),
    )
      .map(
        (e) =>
          bluRayComPhysicalReleaseInputSchema.parse({
            name: Array.from(e.childNodes).at(-1)?.textContent?.trim(),
            code: e.id.replace("country_", "").trim(),
          }),
      ).find((item) => item.code === term);
  }

  async fetchInput(dbInput: DbInputsTable): Promise<Input> {
    const parsed = bluRayComPhysicalReleaseInputSchema.parse(
      JSON.parse(dbInput.raw),
    );

    const input = await this.lookupInput(parsed.code);

    if (!input) throw new Error("Unable to fetch by input", { cause: dbInput });

    return input;
  }

  #extractMovieListItems(
    html: string,
    releaseType: EBLuRayComPhysicalReleaseType,
  ) {
    const items: BluRayComPhysicalReleaseOutputWithExtra[] = [];
    const htmlDoc = new DOMParser().parseFromString(html, "text/html");
    const regexp = /movies\[[0-9]+\]\s*\=\s*(\{.*\})/ig;

    const movieListScript = Array.from(htmlDoc.querySelectorAll("script")).find(
      (e) => e.innerText.includes("function movielist()"),
    )?.textContent;

    if (!movieListScript) {
      return items;
    }

    let match;

    // TODO: Replace with matchAll and for of.
    while ((match = regexp.exec(movieListScript)) !== null) {
      const release = bluRayComPhysicalReleaseOutputSchema.parse(
        eval(`(${match[1]!.trim()})`),
      );
      items.push(bluRayComPhysicalReleaseOutputWithExtraSchema.parse({
        ...release,
        extra: {
          type: releaseType,
          link:
            `https://www.blu-ray.com/movies/${release.title_keywords}-Blu-ray/${release.id}/`,
          artworkUrl: release.artworkurl ??
            `https://images.blu-ray.com/movies/covers/${release.id}_medium.jpg`,
        },
      }));
    }

    return items;
  }

  async fetchOutputs(input: Input): Promise<Output[]> {
    const now = new Date();
    const year = now.getFullYear();
    const month = now.getMonth() + 1;

    const [blurays, dvds] = await Promise.all([
      this.#props.httpClient.fetchText(
        `https://www.blu-ray.com/movies/releasedates.php?year=${year}&month=${month}`,
        {
          headers: {
            "User-Agent":
              "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36",
            "Cookie": `country=${input.code}`,
          },
        },
      ).then((text) =>
        this.#extractMovieListItems(text, EBLuRayComPhysicalReleaseType.BLURAY)
      ),
      this.#props.httpClient.fetchText(
        `https://www.blu-ray.com/dvd/releasedates.php?year=${year}&month=${month}`,
        {
          headers: {
            "User-Agent":
              "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36",
            "Cookie": `country=${input.code}`,
          },
        },
      ).then((text) =>
        this.#extractMovieListItems(text, EBLuRayComPhysicalReleaseType.DVD)
      ),
    ]);

    return [...blurays, ...dvds];
  }

  renderInputListItem(row: DbInputsTable): VNode {
    return inputListItem({
      input: bluRayComPhysicalReleaseInputSchema.parse(JSON.parse(row.raw)),
    });
  }

  renderOutputListItem(row: DbOutputsTable): VNode {
    return outputListItem({
      output: this.fromPersistenceToOutput(row),
      outputRow: row,
    });
  }

  // deno-lint-ignore require-await
  async queryOutputs(
    { pagination }: IProviderRepositoryQueryOutputsProps,
  ): Promise<DbOutputsTable[]> {
    return this.#props.database.sql<DbOutputsTable>`
      select id, input_id, provider, json(outputs.raw) as raw
      from outputs
      where provider = ${EInputProvider.BLU_RAY_COM_PHYSICAL_RELEASE}
      order by outputs.raw->>'releasedate' desc, "rowid" desc
      limit ${pagination.limit}
      offset ${pagination.page * pagination.limit}
    `;
  }

  getOutputsFeed({ queries }: IProviderFeedGetOutputsFeedProps): Feed {
    const { type, country } = z.object({
      type: z.enum(EBLuRayComPhysicalReleaseType).optional(),
      country: z.string().min(1).optional(),
    })
      .parse(queries ?? {});

    const rows = this.#props.database.sql<DbOutputsTable>`
        select id, input_id, provider, json(outputs.raw) as raw
        from outputs
        where provider = ${EInputProvider.BLU_RAY_COM_PHYSICAL_RELEASE}
        and (${type ?? null} is null or outputs.raw->'extra'->>'type' = ${
      type ?? null
    })
        and (${country ?? null} is null or input_id = ${country ?? null})
        and outputs.raw->>'releasedate' <= strftime('%Y-%m-%dT%H:%M:%fZ' , 'now')
        order by outputs.raw->>'releasedate' desc, "rowid" desc
        limit 200
      `;

    const prefix = [EInputProvider.BLU_RAY_COM_PHYSICAL_RELEASE, type, country]
      .filter(
        Boolean,
      );

    const feed = new Feed({
      title: `${prefix ? `[${prefix.join(" | ")}] ` : ""}Media follower`,
      description: `Get the latest${
        prefix ? ` ${prefix.join(" and ")} ` : " "
      }media releases`,
      id: `media_follower${prefix ? `_${prefix.join("_")}` : ""}`,
      copyright: "Media Follower",
      updated: new Date(),
    });

    const outputs = rows.map((row) => ({
      ...this.fromPersistenceToOutput(row),
      row,
    }));

    for (const output of outputs) {
      feed.addItem({
        date: output.releasedate,
        link: output.extra.link,
        title: output.title,
        id:
          `${EInputProvider.BLU_RAY_COM_PHYSICAL_RELEASE}@${output.extra.type}@${output.id}`,
        image: output.extra.artworkUrl,
      });
    }

    return feed;
  }

  fromInputToPersistence(item: Input): DbInputsTable {
    return {
      id: item.code,
      provider: EInputProvider.BLU_RAY_COM_PHYSICAL_RELEASE,
      raw: JSON.stringify(item),
    };
  }

  fromPersistenceToInput(row: DbInputsTable): Input {
    const parsed = JSON.parse(row.raw);

    return bluRayComPhysicalReleaseInputSchema.parse(parsed);
  }

  fromOutputToPersistence(row: DbInputsTable, item: Output): DbOutputsTable {
    return {
      id: String(item.id),
      input_id: row.id,
      provider: row.provider,
      raw: JSON.stringify(item),
    };
  }

  fromPersistenceToOutput(row: DbOutputsTable): Output {
    const parsed = JSON.parse(row.raw);

    return bluRayComPhysicalReleaseOutputWithExtraSchema.parse(parsed);
  }
}
