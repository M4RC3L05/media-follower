import { routes } from "../router.ts";
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
          <a href="/inputs">Go to inputs</a>
          <a href="/outputs">Go to outputs</a>
          <form method="post" action={routes.auth.logout.href()}>
            <input type="submit" value="Logout" />
          </form>
        </nav>
      </header>
    </Page.Body>
  </Page>
);

export const indexPage = () => <IndexPage />;
