package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const rateAPIURL = "https://navi.gaitame.com/v3/info/prices/rate"

type rateAPIResponse struct {
	Status int        `json:"status"`
	Data   []rateItem `json:"data"`
}

type rateItem struct {
	Pair  string  `json:"pair"`
	Bid   float64 `json:"bid"`
	Ask   float64 `json:"ask"`
	Open  float64 `json:"open"`
	High  float64 `json:"high"`
	Low   float64 `json:"low"`
	Diff  float64 `json:"diff"`
	Close float64 `json:"close"`
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	outDir := flag.String("outdir", "data/usdjpy", "出力ディレクトリ")
	dateStr := flag.String("date", "", "保存ファイル名に使う日付(YYYY-MM-DD)。未指定ならローカル日付")
	timeout := flag.Duration("timeout", 10*time.Second, "HTTPタイムアウト")
	flag.Parse()

	name := *dateStr
	if name == "" {
		name = time.Now().Format("2006-01-02")
	}

	rate, err := fetchUSDJPY(*timeout)
	if err != nil {
		exitf("fetch failed: %v", err)
	}

	absOutDir, err := filepath.Abs(*outDir)
	if err != nil {
		exitf("failed to resolve absolute path for %q: %v", *outDir, err)
	}

	if err := os.MkdirAll(absOutDir, 0o755); err != nil {
		exitf("mkdir failed for %q: %v", absOutDir, err)
	}

	outPath := filepath.Join(*outDir, name+".json")
	out := struct {
		Date string `json:"date"`
		rateItem
	}{
		Date:     name,
		rateItem: rate,
	}

	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		exitf("marshal failed: %v", err)
	}

	if err := os.WriteFile(outPath, append(b, '\n'), 0o644); err != nil {
		exitf("write failed: %v", err)
	}

	log.Printf("saved USDJPY data to %s", outPath)
	fmt.Println(outPath)
}

func fetchUSDJPY(timeout time.Duration) (rateItem, error) {
	client := &http.Client{Timeout: timeout}
	req, err := http.NewRequest(http.MethodGet, rateAPIURL, nil)
	if err != nil {
		return rateItem{}, err
	}
	req.Header.Set("User-Agent", "fx-data-analysis/1.0")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return rateItem{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return rateItem{}, fmt.Errorf("http status %d: %s", resp.StatusCode, string(body))
	}

	var decoded rateAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return rateItem{}, err
	}
	if decoded.Status != 200 {
		return rateItem{}, fmt.Errorf("api status %d", decoded.Status)
	}

	for _, item := range decoded.Data {
		if item.Pair == "USDJPY" {
			return item, nil
		}
	}
	return rateItem{}, errors.New("USDJPY not found")
}

func exitf(format string, args ...any) {
	log.Printf(format, args...)
	os.Exit(1)
}
