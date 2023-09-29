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
        nord0: "#2e3440",
        nord1: "#3b4252",
        nord2: "#434c5e",
        nord3: "#4c566a",
        nord4: "#d8dee9",
        nord5: "#e5e9f0",
        nord6: "#eceff4",
        nord7: "#8fbcbb",
        nord8: "#88c0d0",
        nord9: "#81a1c1",
        nord10: "#5e81ac",
        nord11: "#bf616a",
        nord12: "#d08770",
        nord13: "#ebcb8b",
        nord14: "#a3be8c",
        nord15: "#b48ead",
      },
      fontFamily: {
        fugaz: "'Fugaz One', cursive",
      },
      boxShadow: {
        project_card: `
        rgba(0, 0, 0, 0.2) 0px 12px 28px 0px, 
        rgba(0, 0, 0, 0.1) 0px 2px 4px 0px,
        rgba(255, 255, 255, 0.05) 0px 0px 0px 1px inset;
        `,
      },
    },
  },
  plugins: [],
};
