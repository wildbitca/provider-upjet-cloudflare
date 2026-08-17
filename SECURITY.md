# Security Policy

## Reporting a vulnerability

**Do not open a public issue for a security problem.** An issue is visible to everyone the
moment it is created, including to whoever would exploit it.

Use GitHub's private vulnerability reporting instead — it is enabled on this repository:

1. Go to the **Security** tab.
2. **Report a vulnerability**.
3. Describe what you found, how to reproduce it, and what an attacker gains.

The report is visible only to the maintainers. You will get an acknowledgement, and a fix
or an explanation of why it is not one, before anything is made public.

## Scope

This repository is a Crossplane provider: it holds control-plane code that talks to a
third-party API with credentials supplied by whoever runs it. Reports about credential
handling, privilege escalation through a Composition, or leakage of secret material into
logs or into the status of a managed resource are in scope and are the ones we care about
most.

Vulnerabilities in the upstream Terraform provider that this one wraps should be reported
to that project; we will pick up the fix when we bump it.

## Supported versions

Only the latest released version is supported. There are no backports to older tags: the
fix ships in a new release.
