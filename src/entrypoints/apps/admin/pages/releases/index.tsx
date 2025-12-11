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
  <div>
    <h3>{title}</h3>
    <img
      style={{
        maxWidth: "100%",
        height: "auto",
        aspectRatio: "1/1",
        maxHeight: "256px",
      }}
      src={image}
    />
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
  </div>
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
        <a href="/">Back to home</a>
        <h4>Filter by provider</h4>
        <a href="/releases">All</a> |{" "}
        {Object.values(ReleaseSourceProvider).map((value, i, items) => (
          <>
            <a href={`?provider=${value}`}>{value}</a>
            {i < items.length - 1 ? " | " : ""}
          </>
        ))}
        <hr />
        <h4>Filter by type</h4>
        <a href="/releases">All</a> |{" "}
        {Object.values(ReleaseType).map((value, i, items) => (
          <>
            <a href={`?type=${value}`}>{value}</a>
            {i < items.length - 1 ? " | " : ""}
          </>
        ))}
        <hr />
        <h1>Releases</h1>
        <a href={prevPageLink.toString()}>Prev</a> |{" "}
        <a href={nextPageLink.toString()}>Next</a>
        <hr />
        {releases.map((item, i) => (
          <ReleaseDisplayItems
            key={i}
            release={item}
          />
        ))}
      </Page.Body>
    </Page>
  );
};

export const indexPage = (props: IndexPageProps) => <IndexPage {...props} />;
