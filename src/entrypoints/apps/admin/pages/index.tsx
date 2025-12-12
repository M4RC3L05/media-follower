import { Page } from "./page.tsx";

const IndexPage = () => (
  <Page>
    <Page.Head>
      <title>Admin | Index</title>
    </Page.Head>
    <Page.Body>
      <header>
        <h1>Media Follower</h1>
        <nav>
          <a href="/sources">Sources</a>
          <a href="/releases">Releases</a>
        </nav>
      </header>
    </Page.Body>
  </Page>
);

export const indexPage = () => <IndexPage />;
