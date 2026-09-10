import test from "node:test";
import assert from "node:assert/strict";
import { editorSnapshot, readingStats } from "../static/js/experience.mjs";

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

test("editor change tracking ignores topic order but preserves authored text", () => {
  const original = {
    title: "A good question",
    content: "A considered perspective.",
    topics: ["2", "1"],
  };
  assert.equal(
    editorSnapshot(original),
    editorSnapshot({ ...original, topics: ["1", "2"] }),
  );
  assert.notEqual(
    editorSnapshot(original),
    editorSnapshot({ ...original, title: "A better question" }),
  );
  assert.notEqual(
    editorSnapshot(original),
    editorSnapshot({ ...original, content: original.content + "\n" }),
  );
  assert.notEqual(
    editorSnapshot(original),
    editorSnapshot({ ...original, topics: ["1"] }),
  );
  assert.deepEqual(original.topics, ["2", "1"]);
});
