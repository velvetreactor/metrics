# Dependencies
- svu: `brew install caarlos0/tap/svu`

# Releasing
1. Pull the latest changes from `main`
2. Run `VERSION=$(svu minor) make release`