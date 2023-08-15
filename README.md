# gin + go-assets-builder starter

`html` is output (vite `outDir`) of `npm run build` of https://github.com/wildfluss/vite-singlefile-vue (another boilerplate for building vue-ts template to single HTML file )

```bash
air
```

and go to http://localhost:8080/

## Build single file 

```bash
go install github.com/jessevdk/go-assets-builder@latest
```

then

```bash
go-assets-builder html -o assets.go && \
go build -o assets-in-binary
```

## Setup

Get go and air 

```bash
curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
```


