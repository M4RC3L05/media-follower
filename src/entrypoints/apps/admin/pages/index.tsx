import { Page } from "#src/entrypoints/apps/admin/pages/page.tsx";

const IndexPage = () => (
  <Page>
    <Page.Head>
      <title>Admin | Index</title>
    </Page.Head>
    <Page.Body>
      <header>
        <h1>Media Follower</h1>
        <nav>
          <a href="/inputs">Go to inputs</a>
          <a href="/outputs">Go to outputs</a>
        </nav>
      </header>
    </Page.Body>
  </Page>
);

export const indexPage = () => <IndexPage />;
