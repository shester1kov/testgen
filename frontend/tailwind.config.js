/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        // Graphite backgrounds
        'dark-900': '#0a0a0f',  // Deepest background
        'dark-800': '#12121a',  // Main background
        'dark-700': '#1a1a24',  // Card background
        'dark-600': '#24242e',  // Hover background
        'dark-500': '#2e2e38',  // Border

        // Neon Orange accents
        'neon-orange': '#ff6b35',       // Main accent
        'neon-orange-light': '#ff8555',
        'neon-orange-dark': '#e55a2b',
        'neon-glow': '#ff6b35',         // Glow effect

        // Secondary accents
        'cyber-blue': '#00d9ff',    // Secondary highlight
        'cyber-purple': '#b84fff',  // Tertiary
        'cyber-pink': '#ff2e97',    // Alert/Error

        // Text colors
        'text-primary': '#f1f5f9',   // Main text
        'text-secondary': '#94a3b8', // Secondary text
        'text-muted': '#64748b',     // Muted text
      },
      boxShadow: {
        'neon': '0 0 5px #ff6b35, 0 0 20px #ff6b35',
        'neon-sm': '0 0 3px #ff6b35, 0 0 10px #ff6b35',
        'neon-lg': '0 0 10px #ff6b35, 0 0 40px #ff6b35',
        'cyber-blue': '0 0 5px #00d9ff, 0 0 15px #00d9ff',
      },
      animation: {
        'glow': 'glow 2s ease-in-out infinite alternate',
        'pulse-slow': 'pulse 3s cubic-bezier(0.4, 0, 0.6, 1) infinite',
      },
      keyframes: {
        glow: {
          '0%': { boxShadow: '0 0 5px #ff6b35, 0 0 20px #ff6b35' },
          '100%': { boxShadow: '0 0 10px #ff6b35, 0 0 30px #ff6b35, 0 0 40px #ff6b35' },
        }
      },
      backgroundImage: {
        'gradient-radial': 'radial-gradient(var(--tw-gradient-stops))',
        'gradient-cyber': 'linear-gradient(135deg, #ff6b35 0%, #ff2e97 100%)',
      }
    },
  },
  plugins: [],
}
