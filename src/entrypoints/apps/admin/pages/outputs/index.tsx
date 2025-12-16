import type { FunctionComponent } from "preact";
import { Page } from "#src/entrypoints/apps/admin/pages/page.tsx";
import { EInputProvider } from "#src/common/database/enums/input-provider.ts";
import type { IProviderRender } from "#src/common/providers/interfaces.ts";
import type { DbOutputsTable } from "#src/common/database/types.ts";

type IndexPageProps = {
  outputs: DbOutputsTable[];
  providers: Record<EInputProvider, IProviderRender>;
  url: URL;
  paginatio: { page: number; limit: number };
};

const IndexPage: FunctionComponent<IndexPageProps> = (
  { outputs, providers, url, paginatio: { page } },
) => {
  const prevPageLink = new URL(url);
  prevPageLink.searchParams.set("page", String(Math.max(page - 1, 0)));

  const nextPageLink = new URL(url);
  nextPageLink.searchParams.set("page", String(page + 1));

  return (
    <Page>
      <Page.Head>
        <title>Media Follower | Outputs</title>
      </Page.Head>
      <Page.Body>
        <header>
          <h1>Outputs</h1>
          <nav>
            <a href="/">Back to home</a>
          </nav>
        </header>

        <main>
          <section class="filters-section">
            <h6 class="filters-section__item">
              Filter by provider:<br />
              <a href="/outputs" class="button">All</a> |{" "}
              {Object.values(EInputProvider).map((value, i, items) => (
                <>
                  <a href={`?provider=${value}`} class="button">{value}</a>
                  {i < items.length - 1 ? " | " : ""}
                </>
              ))}
            </h6>
            <h6 class="filters-section__item">
              Pagination:<br />
              <a href={prevPageLink.toString()} class="button">Prev</a> |{" "}
              <a href={nextPageLink.toString()} class="button">Next</a>
            </h6>
          </section>

          <section>
            {outputs.map((item) =>
              providers[item.provider].renderOutputListItem(item)
            )}
          </section>
        </main>
      </Page.Body>
    </Page>
  );
};

export const indexPage = (props: IndexPageProps) => <IndexPage {...props} />;
