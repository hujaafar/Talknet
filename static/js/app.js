/* Progressive enhancement: navigation, forms, and content work without JavaScript. */
(() => {
  "use strict";
  const root = document.documentElement;
  const reduced = matchMedia("(prefers-reduced-motion: reduce)");
  const motionButton = document.querySelector(".motion-toggle");
  let paused = false;
  try {
    paused = localStorage.getItem("talknet.motion") === "paused";
  } catch (_) {
    /* Storage can be disabled. */
  }
  let scheduled = false;
  const hero = document.querySelector(".hero");
  const progress = document.querySelector(".reading-progress");
  const update = () => {
    scheduled = false;
    if (paused || reduced.matches) return;
    const distance = Math.max(1, root.scrollHeight - innerHeight);
    progress.style.transform = `scaleX(${Math.min(1, Math.max(0, scrollY / distance))})`;
    if (hero && innerWidth > 640 && scrollY < hero.offsetHeight + 100) {
      hero.style.setProperty(
        "--hero-shift",
        `${Math.min(scrollY * 0.12, 65)}px`,
      );
      hero.style.setProperty(
        "--hero-scale",
        `${1 + Math.min(scrollY / 10000, 0.045)}`,
      );
    }
  };
  const syncMotion = () => {
    const disabled = paused || reduced.matches;
    root.dataset.motion = disabled ? "paused" : "running";
    motionButton.hidden = false;
    motionButton.textContent = reduced.matches
      ? "Reduced motion enabled"
      : paused
        ? "Enable motion"
        : "Pause motion";
    motionButton.disabled = reduced.matches;
    motionButton.setAttribute("aria-pressed", String(disabled));
    update();
  };
  motionButton.addEventListener("click", () => {
    paused = !paused;
    try {
      localStorage.setItem("talknet.motion", paused ? "paused" : "running");
    } catch (_) {}
    syncMotion();
  });
  reduced.addEventListener("change", syncMotion);
  addEventListener(
    "scroll",
    () => {
      if (!scheduled && !paused && !reduced.matches) {
        scheduled = true;
        requestAnimationFrame(update);
      }
    },
    { passive: true },
  );
  addEventListener("resize", update, { passive: true });
  syncMotion();
  if ("IntersectionObserver" in window) {
    // Reveal by animating visible elements, never by hiding essential content.
    const observer = new IntersectionObserver(
      (entries) =>
        entries.forEach((entry) => {
          if (!entry.isIntersecting) return;
          if (!paused && !reduced.matches)
            entry.target.animate(
              [
                { opacity: 0.5, transform: "translateY(18px)" },
                { opacity: 1, transform: "translateY(0)" },
              ],
              { duration: 550, easing: "cubic-bezier(.2,.65,.3,1)" },
            );
          observer.unobserve(entry.target);
        }),
      { threshold: 0.1 },
    );
    document
      .querySelectorAll("[data-reveal]")
      .forEach((el) => observer.observe(el));
  }
  let toastTimer;
  const toast = (message) => {
    const box = document.querySelector(".toast");
    clearTimeout(toastTimer);
    box.textContent = message;
    box.hidden = false;
    toastTimer = setTimeout(() => {
      box.hidden = true;
    }, 5000);
  };
  document.querySelectorAll("[data-reaction-id]").forEach((group) => {
    group.addEventListener("click", async (event) => {
      const button = event.target.closest("[data-vote]");
      if (!button || group.dataset.busy === "true") return;
      group.dataset.busy = "true";
      const buttons = [...group.querySelectorAll("button")];
      buttons.forEach((item) => {
        item.disabled = true;
      });
      try {
        const response = await fetch("/like_dislike", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            "X-CSRF-Token": document.querySelector("meta[name=csrf-token]")
              .content,
          },
          body: JSON.stringify({
            postId: Number(group.dataset.reactionId),
            type: group.dataset.reactionType,
            action: button.dataset.vote,
          }),
        });
        if (response.status === 401) {
          location.assign("/login");
          return;
        }
        if (!response.ok) {
          toast(
            response.status === 403
              ? "Your session changed. Refresh and try again."
              : "Couldn’t save your reaction. Try again.",
          );
          return;
        }
        const data = await response.json();
        group.querySelector("[data-count=like]").textContent = data.likeCount;
        group.querySelector("[data-count=dislike]").textContent =
          data.dislikeCount;
        group
          .querySelector("[data-vote=like]")
          .setAttribute("aria-pressed", String(data.reaction === 1));
        group
          .querySelector("[data-vote=dislike]")
          .setAttribute("aria-pressed", String(data.reaction === 0));
        toast(data.reaction === -1 ? "Reaction removed." : "Reaction saved.");
      } catch (_) {
        toast("Connection lost. Your reaction wasn’t saved.");
      } finally {
        group.dataset.busy = "false";
        buttons.forEach((item) => {
          item.disabled = false;
        });
      }
    });
  });
  document.querySelectorAll("[data-counter]").forEach((field) => {
    const updateCount = () => {
      document.getElementById(field.dataset.counter).textContent =
        `${[...field.value].length.toLocaleString()} / ${field.maxLength.toLocaleString()}`;
    };
    field.addEventListener("input", updateCount);
    updateCount();
  });
  document.querySelectorAll("[data-password]").forEach((button) =>
    button.addEventListener("click", () => {
      const input = document.getElementById(button.dataset.password);
      const show = input.type === "password";
      input.type = show ? "text" : "password";
      button.textContent = show ? "Hide" : "Show";
      button.setAttribute(
        "aria-label",
        show ? "Hide password" : "Show password",
      );
    }),
  );
  document.querySelectorAll(".share-button").forEach((button) =>
    button.addEventListener("click", async () => {
      try {
        await navigator.clipboard.writeText(location.href.split("#")[0]);
        toast("Discussion link copied.");
      } catch (_) {
        toast("Copy the discussion link from your address bar.");
      }
    }),
  );
  const topicInputs = [...document.querySelectorAll(".topic-picker input")];
  topicInputs.forEach((input) =>
    input.addEventListener("change", () => {
      if (topicInputs.filter((item) => item.checked).length > 3) {
        input.checked = false;
        toast("Choose up to three topics.");
      }
    }),
  );
  document.querySelectorAll("form[method=post]").forEach((form) =>
    form.addEventListener("submit", () => {
      const button = form.querySelector("button[type=submit]");
      if (button) {
        button.disabled = true;
        button.setAttribute("aria-busy", "true");
      }
    }),
  );
  addEventListener("pageshow", () =>
    document
      .querySelectorAll("form button[aria-busy=true]")
      .forEach((button) => {
        button.disabled = false;
        button.removeAttribute("aria-busy");
      }),
  );
})();
