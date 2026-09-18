# New Portfolio In Golang
Making a new portfolio website with:
* Go backend
* A paginated /projects JSON endpoint backed by the GitHub API (15 repos per page)
* Svelte frontend

## Local development

The Svelte frontend compiles to `static/js/bundle.js` / `bundle.css`, which are
build output (gitignored) rather than source — `docker build` generates them
in its own stage. To run the server directly on the host, build them once
after pulling frontend changes:

```sh
cd svelte-components && npm install && npm run build
```

Then, from the repo root:

```sh
go run .
```
