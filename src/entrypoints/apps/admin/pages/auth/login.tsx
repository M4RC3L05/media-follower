import type { FunctionComponent } from "preact";
import { Page } from "../../pages/page.tsx";

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
            <a href="/auth/register">Go to register</a>
          </nav>
        </header>

        <main>
          <section>
            <form method="post" action="/auth/login">
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
