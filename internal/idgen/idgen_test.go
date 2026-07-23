package idgen

import (
	"errors"
	"math/big"
	"math/rand"
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

func TestTimeOfRoundTripRandomizedPropertySweep(t *testing.T) {
	// R-WH5F-QJYS
	rng := rand.New(rand.NewSource(0x5eed))
	prefixes := []string{"R", "S", "SPEC", "item42"}

	for i := 0; i < 2_000; i++ {
		ms := rng.Int63n(modulus.Int64())
		want := Epoch.Add(time.Duration(ms) * time.Millisecond)
		for _, prefix := range prefixes {
			id := MintAt(prefix, want)
			got, err := TimeOf(id)
			if err != nil {
				t.Fatalf("TimeOf(MintAt(%q, Epoch+%dms)) returned error: %v", prefix, ms, err)
			}
			if !got.Equal(want) {
				t.Fatalf("TimeOf(MintAt(%q, Epoch+%dms)) = %s, want %s", prefix, ms, got, want)
			}
		}
	}
}

func TestTimeOfIgnoresPrefix(t *testing.T) {
	// R-WN8X-NEO9
	want := Epoch.Add(9_876_543_210 * time.Millisecond)
	minted := MintAt("R", want)
	body := strings.TrimPrefix(minted, "R-")

	for _, prefix := range []string{"R", "S", "SPEC"} {
		got, err := TimeOf(prefix + "-" + body)
		if err != nil {
			t.Fatalf("TimeOf with prefix %q returned error: %v", prefix, err)
		}
		if !got.Equal(want) {
			t.Fatalf("TimeOf with prefix %q = %s, want %s", prefix, got, want)
		}
	}
}

func TestTimeOfRejectsMalformedIDs(t *testing.T) {
	// R-WPOQ-EY5N
	malformed := []string{
		"",
		"R",
		"R-0007J3LA",
		"R-0007-J3L",
		"R-0007-J3LAA",
		"R-000-7J3LA",
		"R--0007-J3LA",
		"-0007-J3LA",
		"R-X-0007-J3LA",
		"R_1-0007-J3LA",
		"R-0007-j3LA",
		"R-0007-J3L!",
		"R-０００７-J3LA",
		" R-0007-J3LA",
		"R-0007-J3LA\n",
	}

	for _, id := range malformed {
		_, err := TimeOf(id)
		if !errors.Is(err, ErrInvalidID) {
			t.Errorf("TimeOf(%q) error = %v, want error wrapping ErrInvalidID", id, err)
		}
	}
}

func TestTimeOfRobustnessSweep(t *testing.T) {
	// R-WQWM-SPWC
	rng := rand.New(rand.NewSource(0xc0ffee))
	inputs := make([]string, 0, 4_000)

	for i := 0; i < 2_000; i++ {
		raw := make([]byte, rng.Intn(65))
		if _, err := rng.Read(raw); err != nil {
			t.Fatalf("generate random input: %v", err)
		}
		inputs = append(inputs, string(raw))
	}

	alphabet := []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	for i := 0; i < 2_000; i++ {
		prefixLength := 1 + rng.Intn(12)
		bodyLength := 6 + rng.Intn(6)
		raw := make([]byte, prefixLength+1+bodyLength)
		for j := 0; j < prefixLength; j++ {
			raw[j] = alphabet[rng.Intn(len(alphabet))]
		}
		raw[prefixLength] = '-'
		for j := prefixLength + 1; j < len(raw); j++ {
			raw[j] = alphabet[rng.Intn(len(alphabet))]
		}
		if len(raw) > prefixLength+5 {
			raw[prefixLength+5] = []byte{'-', '_', '/', 'a'}[rng.Intn(4)]
		}
		inputs = append(inputs, string(raw))
	}

	for _, id := range inputs {
		got, err := TimeOf(id)
		if err != nil {
			if !errors.Is(err, ErrInvalidID) {
				t.Fatalf("TimeOf(%q) error = %v, want error wrapping ErrInvalidID", id, err)
			}
			continue
		}
		if got.Before(Epoch) || !got.Before(Epoch.Add(time.Duration(modulus.Int64())*time.Millisecond)) {
			t.Fatalf("TimeOf(%q) = %s, outside representable range", id, got)
		}
	}
}

func TestMustModInversePanicsWhenValuesAreNotCoprime(t *testing.T) {
	// R-WS4J-6HN1
	defer func() {
		if recover() == nil {
			t.Fatal("mustModInverse did not panic for non-coprime values")
		}
	}()

	mustModInverse(big.NewInt(6), big.NewInt(36))
}
