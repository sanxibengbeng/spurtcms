// tailwind.js
window.tailwind = {
  config: {
    theme: {
      extend: {
        keyframes: {
          dropanime: {
            '0%': {
              opacity: '0',
              transform: 'translateY(-10px) scale(0.9)',
            },
            '70%': {
              opacity: '1',
              transform: 'translateY(5px)',
            },
            '100%': {
              transform: 'translateY(0)',
            },
          },
        },
        animation: {
          dropanime: 'dropanime 0.1s ease-out',
        },
        scrollbar: {
          thin: {
            'scrollbar-width': 'thin',
          }
        }
      }
    },
    plugins: [
      function ({ addUtilities }) {
        addUtilities(
          {
            '.scrollbar-thin': {
              'scrollbar-width': 'thin',
            }
          },
          ['responsive']
        );
      },
    ],
  }
};
