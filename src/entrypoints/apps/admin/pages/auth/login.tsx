import type { FunctionComponent } from "preact";
import { Page } from "../../pages/page.tsx";
import { routes } from "../../router.ts";

const LoginPage: FunctionComponent = () => {
  return (
    <Page>
      <Page.Head>
        <title>Media Follower | Login</title>
      </Page.Head>
      <Page.Body>
        <header>
          <h1>Login</h1>
          <nav>
            <a href={routes.auth.register.index.href()}>Go to register</a>
          </nav>
        </header>

        <main>
          <section>
            <form method="post" action={routes.auth.login.action.href()}>
              <div>
                <label for="username">Username</label>
                <input name="username" id="username" />
              </div>

              <div>
                <label for="password">Password</label>
                <input name="password" id="password" type="password" />
              </div>

              <input type="submit" value="Login" />
            </form>
          </section>
        </main>
      </Page.Body>
    </Page>
  );
};

export const loginPage = () => <LoginPage />;
