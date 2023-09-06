build:
	esbuild --bundle --minify web/css/main.css --outfile=web/static/styles.css --loader:.jpg=dataurl
	go build -o ./bin/main ./...

