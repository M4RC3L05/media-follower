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

const SourceDisplayItem: FunctionComponent<
  { title: string; picture?: string; raw: string }
> = ({ raw, title, picture }) => (
  <article>
    <h3>{title}</h3>
    {picture
      ? <img style={{ aspectRatio: "16/9" }} src={picture} />
      : undefined}
    <details>
      <summary>Raw:</summary>
      <pre>{raw}</pre>
    </details>
  </article>
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
      <header>
        <h1>Sources</h1>
        <nav>
          <a href="/">Back to home</a>
          <a href="/sources/create">Add a new source</a>
        </nav>
      </header>

      <main>
        <section
          style={{
            paddingTop: "2rem",
            textAlign: "center",
            position: "sticky",
            top: 0,
            background: "var(--bg)",
            zIndex: 2,
          }}
        >
          <h6 style={{ marginTop: 0, marginBottom: 0 }}>
            Filter by provider: <a href="/sources" class="button">All</a> |{" "}
            {Object.values(ReleaseSourceProvider).map((item, i, items) => (
              <>
                <a href={`?provider=${item}`} class="button">{item}</a>
                {i + 1 < items.length ? ` | ` : undefined}
              </>
            ))}
          </h6>
        </section>

        <section>
          {sources.map((item, i) => (
            <SourceDisplayItems
              key={i}
              source={item}
            />
          ))}
        </section>
      </main>
    </Page.Body>
  </Page>
);

export const indexPage = (props: IndexPageProps) => <IndexPage {...props} />;
