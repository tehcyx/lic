# Validate locally before pushing

Run `tools/verify-local.sh` from this checkout. The default runs every validation group. Available groups: all. A named group is partial evidence, not a passing full suite.

Install Go. The command runs tests and vet. Dependency PRs now require local verification and human merge; the unattended merge helper is removed. Public CodeQL remains enabled.

## Hosted behavior

Ordinary personal validation no longer runs on pushes or pull requests. The changed validation workflows retain an explicit workflow_dispatch fallback with their tool setup and checks. Dispatch consumes hosted compute; it is not needed for routine development. Existing publication, release, public security scanning, dependency updates, and operational schedules remain enabled unless stated above. A green publication run now means the artifact built, not that local validation ran. Record the exact tested commit and command in your PR evidence before a human merges.

## Check this policy offline

Install PyYAML 6.0.3 in your development Python environment, then run:

```sh
python3 tools/test_ci_policy.py
```

These tests inspect triggers, job dependencies, manual validation gates, and the local command interface. They do not dispatch Actions, publish an image, or replace application tests. Local validation does not grant permission to merge or deploy.

## Rollback

Revert this repository's local-validation PR to restore its previous automatic checks. Restore any removed dependency-auto-merge helper only together with its test gate. No database, image tag, deployed chart, or persistent data needs to change for rollback.
