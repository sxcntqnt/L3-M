package utils

import (
    "bufio"
    "fmt"
    "os"
    "strings"

    "github.com/PuerkitoBio/goquery"
)

// ---------------------------
// Bookie Interface
// ---------------------------

// Bookie defines common behavior that each sportsbook must implement
type Bookie interface {
    Name() string
    URL() string
    SetURL(string) // inject URL from config
    Verify(*goquery.Document) map[string]string
}

// ---------------------------
// Registry
// ---------------------------

var registry = map[string]Bookie{}

// Register is called by each bookie's init()
func Register(b Bookie) {
    registry[b.Name()] = b
}

// GetBookie returns a single bookie by name
func GetBookie(name string) (Bookie, bool) {
    b, ok := registry[name]
    return b, ok
}

// AllRegistered returns all bookies regardless of config
func AllRegistered() []Bookie {
    out := []Bookie{}
    for _, b := range registry {
        out = append(out, b)
    }
    return out
}

// ---------------------------
// Enabled Bookies (from bookies.txt)
// ---------------------------

// EnabledBookies parses bookies.txt (name,url), skips blank lines & comments
func EnabledBookies(filename string) ([]Bookie, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to open %s: %w", filename, err)
    }
    defer file.Close()

    var enabled []Bookie
    scanner := bufio.NewScanner(file)

    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "" || strings.HasPrefix(line, "#") {
            continue // skip blank lines & comments
        }

        parts := strings.Split(line, ",")
        if len(parts) < 2 {
            fmt.Printf("⚠️ Invalid line in %s (need 'name,url'): %s\n", filename, line)
            continue
        }

        name := strings.TrimSpace(parts[0])
        url := strings.TrimSpace(parts[1])

        if b, ok := registry[name]; ok {
            b.SetURL(url)        // inject URL from config
            enabled = append(enabled, b)
        } else {
            fmt.Printf("⚠️ Bookie '%s' listed in %s but not registered in code\n", name, filename)
        }
    }

    if err := scanner.Err(); err != nil {
        return nil, fmt.Errorf("error reading %s: %w", filename, err)
    }

    return enabled, nil
}

