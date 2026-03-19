# FAQ Plugin Test Suite

## Setup

1. Install dependencies:

```bash
npm install
```

2. Copy environment template:

```bash
cp .env.example .env
```

3. Set `BASE_URL` and `API_URL` for your environment.

## Running Tests

```bash
npm run test:e2e
```

```bash
npx playwright test --headed
```

```bash
npx playwright test --debug
```

## Structure

- `tests/e2e/` - End-to-end and API smoke tests
- `tests/support/fixtures/` - Fixtures and factories with auto-cleanup
- `tests/support/helpers/` - Shared helper functions

## Best Practices

- Use `data-testid` selectors for stability.
- Register network intercepts before navigation.
- Keep tests < 300 lines and < 1.5 minutes.
- Seed data via API helpers, not UI flows.

## Notes

- Reports: `playwright-report/`
- Artifacts: `test-results/`
