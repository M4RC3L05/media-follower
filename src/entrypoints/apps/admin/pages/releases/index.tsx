import type { FunctionComponent } from "preact";
import { Page } from "../page.tsx";
import type { DbReleasesTable } from "../../../../../common/database/types.ts";
import { ReleaseSourceProvider } from "../../../../../common/database/enums/release-source-provider.ts";
import { ReleaseType } from "../../../../../common/database/enums/release-type.ts";
import {
  type BluRayComReleaseWithExtra,
  bluRayComReleaseWithExtraSchema,
  type ItunesLookupAlbumModelWithExtra,
  itunesLookupAlbumModelWithExtraSchema,
  type ItunesLookupSongModelWithExtra,
  itunesLookupSongModelWithExtraSchema,
} from "../../../../../common/services/service.ts";

type IndexPageProps = {
  releases: DbReleasesTable[];
  url: URL;
  paginatio: { page: number; limit: number };
};

const ReleaseDisplayItem: FunctionComponent<
  { image: string; title: string; releasedAt: Date; link: string; raw: string }
> = ({ image, link, raw, releasedAt, title }) => (
  <article>
    <h3>{title}</h3>
    <img style={{ aspectRatio: "16/9" }} src={image} />
    <p>
      {releasedAt.toDateString()}
      {releasedAt > new Date()
        ? (
          <>
            {" "}
            <em>(To be released)</em>
          </>
        )
        : undefined}
    </p>
    <a target="_blank" href={link}>View on source</a>
    <details>
      <summary>Raw:</summary>
      <pre>{raw}</pre>
    </details>
  </article>
);

const BluRayComReleaseDisplayItem: FunctionComponent<
  { release: Omit<DbReleasesTable, "raw"> & { raw: BluRayComReleaseWithExtra } }
> = ({ release }) => (
  <ReleaseDisplayItem
    image={release.raw.extra.artworkUrl}
    link={release.raw.extra.link}
    raw={JSON.stringify(release.raw, null, 2)}
    releasedAt={release.raw.releasedate}
    title={`[${release.provider} | ${release.type} | ${release.raw.extra.country.toUpperCase()}] ${release.raw.title}`}
  />
);

const ItunesReleaseDisplayItem: FunctionComponent<
  {
    release: Omit<DbReleasesTable, "raw"> & {
      raw: ItunesLookupAlbumModelWithExtra | ItunesLookupSongModelWithExtra;
    };
  }
> = ({ release }) => {
  const image = release.raw.artworkUrl100
    .split("/")
    .map((segment, index, array) =>
      index === array.length - 1 ? "512x512bb.jpg" : segment
    ).join("/");
  const title = release.type === ReleaseType.SONG
    ? (release.raw as ItunesLookupSongModelWithExtra).trackName
    : release.raw.collectionName;
  const link = release.type === ReleaseType.SONG
    ? (release.raw as ItunesLookupSongModelWithExtra).trackViewUrl
    : release.raw.collectionViewUrl;

  return (
    <ReleaseDisplayItem
      image={image}
      link={link}
      raw={JSON.stringify(release.raw, null, 2)}
      releasedAt={release.raw.releaseDate}
      title={`[${release.provider} | ${release.type}] ${title}`}
    />
  );
};

const ReleaseDisplayItems: FunctionComponent<{ release: DbReleasesTable }> = (
  { release },
) => {
  switch (release.provider) {
    case ReleaseSourceProvider.BLU_RAY_COM: {
      return (
        <BluRayComReleaseDisplayItem
          release={{
            ...release,
            raw: bluRayComReleaseWithExtraSchema.parse(
              JSON.parse(release.raw),
            ),
          }}
        />
      );
    }
    case ReleaseSourceProvider.ITUNES: {
      return (
        <ItunesReleaseDisplayItem
          release={{
            ...release,
            raw: release.type === ReleaseType.ALBUM
              ? itunesLookupAlbumModelWithExtraSchema.parse(
                JSON.parse(release.raw),
              )
              : itunesLookupSongModelWithExtraSchema.parse(
                JSON.parse(release.raw),
              ),
          }}
        />
      );
    }
  }
};

const IndexPage: FunctionComponent<IndexPageProps> = (
  { releases, url, paginatio: { page } },
) => {
  const prevPageLink = new URL(url);
  prevPageLink.searchParams.set("page", String(Math.max(page - 1, 0)));

  const nextPageLink = new URL(url);
  nextPageLink.searchParams.set("page", String(page + 1));

  return (
    <Page>
      <Page.Head>
        <title>Media Follower | Releases</title>
      </Page.Head>
      <Page.Body>
        <header>
          <h1>Releases</h1>
          <nav>
            <a href="/">Back to home</a>
          </nav>
        </header>

        <main>
          <section
            style={{
              paddingTop: "2rem",
              textAlign: "left",
              position: "sticky",
              top: 0,
              background: "var(--bg)",
              zIndex: 2,
            }}
          >
            <h6 style={{ marginTop: 0, marginBottom: 0 }}>
              Filter by provider: <a href="/releases" class="button">All</a> |
              {" "}
              {Object.values(ReleaseSourceProvider).map((value, i, items) => (
                <>
                  <a href={`?provider=${value}`} class="button">{value}</a>
                  {i < items.length - 1 ? " | " : ""}
                </>
              ))}
            </h6>
            <br />
            <h6 style={{ marginTop: 0, marginBottom: 0 }}>
              Filter by type: <a href="/releases" class="button">All</a> |{" "}
              {Object.values(ReleaseType).map((value, i, items) => (
                <>
                  <a href={`?type=${value}`} class="button">{value}</a>
                  {i < items.length - 1 ? " | " : ""}
                </>
              ))}
            </h6>
            <br />
            <h6 style={{ marginTop: 0, marginBottom: 0 }}>
              Pagination:{" "}
              <a href={prevPageLink.toString()} class="button">Prev</a> |{" "}
              <a href={nextPageLink.toString()} class="button">Next</a>
            </h6>
          </section>

          <section>
            {releases.map((item, i) => (
              <ReleaseDisplayItems
                key={i}
                release={item}
              />
            ))}
          </section>
        </main>
      </Page.Body>
    </Page>
  );
};

export const indexPage = (props: IndexPageProps) => <IndexPage {...props} />;
