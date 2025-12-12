import type { FunctionComponent } from "preact";
import { ReleaseSourceProvider } from "../../../../../common/database/enums/release-source-provider.ts";
import { Page } from "../page.tsx";

const BluRayComSourceCreateForm = () => (
  <>
    <h3>Add "{ReleaseSourceProvider.BLU_RAY_COM}" source</h3>

    <form method="post" action="/sources/create">
      <input
        type="hidden"
        value={ReleaseSourceProvider.BLU_RAY_COM}
        name="provider"
      />

      <div>
        <label for="country">Country:</label>
        <input id="country" name="country" />
      </div>

      <button type="submit">Add</button>
    </form>
  </>
);

const ItunesSourceCreateForm = () => (
  <>
    <h3>Add "{ReleaseSourceProvider.ITUNES}" source</h3>

    <form method="post" action="/sources/create">
      <input
        type="hidden"
        value={ReleaseSourceProvider.ITUNES}
        name="provider"
      />

      <div>
        <label for="artistId">Artist ID:</label>
        <input id="artistId" name="artistId" />
      </div>

      <button type="submit">Add</button>
    </form>
  </>
);

const CreateFormItems = () =>
  Object.entries(ReleaseSourceProvider).map(([_, value]) => {
    switch (value) {
      case ReleaseSourceProvider.BLU_RAY_COM: {
        return <BluRayComSourceCreateForm />;
      }
      case ReleaseSourceProvider.ITUNES: {
        return <ItunesSourceCreateForm />;
      }
    }
  });

const CreatePage: FunctionComponent = () => (
  <Page>
    <Page.Head>
      <title>Media follower | Create source</title>
    </Page.Head>
    <Page.Body>
      <header>
        <h1>Create Source</h1>
      </header>

      <main>
        <CreateFormItems />
      </main>
    </Page.Body>
  </Page>
);

export const createPage = () => <CreatePage />;
