import test from "node:test";
import assert from "node:assert/strict";
import { readingStats } from "../static/js/experience.mjs";

test("an empty editor has no words and a minimum one-minute estimate", () => {
  assert.deepEqual(readingStats(" \n\t "), { words: 0, minutes: 1 });
});
test("reading statistics handle line breaks, repeated whitespace, and estimates", () => {
  assert.deepEqual(readingStats("A thoughtful\n\n reply\t here"), {
    words: 4,
    minutes: 1,
  });
  assert.deepEqual(readingStats("word ".repeat(201)), {
    words: 201,
    minutes: 2,
  });
});
