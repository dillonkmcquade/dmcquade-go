import Typewriter from "typewriter-effect/dist/core";

function showTypewriter() {
  const target = document.getElementById("typewriter");
  const typewriter = new Typewriter(target, {
    strings: [
      "Developer",
      "Husband",
      "Adventurer",
      "Life-long learner",
      "Lawn nerd",
    ],
    loop: true,
    autoStart: true,
  });
  typewriter.start();
}
showTypewriter();
