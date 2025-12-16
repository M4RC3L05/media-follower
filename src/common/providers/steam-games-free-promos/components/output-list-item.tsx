import type { FunctionComponent } from "preact";
import type { Output } from "../types.ts";
import { EInputProvider } from "#src/common/database/enums/mod.ts";

export type OutputListItemProps = {
  output: Output;
};

export const OutputListItem: FunctionComponent<OutputListItemProps> = (
  { output },
) => {
  const title =
    `[${EInputProvider.STEAM_GAMES_FREE_PROMOS} | ${output.promoType}] ${output.name}`;
  return (
    <article>
      <h3>{title}</h3>
      <img src={output.image} />
      <p>{output.promoType}</p>
      <p>
        From {output.startDate.toDateString()} to{" "}
        {output.endDate.toDateString()}
      </p>
      <a target="_blank" href={output.link}>View on source</a>
      <details>
        <summary>Raw:</summary>
        <pre>{JSON.stringify(output, null, 2)}</pre>
      </details>
    </article>
  );
};
