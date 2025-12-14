import { FunctionComponent } from "preact";
import { Input } from "../types.ts";
import { EInputProvider } from "../../../database/enums/mod.ts";

export type InputListItemProps = { input: Input };

export const InputListItem: FunctionComponent<InputListItemProps> = (
  { input },
) => (
  <article>
    <h3>[{EInputProvider.ITUNES_MUSIC_RELEASE}] - {input.artistName}</h3>
    {
      <img
        style={{ aspectRatio: "1/1", maxWidth: "256px" }}
        src={input.extra.artistImage}
      />
    }
    <details>
      <summary>Raw:</summary>
      <pre>{JSON.stringify(input, null, 2)}</pre>
    </details>
  </article>
);
