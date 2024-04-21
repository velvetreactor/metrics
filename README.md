# Dependencies
- svu: `brew install caarlos0/tap/svu`
- golang-migrate: `brew install golang-migrate`

# Development
## Database migrations
- creating a migration: `make migrate-create NAME=[migration name]`
- running migrations: `make migrate`

# Releasing
1. Pull the latest changes from `main`
2. Run `VERSION=$(svu minor) make release`