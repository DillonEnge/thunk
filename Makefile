templ:
	go run github.com/a-h/templ/cmd/templ@latest generate -watch -proxy="http://localhost:$(PORT)" -open-browser=false ;

templ_no_watch:
	go run github.com/a-h/templ/cmd/templ@latest generate;

tailwind:
	npx tailwindcss -i templates/static/input.css -o templates/static/output.css --watch=always --minify

tailwind_no_watch:
	npx tailwindcss -i templates/static/input.css -o templates/static/output.css --minify

sync_static:
	go run github.com/air-verse/air@latest \
	--build.cmd "templ generate --notify-proxy" \
	--build.bin "true" \
	--build.delay "100" \
	--build.exclude_dir "" \
	--build.include_dir "templates/static" \
	--build.include_ext "css" ;

air:
	go run github.com/air-verse/air@latest ;

dev:
	make -j3 templ tailwind air

test:
	go test ./... ;

lint:
	golangci-lint run

regen: templ_no_watch tailwind_no_watch
