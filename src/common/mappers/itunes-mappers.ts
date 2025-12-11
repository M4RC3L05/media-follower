import { ReleaseSourceProvider } from "../database/enums/release-source-provider.ts";
import { ReleaseType } from "../database/enums/release-type.ts";
import { DbReleaseSourcesTable, DbReleasesTable } from "../database/types.ts";
import {
  type ItunesLookupAlbumModelWithExtra,
  itunesLookupAlbumModelWithExtraSchema,
  type ItunesLookupArtistModel,
  itunesLookupArtistModelSchema,
  type ItunesLookupSongModel,
  type ItunesLookupSongModelWithExtra,
  itunesLookupSongModelWithExtraSchema,
} from "../services/service.ts";

export const fromReleaseSourceToPersistance = (
  item: ItunesLookupArtistModel,
): DbReleaseSourcesTable => {
  return {
    id: crypto.randomUUID(),
    provider: ReleaseSourceProvider.ITUNES,
    raw: JSON.stringify(item),
  };
};

export const fromReleaseToPersistance = (
  item: ItunesLookupSongModelWithExtra | ItunesLookupAlbumModelWithExtra,
): DbReleasesTable => {
  return {
    releasedAt: item.releaseDate.toISOString(),
    provider: item.extra.provider,
    type: item.extra.type,
    raw: JSON.stringify(item),
    id: String(
      item.wrapperType === "collection"
        ? item.collectionId
        : (item as ItunesLookupSongModel).trackId,
    ),
  };
};

export const fromPersistanceToRelease = (
  row: DbReleasesTable,
): ItunesLookupAlbumModelWithExtra | ItunesLookupSongModelWithExtra => {
  if (row.provider !== ReleaseSourceProvider.ITUNES) {
    throw new Error(`Release provider "${row.provider}" not supported`);
  }

  switch (row.type) {
    case ReleaseType.SONG: {
      return itunesLookupSongModelWithExtraSchema.parse(JSON.parse(row.raw));
    }
    case ReleaseType.ALBUM: {
      return itunesLookupAlbumModelWithExtraSchema.parse(JSON.parse(row.raw));
    }
    default: {
      throw new Error(`Relase type "${row.type}" not supported`);
    }
  }
};

export const fromPersistanceToReleaseSurce = (
  row: DbReleaseSourcesTable,
): ItunesLookupArtistModel => {
  return itunesLookupArtistModelSchema.parse(JSON.parse(row.raw));
};
