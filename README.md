```markdown
# Diago – Bookie Verification CLI

Welcome to **Diago**, a developer-friendly CLI tool to **generate configs** for multiple bookmakers and **verify selectors** on their pages.

---

## **📁 Folder Structure**

```

diago/
├── cmd/
│    └── root.go               # CLI command
├── config/
│    └── generator.go          # Config generator + overrides
├── fetch/
│    └── fetch.go              # Page fetch + selector verification
├── report/
│    └── report.go             # JSON + Markdown report writer
├── utils/
│    └── utils.go              # Bookie registry & loader
├── bookies/
│    ├── bookie.go             # Common Bookie interface
│    ├── bet365.go
│    ├── onexbet.go
│    ├── betway.go
│    ├── sportpesa.go
│    └── williamhill.go
├── bookies.txt                # Bookie names + base URLs
├── bookies-overrides.yaml     # Optional overrides for credentials/selectors
└── main.go                    # Entry point

````

---

## **⚡ Quick Start**

### 1️⃣ Auto-generate configs and fetch

If you have **no config files yet**, or just want a full run:

```bash
go run main.go \
  --mode=auto \
  --bookies-file=bookies.txt \
  --output-dir=EMC
````

* Checks for missing configs → generates them → fetches & verifies in one go.
* Reports are saved in `EMC/report.json` and `EMC/report.md`.

---

### 2️⃣ Generate configs only

```bash
go run main.go \
  --mode=generate \
  --bookies-file=bookies.txt \
  --output-dir=EMC
```

Creates a folder for each bookie in `EMC/<bookie>/config.yaml`.

---

### 3️⃣ Fetch and verify only

```bash
go run main.go \
  --mode=fetch \
  --bookies-file=bookies.txt \
  --output-dir=EMC
```

* Fetches pages from all bookies listed in `bookies.txt`.
* Uses **existing configs** in `EMC`.
* Checks **selectors** (login fields, buttons, betting options).

---

### 4️⃣ Using overrides

You can override defaults in `bookies-overrides.yaml`:

```yaml
Bet365:
  Username: "myuser"
  Password: "mypassword"
  Selectors:
    Login:
      UsernameInput: "input#login-username"
      PasswordInput: "input#login-password"
      LoginButton: "button#login-submit"
```

Overrides are merged automatically when generating configs.

---

### 5️⃣ Adding new bookies

1. Add a line in `bookies.txt`:

```
NewBookie https://www.newbookie.com
```

2. Create `bookies/newbookie.go`:

```go
package bookies

import "github.com/PuerkitoBio/goquery"

type NewBookie struct {
    url string
}

func (b *NewBookie) Name() string { return "NewBookie" }
func (b *NewBookie) URL() string { return b.url }
func (b *NewBookie) SetURL(u string) { b.url = u }
func (b *NewBookie) Verify(doc *goquery.Document) map[string]string {
    return map[string]string{
        "HomeCheck": "✅",
    }
}

func init() {
    Register(&NewBookie{})
}
```

✅ Works dynamically without touching existing bookies.

---

### 6️⃣ Report Example

After running fetch, you’ll get **JSON** and **Markdown** reports.

**Markdown Summary Example (`report.md`)**:

```markdown
# Verification Report

## 📊 Summary
| Bookie      | URL                         | Status |
|------------|-----------------------------|--------|
| Bet365     | https://www.bet365.com      | ✅     |
| Betway     | https://www.betway.com      | ❌     |

---

## Bet365 (https://www.bet365.com)
- UsernameInput: ✅
- PasswordInput: ✅
- LoginButton: ✅
- SportDropdown: ✅
- DatePicker: ✅
- SearchButton: ✅
- Moneyline: ✅
- Spread: ✅
- Totals: ✅
Overall: ✅ Passed

## Betway (https://www.betway.com)
- UsernameInput: ✅
- PasswordInput: ✅
- LoginButton: ❌
- SportDropdown: ✅
- DatePicker: ✅
- SearchButton: ✅
- Moneyline: ✅
- Spread: ❌
- Totals: ✅
Overall: ❌ Failed
```

**JSON Example (`report.json`)**:

```json
{
  "summary": [
    { "name": "Bet365", "url": "https://www.bet365.com", "all_pass": true },
    { "name": "Betway", "url": "https://www.betway.com", "all_pass": false }
  ],
  "details": [
    {
      "name": "Bet365",
      "url": "https://www.bet365.com",
      "all_pass": true,
      "results": [
        { "label": "UsernameInput", "status": "✅" },
        { "label": "PasswordInput", "status": "✅" }
      ]
    }
  ]
}
```

---

### 7️⃣ Notes

* `generate` mode is **idempotent**: only updates changed configs.
* `fetch` mode requires **configs to exist**, else generate first (or use `auto` mode).
* Reports include **summary + detailed selector checks**.
* Default browser path and base URL can be customized in `config/generator.go`.
* Designed for **CI/CD pipelines**, incremental updates, and scalable bookie management.

---

Made with ❤️ for developers and testing automation.

```

This now gives **new developers a visual reference** of what to expect in both Markdown and JSON reports.  

I can also **add a diagram showing the workflow: generate → fetch → report** to make it even more beginner-friendly. Do you want me to do that?
```

