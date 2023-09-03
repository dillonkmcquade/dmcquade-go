build:
	esbuild --bundle --minify web/static/main.css --outfile=web/static/styles.css --loader:.jpg=dataurl

