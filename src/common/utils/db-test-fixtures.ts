import type { IDatabase } from "#src/common/database/database.ts";
import { ReleaseSourceProvider } from "#src/common/database/enums/release-source-provider.ts";
import type {
  DbReleaseSourcesTable,
  DbReleasesTable,
} from "#src/common/database/types.ts";
import {
  type BluRayComCountry,
  bluRayComCountrySchema,
  type BluRayComReleaseWithExtra,
  bluRayComReleaseWithExtraSchema,
  type ItunesLookupAlbumModelWithExtra,
  itunesLookupAlbumModelWithExtraSchema,
  type ItunesLookupArtistModelWithExtra,
  itunesLookupArtistModelWithExtraSchema,
  type ItunesLookupSongModelWithExtra,
  itunesLookupSongModelWithExtraSchema,
} from "#src/common/services/service.ts";
import { ReleaseType } from "../database/enums/release-type.ts";

type DeepPartial<T> = {
  [P in keyof T]?: T[P] extends object ? DeepPartial<T[P]>
    : T[P];
};

export const generateReleaseBluRayCom = (
  data?: DeepPartial<BluRayComReleaseWithExtra>,
) =>
  bluRayComReleaseWithExtraSchema.parse({
    id: data?.id ?? 1,
    title_sort: data?.title_sort ?? "foo",
    title: data?.title ?? "bar",
    title_keywords: data?.title_keywords ?? "biz",
    year: data?.year ?? 2,
    yearend: data?.yearend ?? 3,
    releasedate: data?.releasedate ?? new Date().toISOString(),
    popularity: data?.popularity ?? 4,
    width: data?.width ?? 5,
    height: data?.height ?? 6,
    extra: {
      ...data?.extra,
      type: data?.extra?.type ?? "bluray",
      provider: data?.extra?.provider ?? "blu-ray-com",
      country: data?.extra?.country ?? "buz",
      artworkUrl: data?.extra?.artworkUrl ??
        "https://example.com/325994_medium.jpg",
      link: data?.extra?.link ?? "https://example.com/325994/",
    },
  });

export const generateReleaseSongItunes = (
  data?: DeepPartial<ItunesLookupSongModelWithExtra>,
) =>
  itunesLookupSongModelWithExtraSchema.parse({
    wrapperType: "track",
    kind: "song",
    artistId: data?.artistId ?? 1,
    collectionId: data?.collectionId ?? 2,
    trackId: data?.trackId ?? 3,
    artistName: data?.artistName ?? "Example Artist",
    collectionName: data?.collectionName ?? "Example Album",
    trackName: data?.trackName ?? "Example Song",
    collectionCensoredName: data?.collectionCensoredName ?? "Example Album",
    trackCensoredName: data?.trackCensoredName ?? "Example Song",
    artistViewUrl: data?.artistViewUrl ?? "https://example.com/artist",
    collectionViewUrl: data?.collectionViewUrl ?? "https://example.com/album",
    trackViewUrl: data?.trackViewUrl ?? "https://example.com/song",
    previewUrl: data?.previewUrl ?? "https://example.com/preview",
    artworkUrl30: data?.artworkUrl30 ?? "https://example.com/artwork30.jpg",
    artworkUrl60: data?.artworkUrl60 ?? "https://example.com/artwork60.jpg",
    artworkUrl100: data?.artworkUrl100 ?? "https://example.com/artwork100.jpg",
    releaseDate: data?.releaseDate ?? new Date().toISOString(),
    collectionExplicitness: data?.collectionExplicitness ?? "notExplicit",
    trackExplicitness: data?.trackExplicitness ?? "notExplicit",
    discCount: data?.discCount ?? 1,
    discNumber: data?.discNumber ?? 1,
    trackCount: data?.trackCount ?? 10,
    trackNumber: data?.trackNumber ?? 1,
    trackTimeMillis: data?.trackTimeMillis ?? 180000,
    country: data?.country ?? "US",
    currency: data?.currency ?? "USD",
    primaryGenreName: data?.primaryGenreName ?? "Pop",
    isStreamable: data?.isStreamable ?? true,
    extra: {
      ...data?.extra,
      type: data?.extra?.type ?? "song",
      provider: data?.extra?.provider ?? "itunes",
    },
  });

export const generateReleaseAlbumItunes = (
  data?: DeepPartial<ItunesLookupAlbumModelWithExtra>,
) =>
  itunesLookupAlbumModelWithExtraSchema.parse({
    wrapperType: "collection",
    collectionType: "Album",
    artistId: data?.artistId ?? 1,
    collectionId: data?.collectionId ?? 2,
    artistName: data?.artistName ?? "Example Artist",
    collectionName: data?.collectionName ?? "Example Album",
    collectionCensoredName: data?.collectionCensoredName ?? "Example Album",
    artistViewUrl: data?.artistViewUrl ?? "https://example.com/artist",
    collectionViewUrl: data?.collectionViewUrl ?? "https://example.com/album",
    artworkUrl60: data?.artworkUrl60 ?? "https://example.com/artwork60.jpg",
    artworkUrl100: data?.artworkUrl100 ?? "https://example.com/artwork100.jpg",
    collectionPrice: data?.collectionPrice ?? 9.99,
    collectionExplicitness: data?.collectionExplicitness ?? "notExplicit",
    trackCount: data?.trackCount ?? 10,
    copyright: data?.copyright ?? "© 2023 Example Artist",
    country: data?.country ?? "US",
    currency: data?.currency ?? "USD",
    releaseDate: data?.releaseDate ?? new Date().toISOString(),
    primaryGenreName: data?.primaryGenreName ?? "Pop",
    extra: {
      ...data?.extra,
      type: data?.extra?.type ?? "album",
      provider: data?.extra?.provider ?? "itunes",
    },
  });

export const generateReleaseSourceBluRayComCountry = (
  data?: DeepPartial<BluRayComCountry>,
) =>
  bluRayComCountrySchema.parse({
    code: data?.code ?? "es",
    name: data?.name ?? "Spain",
  });

export const generateReleaseSourceItunesArtist = (
  data?: DeepPartial<ItunesLookupArtistModelWithExtra>,
) =>
  itunesLookupArtistModelWithExtraSchema.parse({
    wrapperType: data?.wrapperType ?? "artist",
    artistType: data?.artistType ?? "Artist",
    artistName: data?.artistName ?? "Foo bar",
    artistLinkUrl: data?.artistLinkUrl ?? "https://example.com",
    artistId: data?.artistId ?? 1,
    amgArtistId: data?.amgArtistId ?? undefined,
    primaryGenreName: data?.primaryGenreName ?? "Biz",
    primaryGenreId: data?.primaryGenreId ?? 7,
    extra: data?.extra ?? {
      artistImage: data?.extra?.artistImage ??
        "https://example.com/256x256.png",
    },
  });

export const loadReleaseSource = (
  db: IDatabase,
  data?: DeepPartial<DbReleaseSourcesTable>,
) => {
  const provider = data?.provider ?? ReleaseSourceProvider.BLU_RAY_COM;
  const raw = data?.raw ??
    (provider === ReleaseSourceProvider.BLU_RAY_COM
      ? JSON.stringify(generateReleaseSourceBluRayComCountry())
      : JSON.stringify(generateReleaseSourceItunesArtist()));

  const [item] = db.sql<DbReleaseSourcesTable>`
    insert into release_sources
      (id, provider, raw)
    values
      (
        ${data?.id ?? String(Math.floor(Math.random() * 10000000))},
        ${provider},
        jsonb(${raw})
      )
    returning id, provider, json(raw) as raw;
  `;

  if (!item) throw new Error("Could not insert fixture");

  return item;
};

export const generateRelease = (
  data?: DeepPartial<DbReleasesTable>,
): DbReleasesTable => {
  const provider = data?.provider ?? ReleaseSourceProvider.BLU_RAY_COM;
  const type = data?.type ??
    (provider === ReleaseSourceProvider.BLU_RAY_COM
      ? ReleaseType.BLURAY
      : ReleaseType.ALBUM);
  const raw = data?.raw ??
    (provider === ReleaseSourceProvider.BLU_RAY_COM
      ? JSON.stringify(generateReleaseBluRayCom())
      : type === ReleaseType.ALBUM
      ? JSON.stringify(generateReleaseAlbumItunes())
      : JSON.stringify(generateReleaseSongItunes()));

  return {
    id: data?.id ?? String(Math.floor(Math.random() * 10000000)),
    provider,
    type,
    releasedAt: data?.releasedAt ?? new Date().toISOString(),
    raw,
  };
};

export const loadRelease = (
  db: IDatabase,
  data?: DeepPartial<DbReleasesTable>,
): DbReleasesTable => {
  const generated = generateRelease(data);

  const [item] = db.sql<DbReleasesTable>`
    insert into releases
      (id, provider, type, "releasedAt", raw)
    values
      (
        ${generated.id},
        ${generated.provider},
        ${generated.type},
        ${generated.releasedAt},
        jsonb(${generated.raw})
      )
    returning id, provider, type, "releasedAt", json(raw) as raw;
  `;

  if (!item) throw new Error("Could not insert fixture");

  return item;
};
