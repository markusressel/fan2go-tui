# Releasing fan2go-tui

fan2go-tui follows [semantic versioning](https://semver.org/):

- **major** (`1.4.2` → `2.0.0`): breaking changes,
- **minor** (`1.4.2` → `1.5.0`): new features, without breaking changes,
- **patch** (`1.4.2` → `1.4.3`): bug fixes only.

Releases up to `0.4.0` were pre-1.0 releases, without these rules. The first release with breaking changes after
`0.4.0` is `1.0.0`; from then on, the rules above apply strictly.

## During development: recording breaking changes

Every pull request has a **"Breaking changes / upgrade notes"** section (see
[the PR template](.github/pull_request_template.md)). A change is breaking if an existing user may notice it after
upgrading without changing anything: config semantics or defaults, CLI flags, output or exit codes, key bindings, the
saved UI state, or which versions of the fan2go daemon are supported.

For a breaking change:

1. describe it in the PR's "Breaking changes / upgrade notes" section, written for users: what changes, who is
   affected, what they have to do, and
2. add the `breaking-change` label to the PR.

[Release Drafter](.github/release-drafter.yml) collects merged PRs into a draft release, lists PRs labeled
`breaking-change` in a "⚠️ Breaking changes" section and proposes the next version from the labels (see
[the version resolver](.github/release-drafter.yml)).

## Cutting a release

1. **Collect the upgrade notes.** List the merged PRs labeled `breaking-change` since the last release:

   ```shell
   gh pr list --state merged --label breaking-change --search "merged:>$(git log -1 --format=%cI $(git describe --tags --abbrev=0))" --json number,title,url
   ```

   Read their "Breaking changes / upgrade notes" sections.

2. **Add an upgrade warning for the Arch packages** if there are breaking changes. Append an entry for the new version
   at the end of the notes in [`packaging/aur/fan2go-tui.install`](packaging/aur/fan2go-tui.install):

   ```text
   == 1.0.0
   ui.updateInterval was renamed to ui.refreshInterval. Rename it in your fan2go-tui.yaml.
   ```

   pacman shows these lines when a user upgrades across that version (both `fan2go-tui` and `fan2go-tui-git`). Keep
   each entry short and actionable; link to the release notes for details. Merge this change before publishing the
   release.

3. **Update the Justfile default version** (`VERSION :=` in the [`Justfile`](Justfile)) to the new version, so local
   builds from the tag report it. CI and the Arch packages pass the version explicitly, so release artifacts are
   correct either way.

4. **Publish the GitHub release; this creates the tag.** Open the draft release that Release Drafter keeps up to date
   (Releases page, or `gh release list`). Then:
   - fill the "Upgrade notes" section from step 1, or remove it if there are no breaking changes,
   - check the proposed version against the rules above and correct it if needed (e.g. if a PR was labeled wrong),
   - set the tag to the version (e.g. `1.0.0`), targeting `main`, with the same release title, and publish.

   Publishing creates the tag. Don't push the tag yourself: the Go workflow would then create a second release instead
   of attaching the binaries to this one. The tag push triggers:
   - the [Go workflow](.github/workflows/go.yml): builds the binaries (with the version from the tag) and attaches
     them to the published release (a few minutes later), and
   - the [AUR Packages workflow](.github/workflows/aur-packages.yml): test-builds the `fan2go-tui` (from the tag's
     source tarball) and `fan2go-tui-git` packages from [`packaging/aur`](packaging/aur). It doesn't publish anything.

5. **Publish the AUR packages** once the AUR Packages workflow for the tag has passed. On an Arch Linux machine with
   your AUR SSH key:

   ```shell
   packaging/aur/publish.sh 1.0.0
   ```

   The script updates your AUR clones (`$AUR_DIR/fan2go-tui` and `$AUR_DIR/fan2go-tui-git`, default
   `AUR_DIR=~/programming/AUR`; missing clones are created), sets `pkgver` and the checksum of the stable package,
   test-builds both packages in a temporary directory, shows the diff and pushes after you confirm. Answering "no"
   restores the clones. See `packaging/aur/publish.sh --help` for options (e.g. `--pkgrel 2` for a packaging-only fix
   of the same version).

## AUR packages

The PKGBUILDs live in [`packaging/aur`](packaging/aur) and are the source of truth; don't edit them on the AUR
directly, the next release overwrites them.

- `fan2go-tui`: built from the release tarball. Its `pkgver` and checksum are set by `publish.sh`.
- `fan2go-tui-git`: built from `main`. It contains unreleased changes, including breaking ones; its upgrade warnings
  only fire once the corresponding version has been tagged.

Publishing is deliberately manual (`packaging/aur/publish.sh`): an AUR SSH key grants push access to every package of
the account, and a PKGBUILD runs code on users' machines, so the key is not stored in GitHub.

Pull requests that change `packaging/aur` and release tags build both packages in CI.
