import type { FunctionComponent } from "preact";
import type { Output } from "../types.ts";
import type { DbOutputsTable } from "#src/common/database/types.ts";
import { EInputProvider } from "#src/common/database/enums/mod.ts";

export type OutputListItemProps = {
  output: Output;
  outputRow: DbOutputsTable;
};

export const OutputListItem: FunctionComponent<OutputListItemProps> = (
  { output, outputRow },
) => {
  const title =
    `[${EInputProvider.BLU_RAY_COM_PHYSICAL_RELEASE} | ${output.extra.type} | ${outputRow.input_id.toUpperCase()}] ${output.title}`;
  return (
    <article>
      <h3>{title}</h3>
      <div style={{ textAlign: "center" }}>
        <img
          style={{ aspectRatio: "9/16", maxWidth: "256px" }}
          src={output.extra.artworkUrl}
        />
        <p>
          {output.releasedate.toDateString()}
          {output.releasedate > new Date()
            ? (
              <>
                {" "}
                <em>(To be released)</em>
              </>
            )
            : undefined}
        </p>
      </div>
      <a target="_blank" href={output.extra.link}>View on source</a>
      <details>
        <summary>Raw:</summary>
        <pre>{JSON.stringify(output, null, 2)}</pre>
      </details>
    </article>
  );
};
