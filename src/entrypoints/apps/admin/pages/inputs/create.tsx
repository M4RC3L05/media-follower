import type { FunctionComponent } from "preact";
import { EInputProvider } from "#src/common/database/enums/input-provider.ts";
import { Page } from "../page.tsx";
import { routes } from "../../router.ts";

const CreatePage: FunctionComponent = () => (
  <Page>
    <Page.Head>
      <title>Media follower | Create input</title>
    </Page.Head>
    <Page.Body>
      <header>
        <h1>Create Input</h1>
      </header>

      <main>
        <section>
          <form method="post" action={routes.inputs.create.action.href()}>
            <select name="provider" id="provider">
              {Object.values(EInputProvider).map((item) => (
                <option value={item}>
                  {item.replaceAll("-", " ")}
                </option>
              ))}
            </select>

            <div>
              <label for="term">Term:</label>
              <input id="term" name="term" />
            </div>

            <button type="submit">Create</button>
          </form>
        </section>
      </main>
    </Page.Body>
  </Page>
);

export const createPage = () => <CreatePage />;
