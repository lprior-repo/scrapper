/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './packages/webapp/index.html',
    './packages/webapp/src/**/*.{js,ts,jsx,tsx}',
  ],
  theme: {
    extend: {
      colors: {
        'dark-bg': '#0d1117',
        'dark-border': '#30363d',
        'dark-text': '#c9d1d9',
        'dark-text-secondary': '#8b949e',
        'accent-green': '#238636',
        'accent-green-hover': '#2ea043',
        'accent-blue': '#1f6feb',
        'accent-red': '#f85149',
      },
    },
  },
  plugins: [],
}