// @ts-ignore
import "https://unpkg.com/bootstrap@5.3.8/dist/js/bootstrap.bundle.js";

const getPreferredTheme = () => {
  return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
}

const setTheme = () => {
  document.documentElement.setAttribute('data-bs-theme', (window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'))
}

window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', () => {
  setTheme()
})

window.addEventListener('DOMContentLoaded', () => {
  setTheme()
})
