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

## SSH site (`ssh evenbom.se`)

`docker-compose.yml` runs two services side by side:

- `portfolio` — the web app, on port 5000 (put a reverse proxy in front for 80/443).
- `ssh-site` — the [charmbracelet/wish](https://github.com/charmbracelet/wish) TUI portfolio at `ssh-site/`, mapped to host port 22 so visitors can just run `ssh evenbom.se`.

Only one process can bind port 22 on the host, so before deploying this, move the
real admin SSH daemon to another port on the production server:

1. Edit `/etc/ssh/sshd_config` and change `Port 22` to e.g. `Port 2222`.
2. Validate the config: `sudo sshd -t`.
3. Open the new port in the firewall / cloud security group (keep 22 open too — it'll be used by `ssh-site`).
4. Restart sshd (`sudo systemctl restart sshd`), then **in a new terminal, without closing your current session**, confirm `ssh -p 2222 <user>@evenbom.se` still works before disconnecting.
5. `docker compose up -d --build` to bring up `ssh-site` bound to port 22.
6. Confirm `ssh evenbom.se` opens the TUI, and admin access is now `ssh -p 2222 <user>@evenbom.se`.

`ssh-site`'s host key and `access.log` persist in the `./ssh-site-data` volume
(gitignored). Its listen address, host key path, and access log path are
configurable via the `SSH_SITE_ADDR`, `SSH_SITE_HOST_KEY_PATH`, and
`SSH_SITE_ACCESS_LOG` env vars (see `ssh-site/Dockerfile` for the defaults
used in Docker).
