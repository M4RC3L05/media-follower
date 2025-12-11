import { DOMParser } from "@b-fuze/deno-dom";
import {
  type BluRayComCountry,
  bluRayComCountrySchema,
  type BluRayComRelease,
  bluRayComReleaseSchema,
  type BluRayComReleaseWithExtra,
  bluRayComReleaseWithExtraSchema,
  type IBlurayComService,
} from "./service.ts";
import type { IHttpFetch } from "../http/mod.ts";
import type {
  DbReleaseSourcesTable,
  DbReleasesTable,
} from "../database/types.ts";
import { ReleaseSourceProvider } from "../database/enums/release-source-provider.ts";
import { ReleaseType } from "../database/enums/release-type.ts";

type BluRayComServiceProps = {
  httpClient: IHttpFetch;
};

export class BluRayComService implements IBlurayComService {
  #props: BluRayComServiceProps;

  constructor(props: BluRayComServiceProps) {
    this.#props = props;
  }

  #toReleaseWithExtra(
    country: string,
    type: ReleaseType,
    release: BluRayComRelease,
  ): BluRayComReleaseWithExtra {
    return bluRayComReleaseWithExtraSchema.parse({
      ...release,
      extra: {
        country,
        type,
        provider: ReleaseSourceProvider.BLU_RAY_COM,
        link:
          `https://www.blu-ray.com/movies/${release.title_keywords}-Blu-ray/${release.id}/`,
        artworkUrl: release.artworkurl ??
          `https://images.blu-ray.com/movies/covers/${release.id}_medium.jpg`,
      },
    });
  }

  fromReleaseSourceToPersistance(
    item: BluRayComCountry,
  ): DbReleaseSourcesTable {
    return {
      id: item.code,
      provider: ReleaseSourceProvider.BLU_RAY_COM,
      raw: JSON.stringify(item),
    };
  }

  fromReleaseToPersistance(
    item: BluRayComReleaseWithExtra,
  ): DbReleasesTable {
    return {
      releasedAt: item.releasedate.toISOString(),
      id: String(item.id),
      provider: item.extra.provider,
      type: item.extra.type,
      raw: JSON.stringify(item),
    };
  }

  fromPersistanceToRelease(
    row: DbReleasesTable,
  ): BluRayComReleaseWithExtra {
    const parsed = JSON.parse(row.raw);

    return bluRayComReleaseWithExtraSchema.parse(parsed);
  }

  fromPersistanceToReleaseSurce(
    row: DbReleaseSourcesTable,
  ): BluRayComCountry {
    const parsed = JSON.parse(row.raw);

    return bluRayComCountrySchema.parse(parsed);
  }

  async getCountries(): Promise<Array<BluRayComCountry>> {
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
          bluRayComCountrySchema.parse({
            name: Array.from(e.childNodes).at(-1)?.textContent.trim()!,
            code: e.id.replace("country_", "").trim(),
          }),
      );
  }

  #extractMovieListItems(
    html: string,
    releaseType: ReleaseType,
    country: string,
  ) {
    const items: BluRayComReleaseWithExtra[] = [];
    const htmlDoc = new DOMParser().parseFromString(html, "text/html");
    const regexp = /movies\[[0-9]+\]\s*\=\s*(\{.*\})/ig;

    const movieListScript = Array.from(htmlDoc.querySelectorAll("script")).find(
      (e) => e.innerText.includes("function movielist()"),
    )?.textContent;

    if (!movieListScript) {
      return items;
    }

    let match;

    while ((match = regexp.exec(movieListScript)) !== null) {
      items.push(this.#toReleaseWithExtra(
        country,
        releaseType,
        bluRayComReleaseSchema.parse(eval(
          `(${match[1]!.trim()})`,
        )),
      ));
    }

    return items;
  }

  async getBlurayReleasesByCountryForMonth(
    country: string,
    year: number,
    month: number,
  ): Promise<BluRayComReleaseWithExtra[]> {
    const [blurayTextContent, dvdTextContent] = await Promise.all([
      await this.#props.httpClient.fetchText(
        `https://www.blu-ray.com/movies/releasedates.php?year=${year}&month=${month}`,
        {
          headers: {
            "User-Agent":
              "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36",
            "Cookie": `country=${country}`,
          },
        },
      ),
      await this.#props.httpClient.fetchText(
        `https://www.blu-ray.com/dvd/releasedates.php?year=${year}&month=${month}`,
        {
          headers: {
            "User-Agent":
              "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36",
            "Cookie": `country=${country}`,
          },
        },
      ),
    ]);

    return [
      ...this.#extractMovieListItems(
        blurayTextContent,
        ReleaseType.BLURAY,
        country,
      ),
      ...this.#extractMovieListItems(
        dvdTextContent,
        ReleaseType.DVD,
        country,
      ),
    ];
  }
}
