import {
  afterAll,
  afterEach,
  beforeEach,
  describe,
  it,
} from "@std/testing/bdd";
import { config, initConfig } from "#src/common/config/mod.ts";
import { CustomDatabase } from "#src/common/database/mod.ts";
import { dbTestFixturesUtils, dbTestUtils } from "#src/common/utils/mod.ts";
import { ReleaseSourceProvider } from "#src/common/database/enums/release-source-provider.ts";
import { App } from "./app.ts";
import { assertSpyCall, assertSpyCalls, spy } from "@std/testing/mock";
import { makeLogger } from "#src/common/logger/mod.ts";
import { FakeTime } from "@std/testing/time";
import { assertEquals } from "@std/assert";

initConfig();

let database: CustomDatabase;

describe("App", () => {
  beforeEach(async () => {
    database = new CustomDatabase(config().database.path);
    await dbTestUtils.runMigrations(database);
  });

  afterEach(() => {
    database.exec("delete from release_sources");
  });

  afterAll(() => {
    database.close();
  });

  for (const provider of Object.values(ReleaseSourceProvider)) {
    describe(`Provider: ${provider}`, () => {
      it("should do nothing if signal is already aborted", async () => {
        const abortController = new AbortController();
        abortController.abort();

        const fetchReleasesFromSource = spy(() => Promise.resolve([]));

        await new App({
          database,
          provider,
          service: { fetchReleasesFromSource },
          signal: abortController.signal,
        }).execute();

        assertSpyCalls(fetchReleasesFromSource, 0);
      });

      it("should do nothing no sources available", async () => {
        const fetchReleasesFromSource = spy(() => Promise.resolve([]));

        await new App({
          database,
          service: { fetchReleasesFromSource },
          provider,
          signal: new AbortController().signal,
        }).execute();

        assertSpyCalls(fetchReleasesFromSource, 0);
        assertSpyCalls(fetchReleasesFromSource, 0);
      });

      it("should log error if an error occurs", async () => {
        const error = new Error("foo");
        const source = dbTestFixturesUtils.loadReleaseSource(database, {
          provider,
        });
        const fetchReleasesFromSource = spy(() => {
          throw error;
        });

        using logErrorSpy = spy(makeLogger("sync-releases-app"), "error");

        await new App({
          database,
          service: { fetchReleasesFromSource },
          provider,
          signal: new AbortController().signal,
        }).execute();

        assertSpyCalls(
          fetchReleasesFromSource,
          1,
        );
        assertSpyCalls(logErrorSpy, 1);
        assertSpyCall(logErrorSpy, 0, {
          args: ["Could not sync releases for source successfully", {
            source: { id: source.id },
            error,
          }],
        });
      });

      it("should do nothing if service resolves with no releases", async () => {
        const fetchReleasesFromSource = spy(() => Promise.resolve([]));

        dbTestFixturesUtils.loadReleaseSource(database, { provider });

        await new App({
          database,
          service: { fetchReleasesFromSource },
          provider,
          signal: new AbortController().signal,
        }).execute();

        assertSpyCalls(fetchReleasesFromSource, 1);
        assertEquals(database.sql`select * from releases`.length, 0);
      });

      it("should insert new releases", async () => {
        const fetchReleasesFromSource = spy(() =>
          Promise.resolve([dbTestFixturesUtils.generateRelease({ provider })])
        );

        dbTestFixturesUtils.loadReleaseSource(database, { provider });

        await new App({
          database,
          service: { fetchReleasesFromSource },
          provider,
          signal: new AbortController().signal,
        }).execute();

        assertSpyCalls(fetchReleasesFromSource, 1);
        assertEquals(database.sql`select * from releases;`.length, 1);
      });

      it("should update existing releases if exists", async () => {
        const computeRaw = (name: string) => {
          switch (provider) {
            case ReleaseSourceProvider.BLU_RAY_COM:
              return dbTestFixturesUtils.generateReleaseBluRayCom({
                title: name,
              });
            case ReleaseSourceProvider.ITUNES:
              return dbTestFixturesUtils.generateReleaseSongItunes({
                trackName: name,
              });
          }
        };
        const source = dbTestFixturesUtils.loadReleaseSource(database, {
          provider,
        });
        const noEditRelease = dbTestFixturesUtils.loadRelease(database, {
          provider: source.provider,
          id: "1",
        });
        const release = dbTestFixturesUtils.loadRelease(database, {
          provider: source.provider,
          id: "2",
          raw: JSON.stringify(computeRaw("foo")),
        });

        const toUpdate = dbTestFixturesUtils.generateRelease({
          id: "2",
          provider,
          raw: JSON.stringify(computeRaw("bar")),
        });

        const fetchReleasesFromSource = spy(() => Promise.resolve([toUpdate]));

        await new App({
          database,
          service: { fetchReleasesFromSource },
          provider,
          signal: new AbortController().signal,
        }).execute();

        assertSpyCalls(fetchReleasesFromSource, 1);
        assertEquals(
          database
            .sql`select id, provider, type, "releasedAt", json(raw) as raw from releases where id = ${noEditRelease.id};`[
              0
            ],
          noEditRelease,
        );
        const updated = database
          .sql`select id, provider, type, "releasedAt", json(raw) as raw from releases where id = ${release.id};`[
            0
          ]!;

        assertEquals(updated.raw !== release.raw, true);
        assertEquals(updated.raw, toUpdate.raw);
      });

      it("should delay if more sources exists", async () => {
        let counter = 0;

        dbTestFixturesUtils.loadReleaseSource(database, { provider });
        dbTestFixturesUtils.loadReleaseSource(database, { provider });
        dbTestFixturesUtils.loadReleaseSource(database, { provider });
        dbTestFixturesUtils.loadReleaseSource(database, { provider });

        const fetchReleasesFromSource = spy(() => {
          counter += 1;

          if (counter === 1) {
            return Promise.resolve([]);
          }
          if (counter === 2) {
            throw new Error("foo");
          }
          if (counter === 3) {
            return Promise.resolve([
              dbTestFixturesUtils.generateRelease({ provider }),
            ]);
          }

          return Promise.resolve([]);
        });

        using timer = new FakeTime();
        let resolved = false;

        const p = new App({
          database,
          service: { fetchReleasesFromSource },
          provider,
          signal: new AbortController().signal,
        }).execute().then(() => {
          resolved = true;
        });

        assertEquals(resolved, false);
        await timer.nextAsync();
        assertEquals(resolved, false);
        await timer.nextAsync();
        assertEquals(resolved, false);
        await timer.nextAsync();
        assertEquals(resolved, false);
        await timer.nextAsync();
        assertEquals(resolved, true);
        await p;
      });
    });
  }
});
