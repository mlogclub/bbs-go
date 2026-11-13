# BBS-GO Web

React Router + Vite frontend for BBS-GO. The same source supports SPA output and
SSR output.

## Scripts

```bash
pnpm dev
pnpm build:spa
pnpm build:ssr
pnpm start
```

## Adding components

To add components to your app, run the following command:

```bash
npx shadcn@latest add button
```

This will place the ui components in the `components` directory.

## Using components

To use the components in your app, import them as follows:

```tsx
import { Button } from "@/components/ui/button";
```

## HTTP error statuses

Use `routeApiFetch` in route loaders to propagate API failures as HTTP responses.
Keep `apiFetch` for interactive forms that handle `ApiError` locally. Rendering
`ErrorPage` only changes the UI; loaders must throw a `Response` to set the
initial document's status. Backend resource endpoints return explicit 400/403/404
statuses and matching business codes. SSR preserves upstream 5xx statuses.

The default Docker/Compose entry point uses SSR on port 3000. Use this entry point
when resource-specific statuses matter (for example, a deleted topic must return
404). Client-side navigation cannot change an already received document status.

`pnpm build:spa` generates `.route-patterns.json` for the Go fallback. Unknown URL
shapes and missing static assets return 404; known page routes return the SPA
shell. Pure SPA mode cannot determine resource existence or API availability
before sending that shell, so dynamic resource errors require SSR for correct
document statuses. Rebuild SPA assets when adding or removing routes.

Regression check (mock API, no running backend needed):

```bash
pnpm build:ssr
node scripts/test-http-status.mjs
```
