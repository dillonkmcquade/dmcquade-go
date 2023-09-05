# dmcquade.dev

This project is a conversion of my previous React-based portfolio site to a server side rendered site.

- HTML is served to the client with a Go HTTP server
- Data is injected to the HTML on the server with the html/template package included in the Go standard library
- There are only 15 lines of JavaScript for the typewriter effect
- HTMX is used for making ajax requests to other pages, reducing the need to reinterpret css and JS, making page loads fast
