import z from "@zod/zod";

export const itunesMusicReleasesInputSchema = z.object({
  wrapperType: z.literal("artist"),
  artistType: z.literal("Artist"),
  artistName: z.string(),
  artistLinkUrl: z.url(),
  artistId: z.number(),
  amgArtistId: z.number().optional(),
  primaryGenreName: z.string(),
  primaryGenreId: z.number(),
});

export const itunesMusicReleasesInputWithExtraSchema =
  itunesMusicReleasesInputSchema
    .extend({
      extra: z.object({ artistImage: z.url() }),
    });

export type ItunesMusicReleasesInput = z.infer<
  typeof itunesMusicReleasesInputSchema
>;

export type ItunesMusicReleasesInputWithExtra = z.infer<
  typeof itunesMusicReleasesInputWithExtraSchema
>;

export const itunesMusicReleasesOutputAlbumSchema = z.object({
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
  ]).optional().default(() => new Date(new Date().toISOString())),
  primaryGenreName: z.string(),
  contentAdvisoryRating: z.string().optional(),
});

export type ItunesMusicReleasesOutputAlbum = z.infer<
  typeof itunesMusicReleasesOutputAlbumSchema
>;

export const itunesMusicReleasesOutputSongSchema = z.object({
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
  artistViewUrl: z.string().optional(),
  collectionViewUrl: z.string(),
  trackViewUrl: z.string(),
  previewUrl: z.string().optional(),
  artworkUrl30: z.string(),
  artworkUrl60: z.string(),
  artworkUrl100: z.string(),
  releaseDate: z.union([
    z.date(),
    z.string().pipe(z.coerce.date()).pipe(z.date()),
  ]).optional().default(() => new Date(new Date().toISOString())),
  collectionExplicitness: z.string(),
  trackExplicitness: z.string(),
  discCount: z.number(),
  discNumber: z.number(),
  trackCount: z.number(),
  trackNumber: z.number(),
  trackTimeMillis: z.number().optional(),
  country: z.string(),
  currency: z.string(),
  primaryGenreName: z.string(),
  isStreamable: z.boolean(),
  collectionPrice: z.number().optional(),
  trackPrice: z.number().optional(),
  contentAdvisoryRating: z.string().optional(),
  collectionArtistId: z.number().optional(),
});

export type ItunesMusicReleasesOutputSong = z.infer<
  typeof itunesMusicReleasesOutputSongSchema
>;

export enum ITunesLookupEntityType {
  SONG = "song",
  ALBUM = "album",
}

export type ItunesMusicReleasesOutput<R extends ITunesLookupEntityType> =
  R extends ITunesLookupEntityType.SONG ? ItunesMusicReleasesOutputSong
    : R extends ITunesLookupEntityType.ALBUM ? ItunesMusicReleasesOutputAlbum
    : never;

export type ItunesResponseModel<T> = {
  resultCount: number;
  results: T[];
};

export type Input = ItunesMusicReleasesInputWithExtra;
export type Output =
  | ItunesMusicReleasesOutputAlbum
  | ItunesMusicReleasesOutputSong;
