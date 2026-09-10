/* Run before styles to restore a visitor's theme without a bright first frame. */
try {
  const theme = localStorage.getItem("talknet.theme");
  if (theme === "light" || theme === "dark")
    document.documentElement.dataset.theme = theme;
} catch (_) {
  /* The default dark theme also works without storage. */
}
