# Continue work on another machine

Cursor chat history does **not** sync across devices. Source of truth is **Git** + `docs/`.

Documentation map: [`README.md`](README.md) (this folder).

## Before you leave a machine

```bash
cd /path/to/Sailorport
git status
# update docs/PROGRESS.md status if you finished a step
git add -p   # or git add <files>
git commit -m "…"
git push
```

## On the next machine

```bash
git clone git@github.com:ekkywi/Sailorport.git   # first time
# or: git pull
```

1. Open the repo in Cursor.
2. Read [`PROGRESS.md`](PROGRESS.md) (status block at the top).
3. Open a new chat and paste [`RESUME-PROMPT.md`](RESUME-PROMPT.md).
4. Follow “Step berikutnya” / backlog links from PROGRESS → ROADMAP.

## First-time GitHub remote (once)

```bash
git branch -M main
git remote add origin git@github.com:ekkywi/Sailorport.git
git push -u origin main
```

## After each completed step

1. Update the status block + checklist in `PROGRESS.md`.
2. Update ROADMAP/QC when a backlog or debt item closes.
3. Commit and push.

Do not rely on chat memory alone.
