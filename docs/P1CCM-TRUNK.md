# P1 CCM Trunk — `feat/p1ccm-enbuild-helm-trunk`

This document explains the P1 CCM long-lived integration branch on this
repo (**`feat/p1ccm-enbuild-helm-trunk`**) from the chart-repo's
perspective: why it exists, how its publish flow is wired, and the hard
rules around not letting it collide with `main`.

If you cloned this repo and landed on this branch, **read this before
touching chart versions, CI workflows, or release tags.**

> File location note: this lives at `docs/P1CCM-TRUNK.md`, alongside
> the Hugo site config (`docs/config/`, `docs/content/`). Hugo only
> renders `docs/content/**`, so this file ships in the git tree but is
> deliberately **not** part of the published documentation site at
> <https://vivsoftorg.github.io/enbuild/>. It is repo-internal CI/trunk
> documentation.

---

## TL;DR

- This branch is the P1 CCM contract's iteration channel for the
  `enbuild` umbrella chart. It mirrors the long-lived feature-branch
  ("trunk") pattern already in use on `vivsoft-platform-ui` and
  `platform-one-eks` for the same contract.
- Chart versions on this branch carry a **`-p1ccm-trunk.N`** semver
  pre-release suffix (e.g. `0.0.49-p1ccm-trunk.0`). They never overwrite
  `main`'s stable line and helm will not auto-pick them as "latest".
- A dedicated GitHub Actions workflow,
  [`.github/workflows/helm-release-trunk.yml`](../.github/workflows/helm-release-trunk.yml),
  packages the chart on every chart-touching push and uploads the
  `.tgz` as a **GitHub Release** asset under the tag
  `enbuild-trunk-<version>`.
- It does **not** publish to the gh-pages helm repo
  (<https://vivsoftorg.github.io/enbuild>), does **not** run
  chart-releaser, does **not** mint `enbuild-<version>` tags.
- Consumers fetch the asset with `gh release download` and hand the
  local `.tgz` to the `roll-vendor13-ib.sh` wrapper in
  [`p1-cluster-mgmt`](https://gitlab.com/) via `LOCAL_CHART_TGZ=`.

---

## Why a new long-lived trunk

Until this branch existed, `enbuild-helm` did **not** follow the
single-trunk integration pattern that the rest of the P1 CCM contract
relies on. `main` is gated by team-approval MRs that slow per-week
iteration; in parallel, five feature branches were each rooted at chart
`0.0.45`:

- `feat/p1ccm-hub-grpc-service`
- `feat/p1ccm-keycloak-tenancy-chart`
- `feat/mongo-ha-v7`
- `feat/headlamp-whitelabel-codify`
- `fix/ccm-08b-installagent-config-idempotent`

Meanwhile, the actually-deployed branch
(`fix/ha-hub-connectivity`, chart `0.0.47`) carried production fixes
that the parallel branches did **not** know about. A merge of any
parallel branch as-is would silently regress shipped production state.

**Task #64** is the canonical example: chart `0.0.47` on the deployed
branch was missing the hub-side gRPC `Service` + `VirtualService`
templates. Those templates lived only on parallel branches that hadn't
been merged. The fleet of 4 agents dropped offline as soon as the
deployed chart rolled forward, because there was nowhere for the agents
to dial in. Trunk consolidation prevents this regression class.

The branch base is `fix/ha-hub-connectivity` (commit `5d0950d`,
chart `0.0.47`) which already carries:

- PR #47 (`b319cfd`) — CCM-08b
- `5d0950d` — `cosmosDB` literal-default stomp fix
- `ec1b752` — `mongo_endpoint_override`

**Wave 1 SHIPPED:** commit `c507a89` (closes #64) cherry-picked the
hub-side gRPC `Service` + `VirtualService` (originals `848420d` +
`258e932` — byte-identical across all 5 parallel branches per audit).
Chart bumped `0.0.47 → 0.0.48`, then `0.0.49-p1ccm-trunk.0` once the
pre-release suffix was applied for trunk publishing.

Subsequent waves consolidate the remaining parallel branches into this
trunk via **file-level cherry-pick**, never via `git merge` of the
source branches (which would re-introduce divergence at base
`0.0.45`).

---

## Trunk publish flow

**Workflow:** [`.github/workflows/helm-release-trunk.yml`](../.github/workflows/helm-release-trunk.yml)
(commit `29297bb`, closes #65, 2026-05-31)

**Triggers:**

```yaml
on:
  push:
    branches:
      - feat/p1ccm-enbuild-helm-trunk
    paths:
      - "charts/**"
  workflow_dispatch: {}
```

Pushes to other branches don't fire it. Non-chart pushes (docs, CI,
top-level `README`) don't fire it.

**What it does:**

1. Adds the same helm repos `helm-release.yml` adds (bitnami,
   jupyterhub, open-webui, loki-stack, headlamp) so subchart deps
   resolve.
2. `helm dependency build charts/enbuild`
3. `helm package charts/enbuild --destination /tmp` → produces
   `/tmp/enbuild-<version>.tgz` where `<version>` is read from
   `charts/enbuild/Chart.yaml`.
4. `gh release create enbuild-trunk-<version>` (or `gh release upload
   --clobber` if the tag already exists) with the `.tgz` as the release
   asset.

**Permissions:** `contents: write` only. No `gh-pages` writes. No
chart-releaser-action. No `CR_TOKEN`. The default `GITHUB_TOKEN` is
sufficient.

**Tag prefix `enbuild-trunk-<version>`** is deliberately distinct from
main's `enbuild-<version>` tag namespace — collision is structurally
impossible.

---

## Why GitHub Releases (not gh-pages, not chart-releaser, not OCI)

This was the central design question for #65. Options considered:

| Option | Why rejected |
| --- | --- |
| Add trunk to `helm-release.yml`'s trigger list | **Tried** (commit `510758b`) and **reverted** (commit `9861ccb`). `chart-releaser-action` auto-bumps the patch version on every CI run and pushes to gh-pages. Trunk's per-commit churn would compete with main's tag namespace + `index.yaml` and silently break main's release cadence. **This path is a hard "do not retry".** |
| Separate `gh-pages-trunk` branch | GitHub Pages serves exactly **one** branch per repo. A non-Pages-served branch can't act as a helm repo URL — useless. |
| Subdirectory on `gh-pages` (e.g. `gh-pages/trunk/`) | Requires custom workflow logic that `chart-releaser-action` does not natively support. Plumbing burden + still shares the `index.yaml`. |
| OCI (`ghcr.io`) | Works, but requires `packages: write`, helm OCI URL handling on the consumer, and a separate auth path. More moving parts than the contract needs. |
| **GitHub Releases (chosen)** | Simplest no-collision path: default `GITHUB_TOKEN`, distinct tag namespace, trivially consumed via `gh release download`. Zero interaction with the public helm repo. |

`artifacthub-repo.yml` continues to point at the
`main` + `gh-pages` distribution. Artifact Hub still sees only stable
main releases. Trunk artifacts are **not** advertised publicly.

---

## How to consume a trunk release

The roll wrapper in
[`p1-cluster-mgmt`](https://gitlab.com/)
(`scripts/roll-vendor13-ib.sh`) honors `LOCAL_CHART_TGZ=` and uses it
ahead of any chart-repo / OCI pull.

```bash
# Find the version you want
gh release list --repo vivsoftorg/enbuild | grep enbuild-trunk-

# Pull the asset
gh release download enbuild-trunk-0.0.49-p1ccm-trunk.0 \
  --repo vivsoftorg/enbuild \
  --pattern 'enbuild-*.tgz' \
  --dir /tmp

# Roll with the local tgz — Stage A Mongo overlay is mandatory until
# upstream `enbuild` chart flips `mongodb.enabled` default to false.
KUBECONFIG=/tmp/p1-via-bastion-18444.yaml \
  EXTRA_HELM_ARGS="--values $HOME/code/p1-cluster-mgmt/envs/vendor13-ib/stage-a-mongo-handoff.values.yaml" \
  LOCAL_CHART_TGZ=/tmp/enbuild-0.0.49-p1ccm-trunk.0.tgz \
  ~/code/p1-cluster-mgmt/scripts/roll-vendor13-ib.sh backend <current-or-new-be-sha>
```

The wrapper's chart-resolution order is:

1. `LOCAL_CHART_TGZ` (this trunk path — short-circuits everything else)
2. Chart repo (main's gh-pages — for stable rolls)
3. Legacy fallback

Stage A overlay note: since chart `0.0.48` the overlay carries
`enbuildBk.grpcVirtualService.host: enbuild-ib-vendor13.staging.dso.mil`
because the chart's default leaves that host blank — and the host
**must** match what the agent dials. Removing the overlay will break
the gRPC ingress added in #64.

---

## Hard rules for future agents

1. **Do not add `feat/p1ccm-enbuild-helm-trunk` to the trigger list of
   `.github/workflows/helm-release.yml`.** That re-introduces the
   chart-releaser tag + `index.yaml` + version-bump collisions that
   were specifically engineered out. It was tried, it was reverted, and
   the revert commit (`9861ccb`) is the boundary.
2. **Do not push to `main` directly** for P1 CCM hub-chart work. All
   chart-side P1 CCM iteration lands on this trunk first.
3. **Work in a dedicated worktree.** The primary `~/code/enbuild-helm`
   checkout has HEAD churn from other feature branches; the dedicated
   worktree `~/code/enbuild-helm.trunk` is the safe target.
4. **Consolidate parallel branches via file-level cherry-pick, not
   `git merge`.** The five legacy branches are rooted at `0.0.45` and
   merging them re-introduces base divergence.
5. **Bump the `-p1ccm-trunk.N` suffix on each chart-touching commit you
   want to publish.** The workflow keys off `Chart.yaml`'s `version:`;
   if the version doesn't change, `gh release create` will fail (tag
   already exists) and fall through to `--clobber` of the existing
   asset.
6. **Stage A Mongo overlay must be carried on every roll** via
   `EXTRA_HELM_ARGS=--values .../stage-a-mongo-handoff.values.yaml`
   until upstream chart default flips `mongodb.enabled` to false.

---

## When trunk eventually merges to main

(Weeks away. This is the forward plan, not a current action.)

When trunk is mature, the merge to `main` is a normal `git merge`. At
that point:

- `main`'s `helm-release.yml` fires on the merge commit (chart paths
  touched).
- `chart-releaser-action` resolves the next stable version via
  `semver.next_version('patch')` — it strips pre-release suffixes
  automatically, so `0.0.49-p1ccm-trunk.7` merging into a `main` at
  `0.0.47` produces `0.0.48` (next stable patch), not `0.0.49`.
- The new stable tag is `enbuild-0.0.48` on `main`. The existing
  trunk-side tags `enbuild-trunk-0.0.49-p1ccm-trunk.0..N` stay where
  they are on this branch — independent tag namespaces, no rename, no
  force-push, no cleanup required.
- Once merged and verified, the `helm-release-trunk.yml` workflow and
  this doc can be removed in a follow-up cleanup commit (or kept for
  the next contract — TBD with the user at that time).

---

## Reference paths

- Trunk worktree: `~/code/enbuild-helm.trunk`
- Trunk workflow: [`.github/workflows/helm-release-trunk.yml`](../.github/workflows/helm-release-trunk.yml)
- Main release workflow (do NOT modify for trunk): `.github/workflows/helm-release.yml`
- Audit log of every spoke roll: `~/.local/share/vendor13-ib-rolls.log`
- Per-env files (`p1-cluster-mgmt`):
  - `envs/vendor13-ib/chart-versions.pin`
  - `envs/vendor13-ib/stage-a-mongo-handoff.values.yaml`
  - `envs/vendor13-ib/identity.yaml`
  - `envs/vendor13-ib/NOTES.md`
- Roll wrapper: `~/code/p1-cluster-mgmt/scripts/roll-vendor13-ib.sh`
- Related issues: #64 (gRPC Service + VirtualService restoration),
  #65 (trunk publish flow)
- Related commits: `c507a89` (Wave 1 cherry-pick), `510758b` →
  `9861ccb` (reverted attempt to share `helm-release.yml`),
  `29297bb` (trunk-only publish workflow)
