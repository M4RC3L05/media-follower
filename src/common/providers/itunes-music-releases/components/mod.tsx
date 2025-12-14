import { InputListItem, type InputListItemProps } from "./input-list-item.tsx";
import {
  OutputListItem,
  type OutputListItemProps,
} from "./output-list-item.tsx";

export const inputListItem = (props: InputListItemProps) => (
  <InputListItem {...props} />
);

export const outputListItem = (props: OutputListItemProps) => (
  <OutputListItem {...props} />
);
