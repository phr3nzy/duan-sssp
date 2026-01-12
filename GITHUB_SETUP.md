# GitHub Setup Guide

This guide will help you push this repository to GitHub.

## Prerequisites

- GitHub account created
- Git configured with your username and email:
  ```bash
  git config --global user.name "Your Name"
  git config --global user.email "your.email@example.com"
  ```

## Step 1: Create GitHub Repository

1. Go to https://github.com/new
2. Repository name: `duan-sssp`
3. Description: `O(m log^(2/3) n) Single-Source Shortest Path algorithm (Duan et al., 2025)`
4. Set to **Public** (or Private if you prefer)
5. **DO NOT** initialize with README, .gitignore, or license (we already have these)
6. Click "Create repository"

## Step 2: Initial Commit

```bash
cd /home/phr3nzy/go/src/github.com/phr3nzy/duan-sssp

# Stage all files
git add .

# Create initial commit
git commit -m "Initial commit: O(m log^(2/3) n) SSSP implementation

- Complete implementation of Duan et al. (2025) algorithm
- Comprehensive benchmark suite
- Full documentation (README, QUICKSTART, ALGORITHM, BENCHMARKS)
- CI/CD with GitHub Actions
- Stack overflow fix in FindPivots
- Version 1.0.1"
```

## Step 3: Push to GitHub

```bash
# Set main as default branch
git branch -M main

# Push to GitHub
git push -u origin main
```

## Step 4: Configure GitHub Repository Settings

### Enable GitHub Actions

1. Go to repository Settings → Actions → General
2. Under "Actions permissions", select "Allow all actions and reusable workflows"
3. Save changes

### Enable Issues and Discussions

1. Go to Settings → General
2. Under "Features", ensure these are checked:
   - ✅ Issues
   - ✅ Discussions (optional, for Q&A)
   - ✅ Projects (optional)

### Add Repository Topics

1. Go to repository main page
2. Click the gear icon next to "About"
3. Add topics:
   - `graph-algorithms`
   - `shortest-path`
   - `sssp`
   - `dijkstra`
   - `go`
   - `golang`
   - `algorithms`
   - `research`
   - `performance`

### Set Up Branch Protection (Optional)

1. Go to Settings → Branches
2. Add rule for `main` branch:
   - ✅ Require pull request reviews before merging
   - ✅ Require status checks to pass before merging
   - ✅ Require branches to be up to date before merging
   - ✅ Include administrators

## Step 5: Add Badges to README

Add these badges to the top of your README.md:

```markdown
[![CI](https://github.com/phr3nzy/duan-sssp/workflows/CI/badge.svg)](https://github.com/phr3nzy/duan-sssp/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/phr3nzy/duan-sssp)](https://goreportcard.com/report/github.com/phr3nzy/duan-sssp)
[![GoDoc](https://godoc.org/github.com/phr3nzy/duan-sssp?status.svg)](https://godoc.org/github.com/phr3nzy/duan-sssp)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Paper](https://img.shields.io/badge/arXiv-2504.17033-b31b1b.svg)](https://arxiv.org/abs/2504.17033)
```

Update with:
```bash
# Edit README.md to add badges at the top
# Then commit and push
git add README.md
git commit -m "docs: add status badges"
git push
```

## Step 6: First Tag/Release

Create your first release:

```bash
# Create annotated tag
git tag -a v1.0.1 -m "Release v1.0.1

Fixes:
- Stack overflow in FindPivots
- Cycle detection in tree size calculation

Features:
- Complete SSSP implementation
- Comprehensive benchmarks
- Full documentation
"

# Push tag
git push origin v1.0.1
```

This will trigger the Release workflow and create a GitHub release with binaries.

## Step 7: Register with Go Package Index

Your package will automatically appear on:
- https://pkg.go.dev/github.com/phr3nzy/duan-sssp

After first push, wait ~30 minutes for indexing, then verify at:
https://pkg.go.dev/github.com/phr3nzy/duan-sssp

## Verification Checklist

After pushing, verify:

- [ ] Repository accessible at https://github.com/phr3nzy/duan-sssp
- [ ] README renders correctly
- [ ] CI workflow runs and passes
- [ ] Go packages importable: `go get github.com/phr3nzy/duan-sssp`
- [ ] Documentation appears on pkg.go.dev
- [ ] Benchmarks run in CI
- [ ] License file present

## Common Issues

### Authentication Failed

If you get authentication errors:

**Using HTTPS:**
```bash
# Use personal access token instead of password
# Create token at: https://github.com/settings/tokens
```

**Using SSH:**
```bash
# Change remote to SSH
git remote set-url origin git@github.com:phr3nzy/duan-sssp.git

# Ensure SSH key is added to GitHub
# https://github.com/settings/keys
```

### CI Fails

If CI workflows fail:
1. Check Actions tab for error details
2. Fix issues locally
3. Commit and push fixes
4. CI will re-run automatically

### Go Module Issues

If `go get` fails:
```bash
# Ensure go.mod is correct
cat go.mod

# Should show:
# module github.com/phr3nzy/duan-sssp
# go 1.21

# If not, fix it:
go mod edit -module github.com/phr3nzy/duan-sssp
git add go.mod
git commit -m "fix: correct module path"
git push
```

## Next Steps

1. **Write a blog post** - Explain the algorithm
2. **Create tutorial** - Step-by-step guide
3. **Benchmark other implementations** - Compare performance
4. **Get feedback** - Share on Reddit, HN, etc.
5. **Respond to issues** - Help users get started

## Promotion Ideas

- Share on [r/golang](https://reddit.com/r/golang)
- Share on [r/algorithms](https://reddit.com/r/algorithms)
- Tweet with #golang #algorithms
- Submit to [Hacker News](https://news.ycombinator.com/submit)
- Add to [Awesome Go](https://github.com/avelino/awesome-go)

---

**Congratulations!** Your SSSP implementation is now on GitHub! 🎉
