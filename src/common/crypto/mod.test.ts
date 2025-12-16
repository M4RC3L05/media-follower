import { beforeAll, describe, it } from "@std/testing/bdd";
import { assertEquals, assertRejects } from "@std/assert";
import { encode } from "@std/msgpack";
import { encodeBase64 } from "@std/encoding";
import { extractParts, hash, verify } from "../crypto/mod.ts";
import { config, initConfig } from "../config/mod.ts";

beforeAll(() => {
  initConfig();
});

describe("Crypto", () => {
  describe("hash()", () => {
    it("should hash a password", async () => {
      const p = await hash("foobar");
      const [hf, i, s, h] = extractParts(p);

      assertEquals(hf, "SHA-256");
      assertEquals(i, config().crypto.pbkdf2.iterations);
      assertEquals(s.length, config().crypto.pbkdf2.saltLength);
      assertEquals(h.length * 8, 256);
    });
  });

  describe("verify()", () => {
    it("should throw and error if malformed hash provided", async () => {
      const err = await assertRejects(() => verify("foobar", "foobar"));
      assertEquals(
        (err as Error).message,
        "Malformed hash provided",
      );

      const err2 = await assertRejects(() =>
        verify("foobar", encodeBase64(encode(["foobar"])))
      );
      assertEquals(
        (err2 as Error).message,
        "Malformed hash provided",
      );
    });

    it("should return true if passwords match", async () => {
      const p = await hash("foobar");
      const ok = await verify("foobar", p);

      assertEquals(ok, true);
    });

    it("should return false if passwords do not match", async () => {
      const p = await hash("foobar");
      const ok = await verify("foobar!", p);

      assertEquals(ok, false);
    });

    describe("when hash has different defaults", () => {
      it("should return true if passwords match", async () => {
        // Password = foo, Iterations = 2 & hash function algo = SHA-512 & salt length = 2
        const p =
          "lKdTSEEtNTEyAsQCPKDEQJi3+3SNymkcC0+P5GzRFDN3jw7Mel/SDh+xdtQlNhjhVhPr3GDsTCBsn+QA9vWcBUCPhEjVFmCFqV9qPvm3Oc0=";
        const ok = await verify("foo", p);

        assertEquals(ok, true);
      });

      it("should return false if passwords do not match", async () => {
        // Password = foo, Iterations = 2 & hash function algo = SHA-512 & salt length = 2
        const p =
          "lKdTSEEtNTEyAsQCPKDEQJi3+3SNymkcC0+P5GzRFDN3jw7Mel/SDh+xdtQlNhjhVhPr3GDsTCBsn+QA9vWcBUCPhEjVFmCFqV9qPvm3Oc0=";
        const ok = await verify("foo!", p);

        assertEquals(ok, false);
      });
    });
  });
});
