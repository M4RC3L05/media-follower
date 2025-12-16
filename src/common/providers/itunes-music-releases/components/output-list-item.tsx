import type { FunctionComponent } from "preact";
import { ITunesLookupEntityType, type Output } from "../types.ts";
import { EInputProvider } from "#src/common/database/enums/input-provider.ts";

export type OutputListItemProps = {
  output: Output;
};

export const OutputListItem: FunctionComponent<OutputListItemProps> = (
  { output },
) => {
  const image = output.artworkUrl100
    .split("/")
    .map((segment, index, array) =>
      index === array.length - 1 ? "512x512bb.jpg" : segment
    ).join("/");
  const title = output.wrapperType === "collection"
    ? output.collectionName
    : output.trackName;
  const link = output.wrapperType === "collection"
    ? output.collectionViewUrl
    : output.trackViewUrl;

  return (
    <article>
      <h3>
        [{EInputProvider.ITUNES_MUSIC_RELEASE} |{" "}
        {output.wrapperType === "collection"
          ? ITunesLookupEntityType.ALBUM
          : ITunesLookupEntityType.SONG}] {output.artistName} - {title}
      </h3>
      <div>
        <img src={image} />
        <p>
          {output.releaseDate.toDateString()}
          {output.releaseDate > new Date()
            ? (
              <>
                {" "}
                <em>(To be released)</em>
              </>
            )
            : undefined}
        </p>
      </div>
      <a target="_blank" href={link}>View on source</a>
      <details>
        <summary>Raw:</summary>
        <pre>{JSON.stringify(output, null, 2)}</pre>
      </details>
    </article>
  );
};
