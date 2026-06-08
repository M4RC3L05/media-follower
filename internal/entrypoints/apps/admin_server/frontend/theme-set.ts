const setTheme = () => {
  document.documentElement.setAttribute('data-bs-theme', (window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'))
}

window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', () => {
  setTheme()
})

window.addEventListener('DOMContentLoaded', () => {
  setTheme()
})
