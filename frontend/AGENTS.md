# Frontend Agent Guide

This guide applies to `frontend/`, the TanStack Start frontend for the Camp 2026 game. Follow it when adding, changing, or reviewing frontend code.

## Priority Order

1. Preserve the architecture boundaries.
2. Prefer existing shadcn/ui-owned components and local shared components.
3. Keep colors, spacing, data fetching, schemas, and API errors centralized.
4. Add tests around behavior, not implementation details.
5. Run the quality commands before finishing.

Do not optimize for short-term convenience if it introduces scattered styles, ad hoc fetch calls, or cross-feature imports.

## Commands

Run commands from `frontend/`.

```sh
pnpm dev
pnpm lint
pnpm format
pnpm test
pnpm build
pnpm quality
```

Use `pnpm format:write` and `pnpm lint:fix` only when intentionally rewriting formatting or fixable lint issues.

## Architecture

Current shape:

```text
src/
  app/          app composition and providers
  routes/       TanStack file routes and route API handlers
  pages/        route-level screens composed from features and shared UI
  features/     user-facing feature modules
  shared/       reusable API, config, lib, utilities, and UI primitives
  styles/       global CSS variables and Tailwind theme tokens
  testing/      test setup and shared test helpers
```

Import direction must stay one-way:

```text
routes -> pages -> features -> shared
app -> shared
```

Rules:

- `shared/` must not import from `app/`, `routes/`, `pages/`, or `features/`.
- `features/` must not import from `app/`, `routes/`, or `pages/`.
- `pages/` must not import from `app/` or `routes/`.
- Route files should be thin: URL/search/loader/API handler wiring only.
- Pages should compose features; they should not contain API calls.
- Features own feature-specific schema, query options, hooks, and UI.

## UI And Components

Always check for an existing component before creating a new one.

Use in this order:

1. Existing component in `src/shared/ui`.
2. shadcn/ui component installed into `src/shared/ui`.
3. Small feature-local component under `features/<feature>/ui`.
4. New shared component only when at least two places need it or the abstraction is clearly reusable.

If shadcn/ui has the component, install/copy the shadcn version and adapt it minimally. Do not hand-roll buttons, cards, badges, dialogs, forms, inputs, tables, dropdowns, skeletons, alerts, or tooltips when shadcn already provides the pattern.

Shared UI rules:

- shadcn components live in `src/shared/ui`.
- Keep shadcn component APIs close to upstream unless there is a concrete local reason.
- Wrap shadcn components for app-specific semantics instead of modifying every caller.
- Prefer composition over large prop-heavy components.
- Keep component names concrete: `MetricCard`, `InfoTile`, `StatusBadge`.
- Do not put business-specific wording into shared components.

## Color Palette

Colors must be semantic and centralized.

Allowed color sources:

- CSS variables in `src/styles/app.css`.
- Variant mappings in `src/shared/config/color-palette.ts`.
- Existing shadcn semantic tokens such as `bg-background`, `text-foreground`, `bg-card`, `text-muted-foreground`, `border-border`, `text-primary`, and `text-destructive`.

Disallowed in feature/page code:

- Raw Tailwind palette classes like `bg-emerald-500`, `text-blue-800`, `border-violet-200`.
- Hard-coded hex, rgb, hsl, or oklch values outside `src/styles/app.css`.
- One-off color strings inside component arrays.

When a new color is needed:

1. Add a semantic CSS variable in `src/styles/app.css`.
2. Expose it through `@theme inline` as a Tailwind token.
3. Add or extend a variant mapping in `src/shared/config/color-palette.ts`.
4. Consume the variant through a component such as `StatusBadge`, `StatusDot`, `MetricCard`, or `IconBadge`.

## Icons

- Use `lucide-react` for interface, action, status, navigation, and toolbar icons.
- Use `@iconify/react` only when the icon is from a broader visual set that lucide does not cover well.
- Do not inline SVGs for standard icons.
- Icon-only buttons must have an accessible label.
- Decorative icons must use `aria-hidden`.

## Data Fetching

Do not call `fetch` directly from React components.

Use this flow:

```text
route loader or component hook
  -> feature queryOptions
  -> shared apiClient
  -> Zod parse
  -> TanStack Query cache
  -> UI
```

Rules:

- Query keys must be stable arrays and live beside `queryOptions`.
- Runtime response validation must use Zod at the API boundary.
- Use `apiClient` for HTTP behavior and error normalization.
- Server state belongs in TanStack Query, not React context.
- UI-only state can use local state or a small context when truly shared.

## Zod Rules

- Use Zod 4 APIs.
- Do not use deprecated chained string format methods for URL, email, UUID, ISO datetime, ISO date, ISO time, or ISO duration validation.
- Use top-level format schemas instead:
  - `z.url()`
  - `z.email()`
  - `z.uuid()`
  - `z.iso.datetime()`
  - `z.iso.date()`
  - `z.iso.time()`
  - `z.iso.duration()`
- Keep schemas near the API or feature boundary that owns the contract.
- Prefer `z.output<typeof Schema>` or `z.infer<typeof Schema>` for parsed values.
- Prefer `z.input<typeof Schema>` for form/API inputs before parsing.

## Routes And API Routes

- Keep TanStack route files thin.
- Put page layout in `pages/`.
- Put reusable workflow UI in `features/`.
- Put API handlers under `src/routes/api`.
- Mock API routes are allowed only while the backend contract is missing. Keep their response shape validated by the same Zod schema the client consumes.
- Do not edit `src/routeTree.gen.ts`; it is generated.

## Testing

Use Vitest and Testing Library.

Test priorities:

1. Zod schemas accept valid payloads and reject invalid payloads.
2. Feature UI renders loading, success, and error states.
3. User interactions call the expected callbacks or trigger the expected query/mutation behavior.
4. API adapters normalize success and error payloads.

Testing rules:

- Prefer user-visible text and roles over implementation details.
- Do not test class names unless the behavior is specifically visual and cannot be tested otherwise.
- Each test that uses TanStack Query should get its own `QueryClient`.
- Disable retries in test query clients unless retry behavior is the test subject.

## Formatting And Linting

Prettier owns formatting. Do not manually align code or reorder Tailwind classes by hand.

ESLint owns:

- React Hooks correctness.
- Layer import boundaries.
- Avoiding raw Tailwind palette classes in source files.
- TypeScript recommended correctness rules.

Before finishing frontend changes, run:

```sh
pnpm lint
pnpm format
pnpm test
pnpm build
```

Use `pnpm quality` when you want the full sequence.

## Adding New Code

Before adding a file:

1. Search for an existing component/helper/schema/query that already does the job.
2. Place the file at the lowest layer that can own it.
3. Keep feature-specific code in the feature.
4. Promote to shared only when reuse is real.
5. Add tests if the behavior can break.

Avoid:

- `utils.ts` files that collect unrelated helpers.
- `components/` dumping grounds.
- Cross-feature imports for convenience.
- Duplicating server state in local stores.
- Hidden side effects at module import time.

## Accessibility And UX

- Buttons must use `<Button>`.
- Cards must use `<Card>` and its subcomponents.
- Badges must use `<Badge>` or a semantic wrapper such as `<StatusBadge>` or `<IconBadge>`.
- Loading placeholders must use `<Skeleton>`.
- Keep text responsive and avoid overflow.
- Do not rely on color alone to communicate status.
- Prefer concise Traditional Chinese UI copy for this project.
