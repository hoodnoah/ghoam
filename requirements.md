# Gontra
## Scope
- Minimalist, web-based double-entry accounting system
- Event-sourcing, build state at any given point in time by summing transactions

## Requrements
- Constant state of balance; impossible to become oob
- Chart of accounts
- Cash
  - Reconciliation
- Inventory
  - Specific inventory management
  - Generic inventory management
  - Project focus; create a project, add inventory and expense labor hours
- Reporting
  - Trial Balance
  - Balance Sheet
  - Income Statement
  - Project margin, profit/DLH
- Payables, Receivables
- Year-end close
- Tax Reporting
- Self-Employment tax estimation/reserve/payment
- Vendor management
- Customer management
- Fixed asset reporting
  - with depreciation

## Technology Stack
### Database
- sqlite
  - portable, easy to back up
### Backend
- Go
### Frontend
- HTMX?
