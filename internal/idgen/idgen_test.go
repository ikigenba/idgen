package idgen

import (
	"strings"
	"testing"
	"time"
)

func TestMintAtEpochGoldenVector(t *testing.T) {
	// R-WIDC-4BPH
	const want = "R-0007-J3LA"
	if got := MintAt("R", Epoch); got != want {
		t.Fatalf("MintAt(%q, Epoch) = %q, want %q", "R", got, want)
	}
}

func TestMintAtFixedAbsoluteInstantGoldenVector(t *testing.T) {
	// R-WJL8-I3G6
	instant := time.Date(2026, time.January, 2, 3, 4, 5, 678_000_000, time.UTC)
	const want = "R-HG3W-VAV0"
	if got := MintAt("R", instant); got != want {
		t.Fatalf("MintAt(%q, %s) = %q, want %q", "R", instant, got, want)
	}
}

func TestMintAtPadsSmallMillisecondValues(t *testing.T) {
	// R-WKT4-VV6V
	got := MintAt("item", Epoch.Add(time.Millisecond))
	const want = "item-0183-WVBZ"
	if got != want {
		t.Fatalf("MintAt one millisecond after Epoch = %q, want %q", got, want)
	}
	body := strings.ReplaceAll(strings.TrimPrefix(got, "item-"), "-", "")
	if len(body) != 8 {
		t.Fatalf("encoded body length = %d, want 8 (%q)", len(body), body)
	}
}

func TestMintAtClampsInstantsBeforeEpoch(t *testing.T) {
	// R-WM11-9MXK
	got := MintAt("R", Epoch.Add(-24*time.Hour))
	want := MintAt("R", Epoch)
	if got != want {
		t.Fatalf("MintAt before Epoch = %q, want Epoch encoding %q", got, want)
	}
	if got != "R-0007-J3LA" {
		t.Fatalf("clamped encoding = %q, want golden Epoch encoding", got)
	}
}
