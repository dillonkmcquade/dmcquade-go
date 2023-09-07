/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./template/*.html"],
  theme: {
    extend: {
      height: {
        hero: "90vh",
        hero_mobile: "calc(100vh + 300px)",
      },
      backgroundImage: {
        hero: "url('/static/images/banner.jpg')",
      },
      backgroundColor: {
        head: "rgba(0, 0, 0, 0.4)",
        header: "rgba(0, 0, 0, 0.7)",
      },
      textColor: {
        nord4: "#d8dee9",
        nord7: "#8fbcbb",
        nord8: "#88c0d0",
      },
      fontFamily: {
        fugaz: "'Fugaz One', cursive",
      },
    },
  },
  plugins: [],
};
