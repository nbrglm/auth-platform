function setTheme() {
  let theme = window.matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "light"
  document.documentElement.setAttribute("data-theme", theme);
  Alpine.store('darkMode', theme === "dark")
}

window.matchMedia("(prefers-color-scheme: dark)").addEventListener("change", setTheme);

document.addEventListener("alpine:init", () => {
  setTheme();
})