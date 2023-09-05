build:
	esbuild --bundle --minify web/static/css/main.css --outfile=web/static/css/styles.css --loader:.jpg=dataurl
	go build -o ./bin/main ./cmd/main/main.go

