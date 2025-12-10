import z from "@zod/zod";
import { ReleaseType } from "../database/enums/release-type.ts";
import { ReleaseSourceProvider } from "../database/enums/release-source-provider.ts";

export const itunesLookupArtistModelSchema = z.object({
  wrapperType: z.literal("artist"),
  artistType: z.literal("Artist"),
  artistName: z.string(),
  artistLinkUrl: z.url(),
  artistId: z.number(),
  amgArtistId: z.number(),
  primaryGenreName: z.string(),
  primaryGenreId: z.number(),
});

export const itunesLookupArtistModelWithExtraSchema =
  itunesLookupArtistModelSchema.extend({
    extra: z.object({ artistImage: z.url() }),
  });

export type ItunesLookupArtistModel = z.infer<
  typeof itunesLookupArtistModelSchema
>;

export type ItunesLookupArtistModelWithExtra = z.infer<
  typeof itunesLookupArtistModelWithExtraSchema
>;

export const itunesLookupAlbumModelSchema = z.object({
  wrapperType: z.literal("collection"),
  collectionType: z.literal("Album"),
  artistId: z.number(),
  collectionId: z.number(),
  amgArtistId: z.number().optional(),
  artistName: z.string(),
  collectionName: z.string(),
  collectionCensoredName: z.string(),
  artistViewUrl: z.string(),
  collectionViewUrl: z.string(),
  artworkUrl60: z.string(),
  artworkUrl100: z.string(),
  collectionPrice: z.number().optional(),
  collectionExplicitness: z.string(),
  trackCount: z.number(),
  copyright: z.string().optional(),
  country: z.string(),
  currency: z.string(),
  releaseDate: z.union([
    z.date(),
    z.string().pipe(z.coerce.date()).pipe(z.date()),
  ]),
  primaryGenreName: z.string(),
  contentAdvisoryRating: z.string().optional(),
});

export type ItunesLookupAlbumModel = z.infer<
  typeof itunesLookupAlbumModelSchema
>;

export const itunesLookupAlbumModelWithExtraSchema =
  itunesLookupAlbumModelSchema
    .extend({
      extra: z.object({
        type: z.literal(ReleaseType.ALBUM),
        provider: z.literal(ReleaseSourceProvider.ITUNES),
      }),
    });

export type ItunesLookupAlbumModelWithExtra = z.infer<
  typeof itunesLookupAlbumModelWithExtraSchema
>;

export const itunesLookupSongModelSchema = z.object({
  wrapperType: z.literal("track"),
  kind: z.literal("song"),
  artistId: z.number(),
  collectionId: z.number(),
  trackId: z.number(),
  artistName: z.string(),
  collectionName: z.string(),
  trackName: z.string(),
  collectionCensoredName: z.string(),
  trackCensoredName: z.string(),
  collectionArtistName: z.string().optional(),
  artistViewUrl: z.string(),
  collectionViewUrl: z.string(),
  trackViewUrl: z.string(),
  previewUrl: z.string(),
  artworkUrl30: z.string(),
  artworkUrl60: z.string(),
  artworkUrl100: z.string(),
  releaseDate: z.union([
    z.date(),
    z.string().pipe(z.coerce.date()).pipe(z.date()),
  ]),
  collectionExplicitness: z.string(),
  trackExplicitness: z.string(),
  discCount: z.number(),
  discNumber: z.number(),
  trackCount: z.number(),
  trackNumber: z.number(),
  trackTimeMillis: z.number(),
  country: z.string(),
  currency: z.string(),
  primaryGenreName: z.string(),
  isStreamable: z.boolean(),
  collectionPrice: z.number().optional(),
  trackPrice: z.number().optional(),
  contentAdvisoryRating: z.string().optional(),
  collectionArtistId: z.number().optional(),
});

export type ItunesLookupSongModel = z.infer<
  typeof itunesLookupSongModelSchema
>;

export const itunesLookupSongModelWithExtraSchema = itunesLookupSongModelSchema
  .extend({
    extra: z.object({
      type: z.literal(ReleaseType.SONG),
      provider: z.literal(ReleaseSourceProvider.ITUNES),
    }),
  });

export type ItunesLookupSongModelWithExtra = z.infer<
  typeof itunesLookupSongModelWithExtraSchema
>;

export type ItunesResponseModel<T> = {
  resultCount: number;
  results: T[];
};

export enum ITunesLookupEntityType {
  SONG = "song",
  ALBUM = "album",
}

export type ItunesLookupType<R extends ITunesLookupEntityType> = R extends
  ITunesLookupEntityType.SONG ? ItunesLookupSongModel
  : R extends ITunesLookupEntityType.ALBUM ? ItunesLookupAlbumModel
  : never;

export type ItunesLookupTypeWithExtra<R extends ITunesLookupEntityType> =
  R extends ITunesLookupEntityType.SONG ? ItunesLookupSongModelWithExtra
    : R extends ITunesLookupEntityType.ALBUM ? ItunesLookupAlbumModelWithExtra
    : never;

export interface IItunesService {
  lookupArtistById(
    id: number,
  ): Promise<ItunesLookupArtistModelWithExtra | undefined>;
  lookupLatestReleasesByArtist<E extends ITunesLookupEntityType>(
    artistId: string,
    entity: E,
    limit: number,
  ): Promise<Array<ItunesLookupTypeWithExtra<E>>>;
}

export type TBluRayComRelease = {
  id: number;
  casing: string;
  artworkurl: string;
  title_sort: string;
  title: string;
  edition: string;
  extended: string;
  title_keywords: string;
  studio: string;
  year: string;
  yearend: string;
  releasedate: string;
  popularity: number;
  width: number;
  height: number;
};

export const bluRayComReleaseSchema = z.object({
  id: z.number(),
  casing: z.string().trim().transform((val) => val === "" ? undefined : val)
    .optional(),
  artworkurl: z.string().trim().transform((val) => val === "" ? undefined : val)
    .optional(),
  title_sort: z.string().trim(),
  title: z.string().trim(),
  edition: z.string().trim().transform((val) => val === "" ? undefined : val)
    .optional(),
  extended: z.string().trim().transform((val) => val === "" ? undefined : val)
    .optional(),
  title_keywords: z.string().trim(),
  studio: z.string().trim().transform((val) => val === "" ? undefined : val)
    .optional(),
  year: z.union([
    z.number(),
    z.string().trim().pipe(z.coerce.number()).pipe(z.number()),
  ]),
  yearend: z.union([
    z.number(),
    z.string().trim().pipe(z.coerce.number()).pipe(z.number()),
  ]),
  releasedate: z.union([
    z.date(),
    z.string().trim().pipe(z.coerce.date()).pipe(z.date()),
  ]),
  popularity: z.number(),
  width: z.number(),
  height: z.number(),
});

export type BluRayComRelease = z.infer<typeof bluRayComReleaseSchema>;

export const bluRayComReleaseWithExtraSchema = bluRayComReleaseSchema.extend({
  extra: z.object({
    country: z.string(),
    artworkUrl: z.url(),
    type: z.enum([ReleaseType.BLURAY, ReleaseType.DVD]),
    link: z.url(),
    provider: z.literal(ReleaseSourceProvider.BLU_RAY_COM),
  }),
});

export type BluRayComReleaseWithExtra = z.infer<
  typeof bluRayComReleaseWithExtraSchema
>;

export const bluRayComCountrySchema = z.object({
  code: z.string(),
  name: z.string(),
});

export type BluRayComCountry = z.infer<typeof bluRayComCountrySchema>;

export interface IBlurayComService {
  getCountries(): Promise<Array<BluRayComCountry>>;
  getBlurayReleasesByCountryForMonth(
    country: string,
    year: number,
    month: number,
  ): Promise<BluRayComReleaseWithExtra[]>;
}

export type ProviderRelease =
  | ItunesLookupAlbumModelWithExtra
  | ItunesLookupSongModelWithExtra
  | BluRayComReleaseWithExtra;
