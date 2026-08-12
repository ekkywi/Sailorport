# Data directory

| Path | Purpose |
|------|---------|
| `workspaces/` | Output scaffold (service folders). Gitignored contents. |

## Development (`go run` API on host)

The API creates `workspaces/` as the user running the process. If you previously
ran full Compose with a bind mount, Docker may have left this folder owned by
`root`/`nobody` — then scaffold fails with `permission denied`.

Fix once:

```bash
sudo chown -R "$USER:$USER" data
mkdir -p data/workspaces
```

Daily dev usually only starts Postgres in Compose (no workspace mount), so the
API owns the folder itself — no recurring `chown`.

## Self-host Compose

Workspaces use the **named volume** `sailorport_workspaces` (not a bind mount).
No host `chown` required for scaffold inside the API container.

Note: agent-on-host deploy works best with the API also on the host (dev mode),
so `workspace_path` in the DB is a host filesystem path the agent can `chdir` to.
