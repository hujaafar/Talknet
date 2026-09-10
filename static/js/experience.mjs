export function readingStats(text) {
  const words = text.trim().split(/\s+/u).filter(Boolean).length;
  return { words, minutes: Math.max(1, Math.ceil(words / 200)) };
}

export function initExperience(toast, motionRunning) {
  const root = document.documentElement;
  const themeButton = document.querySelector(".theme-toggle");
  const syncTheme = () =>
    themeButton.setAttribute(
      "aria-label",
      `Switch to ${root.dataset.theme === "dark" ? "light" : "dark"} theme`,
    );
  themeButton.hidden = false;
  syncTheme();
  themeButton.addEventListener("click", () => {
    root.dataset.theme = root.dataset.theme === "dark" ? "light" : "dark";
    try {
      localStorage.setItem("talknet.theme", root.dataset.theme);
    } catch (_) {}
    syncTheme();
  });

  const dialog = document.querySelector(".command-palette");
  const query = dialog.querySelector("input[type=search]");
  const results = dialog.querySelector(".palette-results");
  const status = dialog.querySelector(".palette-status");
  const trigger = document.querySelector(".palette-trigger");
  let timer,
    controller,
    generation = 0;
  const cancelSearch = () => {
    clearTimeout(timer);
    controller?.abort();
    generation++;
  };
  const open = () => {
    if (!dialog.open) dialog.showModal();
    query.focus();
  };
  trigger.hidden = false;
  trigger.addEventListener("click", open);
  dialog
    .querySelector(".palette-close")
    .addEventListener("click", () => dialog.close());
  dialog.addEventListener("close", cancelSearch);
  dialog.addEventListener("click", (event) => {
    if (event.target === dialog) {
      const box = dialog.getBoundingClientRect();
      if (
        event.clientX < box.left ||
        event.clientX > box.right ||
        event.clientY < box.top ||
        event.clientY > box.bottom
      )
        dialog.close();
    }
  });
  document.addEventListener("keydown", (event) => {
    if ((event.ctrlKey || event.metaKey) && event.key.toLowerCase() === "k") {
      event.preventDefault();
      open();
    }
  });
  query.addEventListener("input", () => {
    cancelSearch();
    results.replaceChildren();
    const term = query.value.trim();
    if ([...term].length < 2) {
      status.textContent = "Type at least two characters to search.";
      return;
    }
    status.textContent = "Searching the community…";
    const requestGeneration = generation;
    timer = setTimeout(async () => {
      controller = new AbortController();
      try {
        const response = await fetch(
          `/api/search?q=${encodeURIComponent(term)}`,
          { signal: controller.signal },
        );
        if (!response.ok) throw new Error("Search failed");
        const data = await response.json();
        if (requestGeneration !== generation || !dialog.open) return;
        results.replaceChildren();
        for (const result of data.results) {
          const item = document.createElement("li");
          const link = document.createElement("a");
          link.href = `/post-details?post_id=${result.id}`;
          const title = document.createElement("strong");
          title.textContent = result.title;
          const meta = document.createElement("span");
          meta.textContent = `${result.topic} · ${result.author}`;
          link.append(title, meta);
          item.append(link);
          results.append(item);
        }
        status.textContent = data.results.length
          ? `${data.results.length} matching ${data.results.length === 1 ? "discussion" : "discussions"}`
          : "No matches yet. Try another word or explore all discussions.";
      } catch (error) {
        if (error.name !== "AbortError" && requestGeneration === generation)
          status.textContent =
            "Search is unavailable. Try again or use the search button.";
      }
    }, 180);
  });
  dialog.addEventListener("keydown", (event) => {
    // Search inputs consume Escape to clear their text in some browsers.
    if (event.key === "Escape") {
      event.preventDefault();
      dialog.close();
      return;
    }
    if (event.key !== "ArrowDown" && event.key !== "ArrowUp") return;
    const links = [...results.querySelectorAll("a")];
    if (
      !links.length ||
      (event.target !== query && !results.contains(event.target))
    )
      return;
    event.preventDefault();
    const index = links.indexOf(document.activeElement);
    const next =
      event.key === "ArrowDown"
        ? (index + 1) % links.length
        : index <= 0
          ? links.length - 1
          : index - 1;
    links[next].focus();
  });

  document.querySelectorAll("[data-bookmark]").forEach((button) => {
    button.hidden = false;
    button.addEventListener("click", async () => {
      const saved = button.getAttribute("aria-pressed") === "true";
      button.disabled = true;
      try {
        const response = await fetch("/bookmarks", {
          method: "POST",
          headers: {
            Accept: "application/json",
            "X-CSRF-Token": document.querySelector('meta[name="csrf-token"]')
              .content,
          },
          body: new URLSearchParams({
            post_id: button.dataset.bookmark,
            action: saved ? "remove" : "save",
          }),
        });
        if (response.status === 401) {
          location.assign("/login");
          return;
        }
        if (!response.ok) throw new Error("Save failed");
        const data = await response.json();
        document
          .querySelectorAll(`[data-bookmark="${button.dataset.bookmark}"]`)
          .forEach((item) =>
            item.setAttribute("aria-pressed", String(data.saved)),
          );
        // Re-fetch the private collection so counts, pagination, and empty states stay correct.
        if (
          !data.saved &&
          document.querySelector(
            '.profile-tabs a[aria-current="page"][href="/profile?tab=saved"]',
          )
        ) {
          location.reload();
          return;
        }
        toast(
          data.saved
            ? "Saved for later. Find it in your profile."
            : "Removed from your saved discussions.",
        );
      } catch (_) {
        toast("Couldn’t update your saved discussions. Try again.");
      } finally {
        button.disabled = false;
      }
    });
  });

  const viewControls = document.querySelector(".feed-view-controls");
  if (viewControls) {
    const feed = document.querySelector(".feed");
    const selectView = (value) => {
      feed.dataset.view = value;
      viewControls
        .querySelectorAll("button")
        .forEach((button) =>
          button.setAttribute(
            "aria-pressed",
            String(button.dataset.view === value),
          ),
        );
    };
    let view = "comfortable";
    try {
      if (localStorage.getItem("talknet.feed") === "compact") view = "compact";
    } catch (_) {}
    selectView(view);
    viewControls.hidden = false;
    viewControls.addEventListener("click", (event) => {
      const button = event.target.closest("[data-view]");
      if (!button) return;
      selectView(button.dataset.view);
      try {
        localStorage.setItem("talknet.feed", button.dataset.view);
      } catch (_) {}
    });
  }

  const editor = document.querySelector(".discussion-editor");
  if (editor) {
    const toolbar = document.querySelector(".editor-tools");
    toolbar.hidden = false;
    const title = editor.querySelector("#title"),
      content = editor.querySelector("#content");
    const preview = editor.querySelector(".editor-preview"),
      fields = editor.querySelector(".editor-fields");
    const stats = editor.querySelector(".editor-stats");
    stats.hidden = false;
    const updatePreview = () => {
      const metrics = readingStats(content.value);
      stats.textContent = `${metrics.words} ${metrics.words === 1 ? "word" : "words"} · ${metrics.minutes} min read`;
      preview.querySelector("h3").textContent =
        title.value.trim() || "Your headline goes here.";
      const body = preview.querySelector(".preview-content");
      body.replaceChildren();
      // Plain text is preserved. Never interpret authored text as markup.
      for (const paragraph of (
        content.value.trim() ||
        "Your perspective will appear here as you write."
      ).split(/\n\s*\n/u)) {
        const p = document.createElement("p");
        p.textContent = paragraph;
        body.append(p);
      }
      const topics = preview.querySelector(".preview-topics");
      topics.replaceChildren();
      editor
        .querySelectorAll(".topic-picker input:checked")
        .forEach((input) => {
          const span = document.createElement("span");
          span.textContent = input.nextElementSibling.textContent;
          topics.append(span);
        });
    };
    const showPreview = (enabled) => {
      updatePreview();
      preview.hidden = !enabled;
      fields.hidden = enabled;
      toolbar
        .querySelector(".editor-write")
        .setAttribute("aria-pressed", String(!enabled));
      toolbar
        .querySelector(".editor-preview-toggle")
        .setAttribute("aria-pressed", String(enabled));
    };
    toolbar
      .querySelector(".editor-write")
      .addEventListener("click", () => showPreview(false));
    toolbar
      .querySelector(".editor-preview-toggle")
      .addEventListener("click", () => showPreview(true));
    toolbar
      .querySelector(".focus-toggle")
      .addEventListener("click", (event) => {
        const container = document.querySelector(".compose-layout");
        const focused = container.dataset.focus !== "on";
        container.dataset.focus = focused ? "on" : "off";
        event.currentTarget.setAttribute("aria-pressed", String(focused));
        event.currentTarget.textContent = focused
          ? "Exit focus mode ↙"
          : "Focus mode ↗";
      });
    editor.addEventListener("input", updatePreview);
    editor.addEventListener("change", updatePreview);
    editor.addEventListener("invalid", () => showPreview(false), true);
    updatePreview();
  }

  const stage = document.querySelector(".hero-visual");
  const hero = document.querySelector(".hero");
  if (hero && matchMedia("(hover: hover) and (pointer: fine)").matches) {
    let frame = 0,
      x = 0,
      y = 0;
    hero.addEventListener("pointermove", (event) => {
      if (!motionRunning()) return;
      const rect = hero.getBoundingClientRect();
      x = Math.max(
        -1,
        Math.min(1, ((event.clientX - rect.left) / rect.width) * 2 - 1),
      );
      y = Math.max(
        -1,
        Math.min(1, ((event.clientY - rect.top) / rect.height) * 2 - 1),
      );
      if (!frame)
        frame = requestAnimationFrame(() => {
          stage.style.setProperty("--pointer-x", `${x * 14}px`);
          stage.style.setProperty("--pointer-y", `${y * 10}px`);
          frame = 0;
        });
    });
    hero.addEventListener("pointerleave", () => {
      if (frame) cancelAnimationFrame(frame);
      frame = 0;
      stage.style.setProperty("--pointer-x", "0px");
      stage.style.setProperty("--pointer-y", "0px");
    });
  }
}
