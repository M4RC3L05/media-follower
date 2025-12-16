import type { FunctionComponent } from "preact";
import { Page } from "./page.tsx";

const ErrorPage: FunctionComponent<{ message?: string | undefined }> = (
  { message },
) => (
  <Page>
    <Page.Head>
      <title>Music Follower | Error</title>
    </Page.Head>
    <Page.Body>
      <h1>{message ?? "Something went wrong"}</h1>
    </Page.Body>
  </Page>
);

export const errorPage = (props?: { message: string | undefined }) => (
  <ErrorPage {...props} />
);
