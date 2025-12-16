import type { FunctionComponent } from "preact";
import type { Input } from "../types.ts";
import { EInputProvider } from "#src/common/database/enums/mod.ts";

export type InputListItemProps = { input: Input };

export const InputListItem: FunctionComponent<InputListItemProps> = (
  { input },
) => (
  <article>
    <h3>[{EInputProvider.ITUNES_MUSIC_RELEASE}] - {input.artistName}</h3>
    <div>
      <img src={input.extra.artistImage} />
    </div>
    <details>
      <summary>Raw:</summary>
      <pre>{JSON.stringify(input, null, 2)}</pre>
    </details>
  </article>
);
