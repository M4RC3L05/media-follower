import { Page } from "./page.tsx";

const IndexPage = () => (
  <Page>
    <Page.Head>
      <title>Admin | Index</title>
    </Page.Head>
    <Page.Body>
      <h1>Media Follower</h1>
      <a href="/sources">Go to sources</a> |{" "}
      <a href="/releases">Go to releases</a>
    </Page.Body>
  </Page>
);

export const indexPage = () => <IndexPage />;
