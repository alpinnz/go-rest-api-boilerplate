# Contributing Guide

Thank you for your interest in contributing to **Go REST API Boilerplate**.

Please follow the steps below to start contributing:

---

## How to Contribute

1. **Fork this repository** to your GitHub account.
2. **Clone** the forked repository to your local machine:

   ```bash
   git clone https://github.com/username/go-rest-api-boilerplate.git
   ````

3. **Create a new branch** for your feature or fix:

   ```bash
   git checkout -b feature/your-feature-name
   ```
4. Make your changes and **commit** with a clear message:

   ```bash
   git commit -m "feat: add user registration endpoint"
   ```
5. **Push** your branch to your forked repository:

   ```bash
   git push origin feature/your-feature-name
   ```
6. Open a **Pull Request** from your forked repository to the main repository.

---

## Commit Message Guidelines

Follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

* `feat:` for new features
* `fix:` for bug fixes
* `docs:` for documentation changes
* `refactor:` for code changes that neither fix a bug nor add a feature
* `test:` for adding or updating tests
* `chore:` for maintenance tasks (deps, configs, etc.)
* `style:` for code style/formatting changes (no logic)
* `perf:` for performance improvements
* `build:` for build system or dependency changes
* `ci:` for CI/CD pipeline/configuration
* `revert:` for reverting a previous commit

**Examples:**

```
feat(auth): add refresh token
fix(user): fix duplicate email validation
```

---

## Coding Style

* Use **Go fmt** for code formatting.
* Ensure all tests pass before submitting:

  ```bash
  go test ./...
  ```

---

## Discussion & Issues

* Use the **Issues** tab on GitHub to report bugs or suggest new features.
* For large features, please discuss with the maintainers before submitting a PR.

---

Thank you for contributing!

```

Would you like me to also **translate and adapt `CHANGELOG.md`** into English in the same style (with Keep a Changelog format)?
```
