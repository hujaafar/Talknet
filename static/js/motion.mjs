// A saved choice overrides the system setting only after the visitor opts in.
export function motionEnabled(preference, prefersReducedMotion) {
  if (preference === "running") return true;
  if (preference === "paused") return false;
  return !prefersReducedMotion;
}

const clamp = (value, min, max) => Math.min(max, Math.max(min, value));

// Scroll transforms are bounded so overscroll and short pages cannot distort the UI.
export function scrollFrame(
  scrollY,
  heroHeight,
  documentHeight,
  viewportHeight,
) {
  const y = Number.isFinite(scrollY) ? Math.max(0, scrollY) : 0;
  const progress = clamp(
    y / Math.max(1, documentHeight - viewportHeight),
    0,
    1,
  );
  const hero = clamp(y / Math.max(1, heroHeight), 0, 1);
  return {
    progress,
    artShift: hero * 95,
    artScale: 1.04 + hero * 0.12,
    artRotate: hero * -3,
    copyShift: hero * -42,
    noteShift: hero * -30,
  };
}
