# Product Requirements Document (PRD)

## FinTrack — Personal Finance SaaS

**Status:** Draft MVP
**Version:** 1.0
**Date:** October 27, 2026
**Product Owner:** FinTrack Team
**MVP Target:** 8 weeks, with development time of 1–2 hours per day

## 1. Executive Summary

FinTrack is a SaaS application designed to help individuals and small households record transactions, understand cash flow, manage budgets, and make financial decisions more easily.

The MVP version focuses on a simple core experience: users create an account, add financial accounts, record income/expenses, and view a summary of their financial condition via a dashboard.

FinTrack is also designed as an end-to-end portfolio project to demonstrate the ability to build a real product using Next.js, Go, and PostgreSQL

## 2. Background and Problem

Many people want to manage their finances, but:

- Transaction recording is often inconsistent.
- Existing financial apps feel too complex.
- Users find it difficult to identify their largest spending categories.
- Users lack a clear overview of their monthly cash flow.
- Spreadsheets are flexible but require manual setup and maintenance.

### Problem Statement

> Users need a simple and fast way to record transactions and understand where their money is going without having to use complex accounting software.

## 3. Product Goals

### Primary Goals

- Help users record transactions quickly.
- Provide an easy-to-understand overview of balances and cash flow.
- Help users control spending through simple budgeting.
- Build a product foundation that can be evolved into a paid SaaS.

### Learning and Portfolio Goals

- Implement Go for a structured and tested REST API backend.
- Build a responsive and usable Next.js frontend.
- Understand authentication, authorization, validation, migration, testing, and deployment.
- Demonstrate security considerations for financial data.

### Out of MVP Scope

- Becoming a professional accounting application.
- Automatically integrating bank accounts.
- Providing investment or financial advice.
- Storing card data, bank credentials, or sensitive payment information.
- Providing automatic account reconciliation

## 4. Target Users

### Primary Persona: Young Professionals or Freelancers

- Have income from one or several sources.
- Use bank accounts and e-wallets.
- Want to know their biggest expenses each month.
- Do not want to learn accounting software.
- Willing to try a free, easy-to-use product.

### Secondary Persona: Small Households

- Need an overview of shared expenses.
- During MVP, can manually use one shared account.
- Multi-user and role features will be considered after initial validation.

## 5. Jobs to Be Done

- When receiving income, I want to record it so that my balance is accurate.
- When making a purchase, I want to record it quickly so I don't forget.
- When evaluating the current month, I want to see the largest spending categories.
- When making a monthly plan, I want to set spending limits per category.
- When making a financial decision, I want to see the current cash flow status.

## 6. User Journey MVP

```text
Landing page
  → Register
  → Login
  → Brief onboarding
  → Create financial account
  → Add first category/transaction
  → View dashboard
  → Review budget and transaction history
```

### Activation Moment

A user is considered active if they have:

1. Created an account.
2. Created at least one financial account.
3. Recorded their first transaction.

## 7. MVP Scope

### P0 — Mandatory

#### Authentication

- Registration with email and password.
- Login and logout.
- Secure user sessions.
- User profile endpoint.
- Clear error messages without leaking sensitive information.

#### Financial Accounts

- Create accounts: cash, bank, e-wallet, or custom types.
- View account list and balances.
- Edit account names.
- Deactivate or delete accounts that have no transactions.

#### Categories

- Default categories for income and expenses.
- Create custom categories.
- Edit or deactivate custom categories.

#### Transactions

- Record income or expenses.
- Select account and category.
- Fill in amount, date, and optional notes.
- View a list of recent transactions.
- Filter by date, type, account, and category.
- Edit and delete transactions.
- (Out of MVP Scope) Transfers between financial accounts (will be implemented as two separate transactions: an expense from Account A and income to Account B).

#### Dashboard

- Total balance.
- Total income for the selected period.
- Total expenses for the selected period.
- Net cash flow.
- List of recent transactions.
- Expense breakdown per category.

#### Budget

- Set monthly budgets per expense category.
- View budget utilization.
- Display status: safe, approaching limit, or exceeded limit.

### P1 — After Initial MVP

- Export transactions to CSV.
- Recurring transactions.
- Household sharing.
- Dark theme.
- Email notifications.
- Savings targets.

### P2 — Business Exploration

- CSV import from banks or e-wallets.
- Open banking integration.
- Premium packages and billing.
- PDF reports.
- Multi-currency.
- Automated insights.

## 8. Functional Requirements

### FR-01 — Registration

Description: Users can create a new account.

Acceptance Criteria:

- [ ] Email is required and must have a valid format.
- [ ] Email must be unique.
- [ ] Password meets minimum length and security rules.
- [ ] Password is stored as a hash, not plaintext.
- [ ] After success, the user can login.
- [ ] Validation errors are displayed securely and are easy to understand.

### FR-02 — Login and Logout

Acceptance Criteria:

- [ ] Users can login with valid credentials.
- [ ] Invalid credentials result in a generic response.
- [ ] Sessions are stored using a secure mechanism.
- [ ] Logout revokes or deletes the active session.
- [ ] Pages requiring authentication cannot be accessed without a session.

### FR-03 — Manage Financial Accounts

Acceptance Criteria:

- [ ] Users can create an account with name, type, and currency.
- [ ] Balance is displayed based on transactions owned by the user.
- [ ] Users can only view and modify their own accounts.
- [ ] Other users' accounts cannot be accessed via ID manipulation.
- [ ] Accounts with transactions cannot be deleted indiscriminately.

### FR-04 — Manage Transactions

Acceptance Criteria:

- [ ] Users can choose between income or expense types.
- [ ] Amount must be greater than zero.
- [ ] Account and category must be valid and owned by the user.
- [ ] Transaction date must be valid.
- [ ] Transactions are saved and appear in the list.
- [ ] Transaction changes update the balance summary.
- [ ] Users cannot read, edit, or delete other users' transactions.

### FR-05 — Dashboard

Acceptance Criteria:

- [ ] Dashboard displays a clear period.
- [ ] Balance, income, expenses, and net cash flow are consistent with transaction data.
- [ ] Empty data has helpful empty states.
- [ ] Charts have understandable labels.
- [ ] Dashboard is usable on mobile screens.

### FR-06 — Budget

Acceptance Criteria:

- [ ] Users can create or modify a budget per category and month.
- [ ] A user cannot have duplicate budgets for the same category and month.
- [ ] Budget usage is calculated from expense transactions in that period.
- [ ] Budget status changes based on actual usage.

## 9. Non-Functional Requirements

### Security

- Passwords use Argon2id or bcrypt.
- Cookies/sessions use secure configurations.
- All private endpoints have authorization checks.
- Validation is performed on both frontend and backend, with the backend as the source of truth.
- Rate limiting is applied to login and sensitive endpoints.
- Secrets are not stored in the repository.
- Logs do not record passwords, tokens, or unnecessary financial data.

### Performance

- Dashboard endpoints respond quickly for MVP datasets.
- Transaction queries use indexes on user_id, date, account, and category as needed.
- Transaction lists use pagination or limits.
- Frontend does not fetch more data than necessary.

### Reliability

- Database migrations can be run repeatedly and safely.
- Database errors are handled without crashing the entire server.
- Health checks are available.
- Database backup and restore are documented before production.

### Accessibility and Usability

- Forms have clear labels.
- Basic keyboard navigation works.
- Color is not the only way to distinguish status.
- Amounts and dates have a consistent format.
- Loading, error, and empty states are available.

## 10. High-Level Data Model

```text
User 1 ─── * Account
User 1 ─── * Category
User 1 ─── * Transaction
Account 1 ─── * Transaction
Category 1 ─── * Transaction
User 1 ─── * Budget
Category 1 ─── * Budget
```

### Core Entities

| Entity | Key Fields | Notes |
|---|---|---|
| User | id, email, password_hash, created_at | Unique email |
| Account | id, user_id, name, type, currency, current_balance | Owned by one user, current_balance for query optimization |
| Category | id, user_id, name, kind | income or expense |
| Transaction | id, user_id, account_id, category_id, amount, kind, occurred_at | Amount as integer/smallest unit |
| Budget | id, user_id, category_id, month, limit_amount | Unique combination of user-category-month |
| Session | id, user_id, token_hash, expires_at | Can be revoked |

### Financial Rules

- Do not use floating point for monetary amounts.
- Use the smallest integer unit for MVP IDR.
- Store currency explicitly.
- Use UTC for time storage.
- All aggregate calculations must originate from server data.
- Do not store bank credentials or card numbers.
- Account balance updates must be performed atomically (using database transactions) alongside transaction creation to maintain data consistency.

## 11. Desain Teknis Tingkat Tinggi

```text
Next.js Web
  ├── Auth Pages
  ├── Dashboard
  ├── Transactions
  ├── Accounts
  └── Budget
          │ HTTPS/JSON
          ▼
Go REST API
  ├── HTTP handlers
  ├── Middleware auth/CORS/rate limit
  ├── Application services
  ├── Repositories
  └── PostgreSQL
```


### Frontend Library Decisions

- Next.js App Router and TypeScript as the application foundation.
- Tailwind CSS for responsive styling.
- shadcn/ui for UI components that reside in the repository and are customizable.
- lucide-react for icons.
- React Hook Form and Zod with @hookform/resolvers for forms and validation.
- TanStack Query for fetching, caching, mutation, and data invalidation from the Go API.
- Recharts for cash flow charts and category breakdowns.
- date-fns for date formatting and period calculations.
- Playwright for smoke tests and end-to-end flows.

For MVP, Redux or Zustand are not yet needed. Use TanStack Query for server state and React state for local component state.

### Minimum API Contract

```text
POST   /api/v1/auth/register
POST   /api/v1/auth/login
POST   /api/v1/auth/logout
GET    /api/v1/me
GET    /api/v1/dashboard/summary
GET    /api/v1/accounts
POST   /api/v1/accounts
PATCH  /api/v1/accounts/:id
GET    /api/v1/categories
POST   /api/v1/categories
GET    /api/v1/transactions
POST   /api/v1/transactions
PATCH  /api/v1/transactions/:id
DELETE /api/v1/transactions/:id
GET    /api/v1/budgets?month=YYYY-MM
PUT    /api/v1/budgets/:categoryID
```

## 12. UX and Pages

### Public Pages

- Landing page.
- Login.
- Register.
- Simple privacy/security note.

### Private Pages

- Dashboard.
- Transactions.
- Accounts.
- Categories.
- Budgets.
- Settings/profile.

### UX Principles

- Adding a transaction should be completed in a few steps.
- Default date is today.
- Transaction types are distinguished by labels and icons, not just color.
- Amounts always have a clear currency format.
- After saving a transaction, users receive confirmation and see the updated data.
- Empty states guide users to the next action

## 13. Analytics and Success Metric

### Product Metrics

- **Activation rate:** percentage of users who create an account, a financial account, and their first transaction.
- **Day-7 retention:** percentage of users who return to the product on the 7th day.
- **Weekly active users:** number of weekly active users.
- **Transactions per active user:** average number of transactions per active user.
- **Budget adoption:** percentage of active users who create a budget.
- **Task completion:** time and success rate of recording a transaction.

### Initial Validation Targets

- At least 5 users try the MVP.
- At least 3 users complete the activation flow without direct assistance.
- At least 2 users return to the application within 7 days.
- No critical findings related to cross-user data access.
- Users can record a main transaction in less than one minute after understanding the UI

## 14. 8-Week Delivery Plan

| Week | Focus | Deliverable |
|---|---|---|
| 1 | Discovery and Setup | Product brief, repo, frontend/backend skeleton |
| 2 | Database and Backend Foundation | Migration, connection, repository, error format |
| 3 | Auth | Register, login, logout, protected route |
| 4 | Core Transaction Flow | Accounts, categories, transaction CRUD |
| 5 | Insight | Dashboard, chart, budget |
| 6 | Quality | UX states, responsive layout, tests, security review |
| 7 | Release | CI, deployment, API docs, seed data |
| 8 | Validation | User testing, feedback, portfolio case study |

## 15. Risks and Mitigations

| Risk | Impact | Mitigation |
|---|---|---|
| Scope too large | MVP delayed | Prioritize P0 and limit external integrations |
| Inconsistent balance calculation | Loss of user trust | Use server-side aggregation and test invariants |
| Cross-user data leak | High security risk | Authorization check on every resource query |
| Sensitive financial data | Privacy risk | Minimize data, transport encryption, secret management |
| Limited learning time | Inconsistent progress | Small daily tasks and clear definition of done |
| Deployment takes too long | Demo delayed | Deploy skeleton earlier in weeks 2–3 |
| No need validation | Product not used | Test with users before adding P1/P2 features |

## 16. MVP Definition of Done

The MVP is considered complete when:

- [ ] New users can register and login.
- [ ] Users can create financial accounts.
- [ ] Users can record, modify, filter, and delete transactions.
- [ ] Dashboard calculates balance and cash flow correctly.
- [ ] Users can create monthly budgets.
- [ ] Authorization tests prove data between users is isolated.
- [ ] Frontend is responsive and has loading/error/empty states.
- [ ] Backend has relevant unit/integration tests.
- [ ] Migration and seed data are available.
- [ ] Application is deployed and can be tried via a demo URL.
- [ ] README and API documentation can be followed by other developers.
- [ ] No secrets or personal data are in the repository.


## 17. Next Steps

1. Review PRD and mark features that are truly necessary for the first version.
2. Create docs/product.md based on the problem, persona, and user journey sections.
3. Convert FRs into small issues/tasks in the repository.
4. Start with the product brief and repository setup in the first week of the roadmap.
5. Before building additional features, validate the register → first transaction → dashboard flow.

## Conclusion

The FinTrack MVP must prove one main thing: users can record their finances easily and gain a better understanding from that data. Any feature outside the core flow must wait until there is evidence of user need, keeping the project realistic to complete while remaining strong as a portfolio and a foundation for a SaaS business.
