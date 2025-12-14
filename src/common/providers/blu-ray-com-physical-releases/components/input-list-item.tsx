import { FunctionComponent } from "preact";
import { Input } from "../types.ts";
import { EInputProvider } from "../../../database/enums/mod.ts";

export type InputListItemProps = { input: Input };

export const InputListItem: FunctionComponent<InputListItemProps> = (
  { input },
) => (
  <article>
    <h3>[{EInputProvider.BLU_RAY_COM_PHYSICAL_RELEASE}] - {input.name}</h3>
    <details>
      <summary>Raw:</summary>
      <pre>{JSON.stringify(input, null, 2)}</pre>
    </details>
  </article>
);
