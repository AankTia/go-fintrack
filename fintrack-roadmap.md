# FinTrack MVP Roadmap (8 Minggu)

## 📋 Overview

| Item | Detail |
|------|--------|
| **Project** | FinTrack – Personal Finance SaaS |
| **Timeline** | 8 Minggu (1-2 jam/hari = 160 jam total) |
| **Backend** | Go (REST API) + PostgreSQL |
| **Frontend** | Next.js App Router + TypeScript |
| **Goal** | MVP launch + Learning portfolio showcase |

---

## 🎯 WEEK 1: Discovery & Setup (Day 1-3)

### Goal: Product understanding + project scaffold ready

#### Day 1: Environment Setup & Project Structure
| Task | Time |
|------|-------|
| Install Go runtime | 15 min |
| Initialize Go module | 2 min |
| Create project structure | 5 min |
| Setup Docker Compose for Postgres local DB | 20 min |

#### Day 2: Configuration & Database Connection
| Task | Time |
|------|-------|
| Create config loader (env vars) | 25 min |
| Setup database connection pool | 30 min |
| Write first migration (table structure) | 40 min |
| Migrate to database | 10 min |

#### Day 3: Basic REST API (Hello World)
| Task | Time |
|------|-------|
| Create HTTP server with `net/http` | 20 min |
| Define API routes (chi or gin) | 30 min |
| Implement `/api/v1/auth/register` endpoint | 45 min |
| Implement `/api/v1/auth/login` endpoint | 50 min |
| Test with Postman collection or curl | 20 min |

---

## 🎯 WEEK 2: Database Design & Core Entities (Day 4-6)

### Goal: All database models + CRUD operations defined

#### Day 4: Database Schema & Models
| Task | Time |
|------|-------|
| Create PostgreSQL migration files | 60 min |
| Define Go structs for DB models | 30 min |
| Create repository interfaces | 25 min |
| Implement repository concrete structs | 40 min |

#### Day 5: Account Management (CRUD)
| Task | Time |
|------|-------|
| Create Account handler (GET/POST) | 40 min |
| Create Account service (business logic) | 30 min |
| Create Account repository implementation | 25 min |
| Test account creation & retrieval | 15 min |

#### Day 6: Category Management (CRUD)
| Task | Time |
|------|-------|
| Create Category handler (POST) | 40 min |
| Create Category service (business rules) | 30 min |
| Create Category repository (DB queries) | 25 min |
| Write unit tests for business logic | 20 min |

---

## 🎯 WEEK 3: Transaction Core (Day 7-9)

### Goal: Full transaction lifecycle complete

#### Day 7: Transaction Create & Validation
| Task | Time |
|------|-------|
| Define Transaction model & validation rules | 20 min |
| Create Transaction handler (POST) | 45 min |
| Implement transaction service logic | 35 min |
| Write validation & business rule tests | 20 min |

#### Day 8: Transaction List & Filtering
| Task | Time |
|------|-------|
| Implement GET list endpoint | 40 min |
| Add pagination & sorting | 30 min |
| Write integration test for filtering logic | 25 min |

#### Day 9: Transaction Update & Delete
| Task | Time |
|------|-------|
| PATCH endpoint for update (idempotent) | 35 min |
| DELETE endpoint | 30 min |
| Verify data integrity after deletion | 15 min |

---

## 🎯 WEEK 4: Dashboard & Analytics (Day 10-12)

### Goal: Visual insights + cash flow calculation

#### Day 10: Dashboard Summary API
| Task | Time |
|------|-------|
| Create Dashboard handler (GET) | 35 min |
| Implement service aggregation logic | 40 min |
| Query recent transactions for timeline | 25 min |

#### Day 11: Category Breakdown & Charts
| Task | Time |
|------|-------|
| Create category breakdown endpoint | 40 min |
| Implement expense breakdown service | 35 min |

#### Day 12: Security & Performance Testing
| Task | Time |
|------|-------|
| Security test: unauthorized access attempt | 30 min |
| Rate limiting for auth endpoints | 25 min |
| Load test dashboard aggregation | 20 min |

---

## 🎯 WEEK 5: Budget Management (Day 13-15)

### Goal: Monthly budget tracking + limits

#### Day 13: Budget Model & CRUD
| Task | Time |
|------|-------|
| Define Budget model (amount, limit, period) | 25 min |
| Create Budget handler (GET/POST) | 40 min |
| Implement budget service logic | 35 min |

#### Day 14: Budget Calculation & Status
| Task | Time |
|------|-------|
| Implement budget status endpoint | 35 min |
| Create budget usage aggregation service | 30 min |

#### Day 15: Budget Override & Notification
| Task | Time |
|------|-------|
| Create API endpoint for budget adjustments | 25 min |
| Write integration tests for all budget scenarios | 30 min |

---

## 🎯 WEEK 6: Frontend Integration (Day 16-18)

### Goal: Connect React frontend to all backend APIs

#### Day 16: Next.js Setup & Auth UI
| Task | Time |
|------|-------|
| Setup Next.js App Router project | 45 min |
| Create authentication pages | 50 min |
| Implement client API calls (TanStack Query) | 30 min |

#### Day 17: Accounts & Transactions Frontend
| Task | Time |
|------|-------|
| Create Accounts page UI (list + create form) | 50 min |
| Implement Transaction list UI with filters | 60 min |
| Create Transaction form (modal/page) | 45 min |

#### Day 18: Dashboard & Budget Frontend
| Task | Time |
|------|-------|
| Implement Dashboard component | 40 min |
| Create Charts using Recharts | 35 min |
| Build Budget management UI | 45 min |

---

## 🎯 WEEK 7: Testing, Polish & Optimization (Day 19-21)

### Goal: Quality assurance + performance tuning

#### Day 19: Unit Tests (Backend)
| Task | Time |
|------|-------|
| Write unit tests for service layer | 60 min |
| Test repository methods with `*testing.T` | 30 min |
| Run all tests locally | 15 min |

#### Day 20: E2E Tests (Playwright)
| Task | Time |
|------|-------|
| Setup Playwright with Next.js | 40 min |
| Write E2E tests for critical flows | 90 min |
| Performance test (Lighthouse) | 15 min |

#### Day 21: Code Quality & Optimization
| Task | Time |
|------|-------|
| Run linter & fix issues | 30 min |
| Format code with gofmt | 10 min |
| Optimize database queries (indexes) | 25 min |
| Review security checklist | 30 min |

---

## 🎯 WEEK 8: Deployment & Launch (Day 22-25)

### Goal: MVP goes live + documentation complete

#### Day 22: Dockerize & CI/CD
| Task | Time |
|------|-------|
| Create Dockerfile for Go backend | 45 min |
| Setup GitHub Actions workflow | 60 min |
| Configure environment variables | 15 min |

#### Day 23: Database Migration & Seeding
| Task | Time |
|------|-------|
| Add default seed data for demo | 30 min |
| Test migration on production DB (separate) | 20 min |

#### Day 24: API Documentation & README
| Task | Time |
|------|-------|
| Document all API endpoints (Swagger/OpenAPI) | 60 min |
| Update README.md with project overview | 30 min |
| Create CONTRIBUTING guide | 15 min |

#### Day 25: Final Testing & Launch Preparation
| Task | Time |
|------|-------|
| Run full regression test suite | 30 min |
| User acceptance testing (internal) | 40 min |
| Deploy MVP to hosting platform | 60 min |
| Document portfolio case study | 30 min |

---

## 📊 Progress Tracking Summary

| Week | Focus Area | Success Metric |
|------|-----------|---------------|
| W1 | Setup + Config | DB connected, Hello World endpoint works |
| W2 | Core Entities | Users can create accounts & categories |
| W3 | Transaction System | Full CRUD for income/expense transactions |
| W4 | Dashboard Insights | Visual analytics showing cash flow |
| W5 | Budget Tracking | Monthly budget limits with status |
| W6 | Frontend Integration | Complete user journey in Next.js UI |
| W7 | Quality & Polish | All tests pass, code formatted, secure |
| W8 | Deployment & Launch | MVP live + documentation complete |

---

## 🎓 Learning Goals (Portfolio Value)

### What You'll Learn & Demonstrate:
1. ✅ **Go Architecture**: Clean architecture with proper separation of concerns
2. ✅ **Database Mastery**: PostgreSQL migrations, indexing, query optimization
3. ✅ **Web Security**: Authentication, authorization, rate limiting, CORS
4. ✅ **Full-Stack Integration**: Backend + frontend communication patterns
5. ✅ **Testing Excellence**: Unit tests, integration tests, E2E tests with Playwright
6. ✅ **DevOps Skills**: Docker, CI/CD pipelines, deployment strategy

---

## 🎯 Definition of Done (MVP Complete)

FinTrack will be considered MVP-ready when ALL these criteria are met:

- [ ] User can register and login
- [ ] User can create at least one financial account
- [ ] User can create, edit, delete transactions
- [ ] Dashboard shows balance, income, expense, net cash flow
- [ ] Users can set monthly budgets per category
- [ ] Budget status updates when transactions are created/updated
- [ ] Security tests confirm data isolation between users
- [ ] Frontend responsive and handles loading/error states
- [ ] Backend has unit/integration tests
- [ ] Migration and seed data available
- [ ] Application deployed and accessible via demo URL
- [ ] README and API documentation clear & complete

---

## 🔑 Key Constraints & Guidelines

### MVP Scope (P0 - Must Have)
1. Authentication & session management
2. Financial accounts (cash, bank, e-wallet)
3. Income & expense transactions
4. Dashboard summary with cash flow visualization
5. Monthly budget tracking per category

### Out of Scope (P1/P2 - Future Enhancement)
- Integration with bank APIs/open banking
- AI-powered financial advice
- Multi-currency support
- Household roles/permissions
- Mobile native applications
- Subscription billing features
- Advanced analytics & reporting

### Technical Constraints
- **Go**: REST API architecture with proper separation (handler/service/repository/model)
- **Database**: PostgreSQL with migrations and test isolation
- **Frontend**: Next.js App Router + TypeScript + Tailwind CSS + shadcn/ui
- **Validation**: Zod on frontend, Go validator library on backend
- **Testing**: TanStack Query for server state, Playwright for E2E

---

*Document generated for FinTrack MVP Planning – 8 Week Development Timeline*
