import type { Config } from "tailwindcss";
const {heroui} = require("@heroui/theme");
export default {
  content: [
    "./pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./node_modules/@heroui/theme/dist/**/*.{js,ts,jsx,tsx}",
    "./components/**/*.{js,ts,jsx,tsx,mdx}",
    "./app/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      colors: {
        background: "var(--background-color)",
        foreground: "#1E1E1E",
        text:"#E0E0E0",
        textsecondary:"#B0B0B0",
        'muted-purple': 'rgba(44, 129, 163, 0.3)',
        'muted-purple-hover': 'rgba(44, 129, 163)',
      },
    },
  },
  darkMode: "class",
  plugins: [
    heroui({
      defaultTheme: 'dark',
      themes: {
        dark: {
          colors : {
            default: {
              DEFAULT: "#006FEE",
              50: "#18181b",
              100: "#1E1E1E",
              200: "#2E2E2E",
              300: "#52525b",
              400: "#71717a",
              500: "#B0B0B0",
              600: "#E0E0E0",
              700: "#e4e4e7",
              800: "#f4f4f5",
              900: "#fafafa"
            },
          }
        },
      },
    }),
  ],
} satisfies Config;
