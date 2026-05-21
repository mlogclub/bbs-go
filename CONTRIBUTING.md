# Contributing to BBS-GO

Thank you for your interest in contributing to BBS-GO.

BBS-GO is an open-source, self-hosted community platform for forums, Q&A, and reusable knowledge. We welcome bug reports, documentation improvements, feature ideas, and code contributions.

## Ways to Contribute

- Report bugs through [GitHub Issues](https://github.com/mlogclub/bbs-go/issues)
- Ask usage questions in [GitHub Discussions](https://github.com/mlogclub/bbs-go/discussions)
- Improve documentation and examples
- Submit fixes and small improvements
- Propose larger features before implementing them

## Before Opening a Pull Request

1. Search existing issues and discussions to avoid duplicates.
2. For larger changes, open an issue or discussion first and describe the proposal.
3. Keep pull requests focused and easy to review.
4. Include screenshots for UI changes when possible.
5. Update documentation when behavior changes.

## Development Notes

- Go backend lives in `server/`.
- React frontend/dashboard lives in `web/`.
- Dashboard routes are under `/dashboard`; the old `admin/` frontend is no longer maintained.
- User-facing strings should support both `en-US` and `zh-CN` where applicable.

Useful checks:

```bash
# Backend
cd server && go test ./...

# Frontend
cd web && pnpm lint && pnpm typecheck
```

For documentation site changes:

```bash
cd bbs-go-docs && pnpm build
```

## Licensing of Contributions

BBS-GO is open source under GPLv3, with commercial licensing available for teams that need proprietary redistribution, custom commercial terms, deployment support, customization, migration, or long-term maintenance.

By contributing to BBS-GO, you agree that your contributions may be distributed under the project's open-source license and may also be used as part of BBS-GO's commercial licensing offerings.

Please only contribute code, documentation, assets, or other materials that you have the right to submit.

## Commercial Licensing and Private Support

For commercial licensing, deployment support, customization, migration, or long-term maintenance, email:

<g330721072@gmail.com>

For security reports, email the same address with `[Security]` in the subject instead of opening a public issue.

## Code of Conduct

Be respectful and constructive. Assume good intent, keep discussions technical and practical, and help make the project welcoming for both users and contributors.
