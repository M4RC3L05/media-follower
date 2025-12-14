import { Page } from "./page.tsx";

const ErrorPage = () => (
  <Page>
    <Page.Head>
      <title>Music Follower | Error</title>
    </Page.Head>
    <Page.Body>
      <h1>Something went wrong</h1>
    </Page.Body>
  </Page>
);

export const errorPage = () => <ErrorPage />;
