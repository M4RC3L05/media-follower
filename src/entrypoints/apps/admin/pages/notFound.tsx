import { Page } from "./page.tsx";

const NotFoundPage = () => (
  <Page>
    <Page.Head>
      <title>Music Follower | 404</title>
    </Page.Head>
    <Page.Body>
      <h1>Page not found</h1>
    </Page.Body>
  </Page>
);

export const notFoundPage = () => <NotFoundPage />;
