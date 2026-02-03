import type { VNode } from "preact";
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
  type Input,
  type Output,
  steamGamesFreePromosInputSchema,
  steamGamesFreePromosOutputSchema,
  SteamGamesFreePromoTypes,
} from "./types.ts";
import { Feed } from "feed";
import { DOMParser } from "@b-fuze/deno-dom";
import { distinctBy } from "@std/collections";
import type { IDatabase } from "#src/common/database/database.ts";
import { inputListItem, outputListItem } from "./components/mod.tsx";

type SteamGamesFreePromosProps = {
  httpClient: IHttpFetch;
  database: IDatabase;
};

export class SteamGamesFreePromosProvider implements
  IProvider<
    EInputProvider.STEAM_GAMES_FREE_PROMOS,
    Input,
    Output
  > {
  readonly provider = EInputProvider.STEAM_GAMES_FREE_PROMOS;
  #url = "https://steamdb.info/upcoming/free/";
  #props: SteamGamesFreePromosProps;

  constructor(props: SteamGamesFreePromosProps) {
    this.#props = props;
  }

  // deno-lint-ignore require-await
  async lookupInput(_term: string): Promise<Input | undefined> {
    return { url: this.#url };
  }

  // deno-lint-ignore require-await
  async fetchInput(_row: DbInputsTable): Promise<Input> {
    return { url: this.#url };
  }

  async fetchOutputs(input: Input): Promise<Output[]> {
    const txtHtml = await this.#props.httpClient.fetchText(input.url, {
      headers: {
        "User-Agent":
          "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36",
      },
    });

    const dom = new DOMParser().parseFromString(txtHtml, "text/html");

    return distinctBy(
      Array.from(
        dom.querySelectorAll("#main .container .row .panel-sale"),
      ).filter((e) => e.getAttribute("data-appid")?.trim() !== "730").map(
        (ele) => {
          const catElement = ele.querySelector("div.cat");
          const promoType = catElement?.className.includes("cat-free-to-keep")
            ? SteamGamesFreePromoTypes.FREE_TO_KEEP
            : catElement?.className.includes("cat-play-for-free")
            ? SteamGamesFreePromoTypes.FREE_TO_PLAY
            : undefined;

          const [startDate, endDate] = Array.from(
            ele.querySelectorAll("div.panel-sale-time"),
          ).map(
            (ele) =>
              ele.querySelector("relative-time")?.getAttribute("datetime"),
          ).filter((item) => !!item).toSorted();

          return steamGamesFreePromosOutputSchema.parse({
            id: ele.getAttribute("data-appid")?.trim(),
            image: ele.querySelector("img.sale-image")?.getAttribute("src")
              ?.trim(),
            link: Array.from(
              ele.querySelector("div.app-history-type")?.children ?? [],
            ).map((x) => x.getAttribute("href")?.trim()).find((link) =>
              link?.startsWith("https://store.steampowered.com/app")
            ),
            name: ele.querySelector("h4.panel-sale-name a")?.textContent.trim(),
            promoType,
            startDate,
            endDate,
          });
        },
      ),
      (item) => item.id,
    );
  }

  // deno-lint-ignore require-await
  async queryOutputs(
    { pagination }: IProviderRepositoryQueryOutputsProps,
  ): Promise<DbOutputsTable[]> {
    return this.#props.database.sql.all`
      select id, input_id, provider, json(raw) as raw
      from outputs
      where provider = ${EInputProvider.STEAM_GAMES_FREE_PROMOS}
      order by raw->>'startDate' desc
      limit ${pagination.limit}
      offset ${pagination.page * pagination.limit}
    ` as DbOutputsTable[];
  }

  getOutputsFeed(_props: IProviderFeedGetOutputsFeedProps): Feed {
    const rows = this.#props.database.sql.all`
      select id, input_id, provider, json(raw) as raw
      from outputs
      where provider = ${EInputProvider.STEAM_GAMES_FREE_PROMOS}
      order by raw->>'startDate' desc
      limit 200
    ` as DbOutputsTable[];

    const prefix = [EInputProvider.STEAM_GAMES_FREE_PROMOS].filter(Boolean);

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
        date: output.startDate,
        link: output.link,
        title: output.name,
        id:
          `${EInputProvider.STEAM_GAMES_FREE_PROMOS}@${output.promoType}@${output.id}`,
        image: output.image,
        description:
          `<p>${output.name} (${output.promoType})</p><p>From: ${output.startDate.toDateString()} to ${output.endDate.toDateString()}</p>`,
      });
    }

    return feed;
  }

  renderInputListItem(row: DbInputsTable): VNode {
    return inputListItem({ input: this.fromPersistenceToInput(row) });
  }

  renderOutputListItem(row: DbOutputsTable): VNode {
    return outputListItem({
      output: this.fromPersistenceToOutput(row),
    });
  }

  fromOutputToJsonPatchPersistance(
    row: DbInputsTable,
    item: Output,
  ): DbOutputsTable {
    return this.fromOutputToPersistence(row, item);
  }

  fromInputToPersistence(item: Input): DbInputsTable {
    return {
      id: item.url,
      provider: EInputProvider.STEAM_GAMES_FREE_PROMOS,
      raw: JSON.stringify({ url: item.url }),
    };
  }

  fromPersistenceToInput(row: DbInputsTable): Input {
    return steamGamesFreePromosInputSchema.parse(JSON.parse(row.raw));
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
    return steamGamesFreePromosOutputSchema.parse(JSON.parse(row.raw));
  }
}
