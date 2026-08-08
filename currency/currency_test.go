package util

import (
	"math/big"
	"testing"
)

func TestCurrencyStringAndSymbol(t *testing.T) {
	tests := []struct {
		currency   Currency
		wantName   string
		wantSymbol string
	}{
		{Currency("eur"), "EUR", "€"},
		{Currency("usd"), "USD", "$"},
		{Currency("gbp"), "GBP", "£"},
		{Currency("jpy"), "JPY", "¥"},
		{Currency("aud"), "AUD", "$"},
		{Currency("brl"), "BRL", "R$"},
		{Currency("chf"), "CHF", "CHF"},
		{Currency("cny"), "CNY", "¥"},
		{Currency("czk"), "CZK", "Kč"},
		{Currency("dkk"), "DKK", "kr"},
		{Currency("huf"), "HUF", "Ft"},
		{Currency("ils"), "ILS", "₪"},
		{Currency("inr"), "INR", "₹"},
		{Currency("krw"), "KRW", "₩"},
		{Currency("pln"), "PLN", "zł"},
		{Currency("rub"), "RUB", "₽"},
		{Currency("thb"), "THB", "฿"},
		{Currency("try"), "TRY", "₺"},
		{Currency("zar"), "ZAR", "R"},
	}
	for _, test := range tests {
		if got := test.currency.String(); got != test.wantName {
			t.Errorf("%q.String() = %q, want %q", test.currency, got, test.wantName)
		}
		if got := test.currency.Symbol(); got != test.wantSymbol {
			t.Errorf("%q.Symbol() = %q, want %q", test.currency, got, test.wantSymbol)
		}
	}
	if got := CurrencySymbol("eur"); got != "€" {
		t.Errorf("CurrencySymbol(\"eur\") = %q, want %q", got, "€")
	}
}

func TestCurrencyMapAndFormatting(t *testing.T) {
	previous := currencyMap
	t.Cleanup(func() { SetCurrencyMap(previous) })

	SetCurrencyMap(CurrencyMap{"XTS": "¤"})
	if got := Currency("xts").Symbol(); got != "¤" {
		t.Errorf("custom XTS symbol = %q, want %q", got, "¤")
	}
	if got := (CurrencyMap{"USD": "$"}).Symbol("missing"); got != "missing" {
		t.Errorf("missing CurrencyMap value = %q, want %q", got, "missing")
	}
	if got := FormatMoney(12.345, "jpy"); got != "12.35 ¥" {
		t.Errorf("FormatMoney() = %q, want %q", got, "12.35 ¥")
	}
	if got := FormatMoneyBig(big.NewFloat(12.345), "jpy"); got != "12.35 ¥" {
		t.Errorf("FormatMoneyBig() = %q, want %q", got, "12.35 ¥")
	}
}
