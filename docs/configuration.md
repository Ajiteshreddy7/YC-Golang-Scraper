# Configuration

## `config/scraper_config.json`

```json
{
  "target_platforms": {
    "greenhouse": ["stripe", "airbnb", "coinbase", "databricks"]
  }
}
```

- The strings inside `greenhouse` are company slugs from their Greenhouse boards.
  - Example: `boards.greenhouse.io/stripe` → slug is `stripe`.
- Add as many as you like; the scraper will iterate them.

## Environment variables

| Name       | Default         | Purpose                                   |
|------------|------------------|-------------------------------------------|
| `DB_PATH`  | `data/jobs.db`   | SQLite database file path                  |
| `LOG_LEVEL`| `INFO`           | `DEBUG`, `INFO`, `WARN`, or `ERROR`        |

### Examples

Windows PowerShell:
```powershell
$env:DB_PATH='../data/jobs.db'; $env:LOG_LEVEL='DEBUG'; go run ./cmd/scraper --config ../config/scraper_config.json
```

macOS/Linux Bash:
```bash
DB_PATH=../data/jobs.db LOG_LEVEL=DEBUG go run ./cmd/scraper --config ../config/scraper_config.json
```
