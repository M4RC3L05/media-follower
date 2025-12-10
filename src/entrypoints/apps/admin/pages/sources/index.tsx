import type { FunctionComponent } from "preact";
import type { DbReleaseSourcesTable } from "#src/common/database/types.ts";
import { Page } from "../page.tsx";
import { ReleaseSourceProvider } from "#src/common/database/enums/release-source-provider.ts";
import {
  type BluRayComCountry,
  bluRayComCountrySchema,
  type ItunesLookupArtistModelWithExtra,
  itunesLookupArtistModelWithExtraSchema,
} from "#src/common/services/service.ts";

type IndexPageProps = {
  sources: DbReleaseSourcesTable[];
};

const ReleaseSourceProviderFilters = () => (
  <div>
    <h4>Filter by provider:</h4>
    <a href="/sources">All</a> |{" "}
    {Object.values(ReleaseSourceProvider).map((item, i, items) => (
      <>
        <a href={`?provider=${item}`}>{item}</a>
        {i + 1 < items.length ? ` | ` : undefined}
      </>
    ))}
  </div>
);

const SourceDisplayItem: FunctionComponent<
  { title: string; picture?: string; raw: string }
> = ({ raw, title, picture }) => (
  <div>
    <h3>{title}</h3>
    {picture
      ? (
        <img
          style={{
            maxWidth: "100%",
            height: "auto",
            aspectRatio: "1/1",
            maxHeight: "256px",
          }}
          src={picture}
        />
      )
      : undefined}
    <details>
      <summary>Raw:</summary>
      <pre>{raw}</pre>
    </details>
  </div>
);

const ITunesDisplayItem: FunctionComponent<
  {
    source: Omit<DbReleaseSourcesTable, "raw"> & {
      raw: ItunesLookupArtistModelWithExtra;
    };
  }
> = ({ source }) => (
  <SourceDisplayItem
    title={`${source.provider} | ${source.raw.artistName}`}
    picture={source.raw.extra.artistImage}
    raw={JSON.stringify(source.raw, null, 2)}
  />
);

const BluRayComDisplayItem: FunctionComponent<
  { source: Omit<DbReleaseSourcesTable, "raw"> & { raw: BluRayComCountry } }
> = ({ source }) => (
  <SourceDisplayItem
    title={`${source.provider} | ${source.raw.name}`}
    raw={JSON.stringify(source.raw, null, 2)}
  />
);

const SourceDisplayItems: FunctionComponent<{ source: DbReleaseSourcesTable }> =
  ({ source }) => {
    switch (source.provider) {
      case ReleaseSourceProvider.BLU_RAY_COM: {
        return (
          <BluRayComDisplayItem
            source={{
              ...source,
              raw: bluRayComCountrySchema.parse(JSON.parse(source.raw)),
            }}
          />
        );
      }
      case ReleaseSourceProvider.ITUNES: {
        return (
          <ITunesDisplayItem
            source={{
              ...source,
              raw: itunesLookupArtistModelWithExtraSchema.parse(
                JSON.parse(source.raw),
              ),
            }}
          />
        );
      }
    }
  };

const IndexPage: FunctionComponent<IndexPageProps> = ({ sources }) => (
  <Page>
    <Page.Head>
      <title>Media Follower | Sources</title>
    </Page.Head>
    <Page.Body>
      <a href="/">Back to home</a>
      <ReleaseSourceProviderFilters />
      <hr />
      <h1>Sources</h1>
      <a href="/sources/create">Add a new source</a>
      {sources.map((item, i) => <SourceDisplayItems key={i} source={item} />)}
    </Page.Body>
  </Page>
);

export const indexPage = (props: IndexPageProps) => <IndexPage {...props} />;
