import type { FunctionComponent } from "preact";
import { Page } from "../page.tsx";

const RegisterPage: FunctionComponent = () => {
  return (
    <Page>
      <Page.Head>
        <title>Media Follower | Register</title>
      </Page.Head>
      <Page.Body>
        <header>
          <h1>Register</h1>
          <nav>
            <a href="/auth/login">Go to login</a>
          </nav>
        </header>

        <main>
          <section>
            <form method="post" action="/auth/register">
              <div>
                <label for="username">Username</label>
                <input name="username" id="username" />
              </div>

              <div>
                <label for="password">Password</label>
                <input name="password" id="password" type="password" />
              </div>

              <input type="submit" value="Register" />
            </form>
          </section>
        </main>
      </Page.Body>
    </Page>
  );
};

export const registerPage = () => <RegisterPage />;
