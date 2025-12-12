import {
  afterAll,
  afterEach,
  beforeEach,
  describe,
  it,
} from "@std/testing/bdd";
import { FakeTime } from "@std/testing/time";
import { assertSpyCall, assertSpyCalls, spy } from "@std/testing/mock";
import { CustomDatabase } from "#src/common/database/mod.ts";
import { config, initConfig } from "#src/common/config/mod.ts";
import { dbTestFixturesUtils, dbTestUtils } from "#src/common/utils/mod.ts";
import { App } from "./app.ts";
import { ReleaseSourceProvider } from "#src/common/database/enums/release-source-provider.ts";
import { makeLogger } from "#src/common/logger/mod.ts";
import { assertEquals } from "@std/assert/equals";
import type { DbReleaseSourcesTable } from "#src/common/database/types.ts";

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

  for (
    const provider of Object.values(ReleaseSourceProvider)
  ) {
    describe(`Provider: ${provider}`, () => {
      it("should do nothing if signal is already aborted", async () => {
        const abortController = new AbortController();
        abortController.abort();

        const fetchReleaseSource = spy(() => Promise.resolve(undefined));

        await new App({
          database,
          service: {
            fetchReleaseSource: fetchReleaseSource,
          },
          provider,
          signal: abortController.signal,
        }).execute();

        assertSpyCalls(fetchReleaseSource, 0);
      });

      it("should do nothing no sources available", async () => {
        const fetchReleaseSource = spy(() => Promise.resolve(undefined));

        await new App({
          database,
          service: {
            fetchReleaseSource: fetchReleaseSource,
          },
          provider,
          signal: new AbortController().signal,
        }).execute();

        assertSpyCalls(fetchReleaseSource, 0);
      });

      it("should log error if an error occures", async () => {
        const error = new Error("foo");
        const source = dbTestFixturesUtils.loadReleaseSource(database, {
          provider,
        });
        const fetchReleaseSource = spy(() => {
          throw error;
        });

        using errorLogStub = spy(
          makeLogger("sync-release-sources-app"),
          "error",
        );

        await new App({
          database,
          service: { fetchReleaseSource },
          provider,
          signal: new AbortController().signal,
        }).execute();

        assertSpyCalls(fetchReleaseSource, 1);
        assertSpyCalls(errorLogStub, 1);
        assertSpyCall(fetchReleaseSource, 0, { args: [source] });
        assertSpyCall(errorLogStub, 0, {
          args: ["Could not sync release source successfully", {
            source: { id: source.id },
            error,
          }],
        });
      });

      it("should do nothing if `fetchReleaseSource()` returns nothing", async () => {
        const fetchReleaseSource = spy(() => Promise.resolve(undefined));
        const source = dbTestFixturesUtils.loadReleaseSource(database, {
          provider,
        });

        await new App({
          database,
          service: { fetchReleaseSource },
          provider,
          signal: new AbortController().signal,
        }).execute();

        assertSpyCalls(fetchReleaseSource, 1);
        assertSpyCall(fetchReleaseSource, 0, { args: [source] });
      });

      it("should update release source", async () => {
        const source = dbTestFixturesUtils.loadReleaseSource(database, {
          provider,
          raw: JSON.stringify(
            dbTestFixturesUtils.generateReleaseSourceItunesArtist({
              artistId: 1,
              artistName: "foo",
            }),
          ),
        });
        const updated = dbTestFixturesUtils.generateReleaseSourceItunesArtist({
          artistId: 1,
          artistName: "bar",
        });
        const fetchReleaseSource = spy(() =>
          Promise.resolve({ ...source, raw: JSON.stringify(updated) })
        );

        await new App({
          database,
          service: { fetchReleaseSource },
          provider,
          signal: new AbortController().signal,
        }).execute();

        assertSpyCalls(fetchReleaseSource, 1);
        assertSpyCall(fetchReleaseSource, 0, { args: [source] });
        assertEquals(
          database.sql<
            DbReleaseSourcesTable
          >`select id, provider, json(raw) as raw from release_sources where id = ${source.id}`[
            0
          ]!,
          {
            id: source.id,
            provider: source.provider,
            raw: JSON.stringify(updated),
          },
        );
      });

      it("should delay if more items exists", async () => {
        let counter = 0;

        dbTestFixturesUtils.loadReleaseSource(database, {
          provider,
        });
        dbTestFixturesUtils.loadReleaseSource(database, {
          provider,
        });
        const three = dbTestFixturesUtils.loadReleaseSource(database, {
          provider,
        });
        dbTestFixturesUtils.loadReleaseSource(database, {
          provider,
        });

        const fetchReleaseSource = spy(() => {
          counter += 1;

          if (counter === 1) return Promise.resolve(undefined);
          if (counter === 2) throw new Error("foo");
          if (counter === 3) {
            return Promise.resolve(
              {
                ...three,
                raw: JSON.stringify(
                  dbTestFixturesUtils.generateReleaseSourceItunesArtist(),
                ),
              },
            );
          }

          return Promise.resolve(undefined);
        });

        using timer = new FakeTime();
        let resolved = false;

        const p = new App({
          database,
          service: { fetchReleaseSource },
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
