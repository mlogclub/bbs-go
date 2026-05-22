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
