import type { FunctionComponent } from "preact";
import { EInputProvider } from "../../../../../common/database/enums/input-provider.ts";
import { Page } from "#src/entrypoints/apps/admin/pages/page.tsx";

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
        <form method="post" action="/inputs/create">
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
      </main>
    </Page.Body>
  </Page>
);

export const createPage = () => <CreatePage />;
