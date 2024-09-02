function switchMode() {
  const html = document.getElementsByTagName("html")[0]
  const theme = html.getAttribute("theme")
  html.setAttribute("theme", theme !== "dark" ? "dark" : "light")
}
