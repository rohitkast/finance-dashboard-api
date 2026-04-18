# Personal Finance Dashboard API

## How I used AI for this

---
**Most of this README**. But, this is me. Any personal messages, assumptions, tradeoffs I took is written by me. Only the code explanation part is claude because I think it is better at it then me.

**Error logging throughout**. 50 percent done with the code, I realized I haven't put loggers anywhere. So I asked AI to do it for me.

**Spelling mistakes and refactoring**. Clearly its better at it. By refactoring I mean, towards the end when I needed to make changes that would affect multiple files, I would explain everything to it. 

---

## What This Project Does

This is a REST API for a personal finance dashboard where users can:
- Track their income and expenses
- Filter and analyze their transactions
- Get financial summaries and trends

Plus, admins can manage users and view all transactions across the platform.

---

## Role-Based Access (Important Assumption!)

The requirements mentioned three roles but didn't specify exactly what each role could do. Here's what I assumed and implemented:

**Viewer** (Read-Only User)
- Can view their own dashboard and transaction data
- Get summaries, trends, and analytics
- **Cannot** create, edit, or delete any transactions
- Use case: Maybe a financial advisor viewing client data, or a family member with view-only access

**Analyst** (Regular Active User)
- Everything a Viewer can do, PLUS:
- **Can create** new transactions
- **Can edit** their own transactions  
- **Can delete** their own transactions
- Use case: This is your typical user managing their day-to-day finances

**Admin** (Platform Manager)
- Everything an Analyst can do, PLUS:
- Manage all users (view, delete)
- View ANY user's transactions
- Delete ANY transaction
- Use case: Platform administrator or household head with full control

**Why I chose this:** The task description wasn't clear if there was a "regular user" who could create transactions, or if only admins could. I assumed Analyst = regular active user, because it doesn't make sense to have a finance app where you can't track your own expenses!

---

## Quick Setup (5 minutes)

### Prerequisites
- **Go 1.25+** installed
- **Docker & Docker Compose** installed
- **PostgreSQL client** (optional, for manual DB access)

### Step 1: Start the Database
```bash
docker compose up -d
```
This spins up PostgreSQL on port **5433** (to avoid conflicts with your local setup).

### Step 2: Configure Environment
Create a `.env` file in the root directory:
```env
DB_URL=postgres://postgres:postgres@localhost:5433/finance?sslmode=disable
PORT=8080
JWT_SECRET=your-super-secret-key-change-this-in-production
```

### Step 3: Run the Application
```bash
go mod download
go run main.go
```

The API will start on `http://localhost:8080`

---

## Testing the API (Recommended: Use Bruno!)

I've included a **Bruno collection** in this repo that has all the endpoints pre-configured. Just:
1. Install [Bruno](https://www.usebruno.com/)
2. Open the collection from this project
3. Start making requests!

It's way easier than manually testing with curl.

---

## How to Use (The Flow)

### Step 1: Create a User
First, create an account. There are three user roles:
- **viewer** - Read-only access to your own data
- **analyst** - Full access to manage your own transactions (this is your typical user!)
- **admin** - Platform manager with access to all users and data

```bash
POST /api/auth/create
{
  "email": "john@example.com",
  "password": "password123",
  "role": "analyst"
}
```

### Step 2: Login
Get your access token:
```bash
POST /api/auth/login
{
  "email": "john@example.com",
  "password": "password123"
}
```
Returns: `{"access_token": "your-jwt-token"}`

### Step 3: Use the Token
Add this to all subsequent requests:
```
Authorization: Bearer your-jwt-token
```

### Step 4: Start Managing Transactions!
Now you can create transactions, filter them, get summaries, etc.

---

## API Endpoints

### User Management
| Method | Endpoint | Description | Roles Allowed |
|--------|----------|-------------|---------------|
| POST | `/api/auth/create` | Register a new user | Public |
| POST | `/api/auth/login` | Login and get JWT token | Public |
| GET | `/api/admin/users` | Get all users | Admin only |
| GET | `/api/admin/getuser/:id` | Get user by ID | Admin only |
| DELETE | `/api/admin/user/:id` | Delete user by ID | Admin only |

### Transaction Management
| Method | Endpoint | Description | Roles Allowed |
|--------|----------|-------------|---------------|
| POST | `/api/transaction/create` | Create a new transaction | Analyst, Admin |
| GET | `/api/transaction/all` | Get all your transactions | Viewer, Analyst, Admin |
| GET | `/api/transaction/:id` | Get specific transaction | Viewer, Analyst, Admin |
| PUT | `/api/transaction/:id` | Update your transaction | Analyst, Admin |
| DELETE | `/api/transaction/:id` | Delete your transaction | Analyst, Admin |
| GET | `/api/transaction/filter` | Get filtered transactions | Viewer, Analyst, Admin |

**Filter Query Parameters:**
- `asc=true/false` - Sort by amount ascending
- `dsc=true/false` - Sort by amount descending
- `month=1-12` - Filter by month
- `year=2024` - Filter by year
- `from=2024-01-01` - Start date (YYYY-MM-DD)
- `to=2024-12-31` - End date (YYYY-MM-DD)
- `exp=true` - Only expenses
- `inc=true` - Only income

### Dashboard & Analytics
| Method | Endpoint | Description | Roles Allowed |
|--------|----------|-------------|---------------|
| GET | `/api/dashboard/summary` | Get financial summary | Viewer, Analyst, Admin |
| GET | `/api/dashboard/recent` | Get recent transactions | Viewer, Analyst, Admin |
| GET | `/api/dashboard/category/:category` | Transactions by category (income/expense) | Viewer, Analyst, Admin |
| GET | `/api/dashboard/trends/:category` | Monthly trends by category | Viewer, Analyst, Admin |

### Admin Routes
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/api/admin/users` | Get all users | Admin only |
| GET | `/api/admin/getuser/:id` | Get specific user | Admin only |
| DELETE | `/api/admin/user/:id` | Delete any user | Admin only |
| GET | `/api/admin/transactions/:uid` | Get user's transactions | Admin only |
| DELETE | `/api/admin/transaction/:id` | Delete any transaction | Admin only |

---

## Design Choices & Tradeoffs

### What I Prioritized
- **Security First**: JWT authentication, role-based access, proper data ownership checks
- **Clean Code**: Used helper functions to avoid repetition (like `GetUserIDFromContext`)
- **Practical Structure**: Separated concerns (api, repository, middleware, models) without over-complicating
- **Database overload**: The database queries deliberately use context timeouts (5 seconds) to prevent hanging connections. Error handling is consistent throughout - log for debugging, return user-friendly messages.

### Tradeoffs
- Didnt use seperate packages as I just needed to build a basic API. It would have been better if i had used different packages for say, admin or dashboard etc.
- Used basic gin logger instead of a third party library or a middleware. But a seperate logger would have been nice giving us details logs.


### What I'd Add With More Time
- Pagination for large transaction lists
- Export to CSV/PDF
- Recurring transactions
- Budget limits and alerts
- Transaction categories with icons/colors
- Better test coverage (unit + integration tests)

---

## Project Structure
```
├── api/              # HTTP handlers (controllers)
├── internal/
│   ├── middleware/   # Auth and role validation
│   ├── models/       # Database models
│   └── repository/   # Database operations
├── routes/           # Route definitions
├── config/           # Configuration loader
├── database/         # DB connection
├── utils/            # Helper functions
├── compose.yaml      # Docker setup
└── main.go           # Entry point
```

---

## Why I Built It This Way

I organized the code so anyone jumping in can quickly understand what's where. Each package has a clear purpose:
- **API handlers** deal with HTTP stuff
- **Repository** handles database logic
- **Middleware** manages security
- **Utils** keeps shared code DRY

---

## Final Thoughts

If you have questions about any design decisions or want to see something different, I'm happy to discuss!

Thanks for reviewing my work
