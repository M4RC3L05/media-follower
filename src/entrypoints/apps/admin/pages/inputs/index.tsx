import type { FunctionComponent } from "preact";
import { Page } from "#src/entrypoints/apps/admin/pages/page.tsx";
import { EInputProvider } from "#src/common/database/enums/input-provider.ts";
import type { IProviderRender } from "#src/common/providers/interfaces.ts";
import type { DbInputsTable } from "#src/common/database/types.ts";

type IndexPageProps = {
  sources: DbInputsTable[];
  providers: Record<EInputProvider, IProviderRender>;
  url: URL;
  paginatio: { page: number; limit: number };
};

const IndexPage: FunctionComponent<IndexPageProps> = (
  { sources, providers, url, paginatio: { page } },
) => {
  const prevPageLink = new URL(url);
  prevPageLink.searchParams.set("page", String(Math.max(page - 1, 0)));

  const nextPageLink = new URL(url);
  nextPageLink.searchParams.set("page", String(page + 1));

  return (
    <Page>
      <Page.Head>
        <title>Media Follower | Inputs</title>
      </Page.Head>
      <Page.Body>
        <header>
          <h1>Inputs</h1>
          <nav>
            <a href="/">Back to home</a>
            <a href="/inputs/create">Add a new input</a>
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
            <h6 style={{ marginTop: 0, marginBottom: 0, textAlign: "center" }}>
              Filter by provider:<br />
              <a href="/inputs" class="button">All</a> |{" "}
              {Object.values(EInputProvider).map((item, i, items) => (
                <>
                  <a href={`?provider=${item}`} class="button">{item}</a>
                  {i + 1 < items.length ? ` | ` : undefined}
                </>
              ))}
            </h6>
            <h6 style={{ marginTop: 0, marginBottom: 0, textAlign: "center" }}>
              Pagination:<br />
              <a href={prevPageLink.toString()} class="button">Prev</a> |{" "}
              <a href={nextPageLink.toString()} class="button">Next</a>
            </h6>
          </section>

          <section>
            {sources.map((item) =>
              providers[item.provider].renderInputListItem(item)
            )}
          </section>
        </main>
      </Page.Body>
    </Page>
  );
};

export const indexPage = (props: IndexPageProps) => <IndexPage {...props} />;
