import test from "node:test";
import assert from "node:assert/strict";
import { motionEnabled, scrollFrame } from "../static/js/motion.mjs";

test("system reduced motion is honored until the visitor explicitly opts in", () => {
  assert.equal(motionEnabled("auto", true), false);
  assert.equal(motionEnabled("auto", false), true);
  assert.equal(motionEnabled("running", true), true);
  assert.equal(motionEnabled("paused", false), false);
  assert.equal(motionEnabled(null, true), false);
});

test("overscroll never creates negative progress or unbounded transforms", () => {
  const before = scrollFrame(-400, 700, 2500, 900);
  assert.equal(before.progress, 0);
  assert.equal(before.artShift, 0);
  const after = scrollFrame(1000000, 700, 2500, 900);
  assert.equal(after.progress, 1);
  assert.equal(after.artShift, 95);
  assert.ok(Math.abs(after.artScale - 1.16) < 1e-12);
  assert.equal(after.copyShift, -42);
});

test("a document shorter than the viewport and a missing hero remain finite", () => {
  const frame = scrollFrame(0, 0, 300, 900);
  assert.equal(frame.progress, 0);
  assert.ok(Object.values(frame).every(Number.isFinite));
  assert.ok(Object.values(scrollFrame(NaN, 0, 0, 0)).every(Number.isFinite));
});

test("art, title, and note move at different depths as the page scrolls", () => {
  const start = scrollFrame(0, 700, 2100, 700);
  const halfway = scrollFrame(350, 700, 2100, 700);
  assert.equal(halfway.progress, 0.25);
  assert.ok(halfway.artShift > start.artShift);
  assert.ok(halfway.artScale > start.artScale);
  assert.ok(halfway.copyShift < start.copyShift);
  assert.notEqual(halfway.noteShift, halfway.copyShift);
});
