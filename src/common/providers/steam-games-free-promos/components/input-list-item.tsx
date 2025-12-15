import type { FunctionComponent } from "preact";
import type { Input } from "../types.ts";
import { EInputProvider } from "#src/common/database/enums/mod.ts";

export type InputListItemProps = { input: Input };

export const InputListItem: FunctionComponent<InputListItemProps> = (
  { input },
) => (
  <article>
    <h3>[{EInputProvider.STEAM_GAMES_FREE_PROMOS}] - {input.url}</h3>
    <details>
      <summary>Raw:</summary>
      <pre>{JSON.stringify(input, null, 2)}</pre>
    </details>
  </article>
);
