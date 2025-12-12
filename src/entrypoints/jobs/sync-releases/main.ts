import z from "@zod/zod";
import { parseArgs } from "@std/cli";
import { config, initConfig } from "#src/common/config/mod.ts";
import { BluRayComService } from "#src/common/services/blu-ray-com-service.ts";
import { HttpFetch } from "#src/common/http/mod.ts";
import { gracefulShutdown } from "#src/common/process/mod.ts";
import { ItunesService } from "#src/common/services/itunes-service.ts";
import { ReleaseSourceProvider } from "#src/common/database/enums/release-source-provider.ts";
import { CustomDatabase } from "#src/common/database/mod.ts";
import { App, type AppProps } from "#src/entrypoints/jobs/sync-releases/app.ts";
import {
  bluRayComCountrySchema,
  itunesLookupArtistModelWithExtraSchema,
  ITunesLookupEntityType,
} from "#src/common/services/service.ts";
import { bluRayComMappers, itunesMappers } from "#src/common/mappers/mod.ts";
import { ReleaseType } from "../../../common/database/enums/release-type.ts";

initConfig();

const { done, signal: shutdownSignal } = gracefulShutdown();

const { provider } = z.object({ provider: z.enum(ReleaseSourceProvider) })
  .parse(parseArgs(Deno.args));

using db = new CustomDatabase(config().database.path);
const httpClient = new HttpFetch({ signal: shutdownSignal });
const now = new Date();
const bluRayComService = new BluRayComService({ httpClient: httpClient });
const itunesService = new ItunesService({ httpClient: httpClient });

const resolveService = (
  provider: ReleaseSourceProvider,
): AppProps["service"] => {
  switch (provider) {
    case ReleaseSourceProvider.BLU_RAY_COM: {
      return {
        fetchReleasesFromSource: async (source) => {
          const parsed = bluRayComCountrySchema.parse(JSON.parse(source.raw));
          const [blurays, dvds] = await Promise.all([
            bluRayComService
              .getBlurayReleasesByCountryForMonth(
                parsed.code,
                now.getFullYear(),
                now.getMonth() + 1,
                ReleaseType.BLURAY,
              ),
            bluRayComService.getBlurayReleasesByCountryForMonth(
              parsed.code,
              now.getFullYear(),
              now.getMonth() + 1,
              ReleaseType.DVD,
            ),
          ]);

          return [
            ...blurays.map((item) =>
              bluRayComMappers.fromReleaseToPersistance(item)
            ),
            ...dvds.map((item) =>
              bluRayComMappers.fromReleaseToPersistance(item)
            ),
          ];
        },
      };
    }
    case ReleaseSourceProvider.ITUNES: {
      return {
        fetchReleasesFromSource: async (source) => {
          const parsed = itunesLookupArtistModelWithExtraSchema.parse(
            JSON.parse(source.raw),
          );

          const [albums, songs] = await Promise.all([
            itunesService.lookupLatestReleasesByArtist(
              String(parsed.artistId),
              ITunesLookupEntityType.ALBUM,
              50,
            ),
            itunesService.lookupLatestReleasesByArtist(
              String(parsed.artistId),
              ITunesLookupEntityType.SONG,
              50,
            ),
          ]);

          return [
            ...albums.map((item) =>
              itunesMappers.fromReleaseToPersistance(item)
            ),
            ...songs.map((item) =>
              itunesMappers.fromReleaseToPersistance(item)
            ),
          ];
        },
      };
    }
  }
};

await new App({
  database: db,
  provider,
  service: resolveService(provider),
  signal: shutdownSignal,
}).execute();

await done();
