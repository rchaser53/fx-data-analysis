package handler

import (
	"math"
	"testing"

	"github.com/rchaser53/fx-data-analysis/internal/model"
)

func TestAggregateUSDJPYRatesByWeek(t *testing.T) {
	rates := []model.USDJPYRate{
		{Date: "2026-03-16", Pair: "USDJPY", Open: 148.0, High: 149.2, Low: 147.8, Close: 148.9, Bid: 148.85, Ask: 148.95},
		{Date: "2026-03-17", Pair: "USDJPY", Open: 148.9, High: 149.5, Low: 148.4, Close: 149.1, Bid: 149.05, Ask: 149.15},
		{Date: "2026-03-18", Pair: "USDJPY", Open: 149.1, High: 150.1, Low: 148.7, Close: 149.8, Bid: 149.75, Ask: 149.85},
		{Date: "2026-03-23", Pair: "USDJPY", Open: 149.7, High: 150.0, Low: 149.0, Close: 149.2, Bid: 149.15, Ask: 149.25},
		{Date: "2026-03-24", Pair: "USDJPY", Open: 149.2, High: 151.0, Low: 149.1, Close: 150.6, Bid: 150.55, Ask: 150.65},
	}

	aggregated, err := aggregateUSDJPYRatesByWeek(rates)
	if err != nil {
		t.Fatalf("aggregateUSDJPYRatesByWeek returned error: %v", err)
	}

	if len(aggregated) != 2 {
		t.Fatalf("expected 2 weekly rates, got %d", len(aggregated))
	}

	first := aggregated[0]
	if first.Date != "2026-03-16" {
		t.Fatalf("expected first week start 2026-03-16, got %s", first.Date)
	}
	if first.Label != "2026-03-16 - 2026-03-18" {
		t.Fatalf("expected first week label to cover grouped range, got %s", first.Label)
	}
	if first.Open != 148.0 || first.Close != 149.8 {
		t.Fatalf("unexpected first week open/close: %+v", first)
	}
	if first.High != 150.1 || first.Low != 147.8 {
		t.Fatalf("unexpected first week high/low: %+v", first)
	}
	if first.Bid != 149.75 || first.Ask != 149.85 {
		t.Fatalf("unexpected first week bid/ask: %+v", first)
	}
	if math.Abs(first.Diff-1.8) > 1e-9 {
		t.Fatalf("expected first week diff 1.8, got %v", first.Diff)
	}

	second := aggregated[1]
	if second.Date != "2026-03-23" {
		t.Fatalf("expected second week start 2026-03-23, got %s", second.Date)
	}
	if second.Label != "2026-03-23 - 2026-03-24" {
		t.Fatalf("expected second week label to cover grouped range, got %s", second.Label)
	}
	if second.Open != 149.7 || second.Close != 150.6 {
		t.Fatalf("unexpected second week open/close: %+v", second)
	}
	if second.High != 151.0 || second.Low != 149.0 {
		t.Fatalf("unexpected second week high/low: %+v", second)
	}
}

func TestNormalizeUSDJPYTimeframe(t *testing.T) {
	tests := map[string]string{
		"":         usdJPYTimeframeDaily,
		"daily":    usdJPYTimeframeDaily,
		"weekly":   usdJPYTimeframeWeekly,
		" WEEKLY ": usdJPYTimeframeWeekly,
		"monthly":  "",
	}

	for input, want := range tests {
		got := normalizeUSDJPYTimeframe(input)
		if got != want {
			t.Fatalf("normalizeUSDJPYTimeframe(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestFilterUSDJPYTradingDays(t *testing.T) {
	rates := []model.USDJPYRate{
		{Date: "2026-03-20", Pair: "USDJPY"},
		{Date: "2026-03-21", Pair: "USDJPY"},
		{Date: "2026-03-22", Pair: "USDJPY"},
		{Date: "2026-03-23", Pair: "USDJPY"},
	}

	filtered, err := filterUSDJPYTradingDays(rates)
	if err != nil {
		t.Fatalf("filterUSDJPYTradingDays returned error: %v", err)
	}

	if len(filtered) != 2 {
		t.Fatalf("expected 2 trading-day rates, got %d", len(filtered))
	}

	if filtered[0].Date != "2026-03-20" {
		t.Fatalf("expected first trading day to remain 2026-03-20, got %s", filtered[0].Date)
	}
	if filtered[1].Date != "2026-03-23" {
		t.Fatalf("expected second trading day to remain 2026-03-23, got %s", filtered[1].Date)
	}
}
