package util

import (
	"fmt"
	"math/big"
	"strings"
	"sync"
)

type Currency string

func (c Currency) Symbol() string {
	// Make sure the currency name is uppercase
	currencyString := c.String()
	switch currencyString {
	case "AUD", "CAD", "NZD", "SGD":
		return "$"
	case "ARS", "MXN":
		return "$"
	case "BRL":
		return "R$"
	case "CHF":
		return "CHF"
	case "CNY", "JPY":
		return "¥"
	case "CZK":
		return "Kč"
	case "DKK", "NOK", "SEK":
		return "kr"
	case "EUR":
		return "€"
	case "HKD":
		return "$"
	case "HUF":
		return "Ft"
	case "IDR":
		return "Rp"
	case "ILS":
		return "₪"
	case "INR":
		return "₹"
	case "USD":
		return "$"
	case "GBP":
		return "£"
	case "KRW":
		return "₩"
	case "MYR":
		return "RM"
	case "PLN":
		return "zł"
	case "RUB":
		return "₽"
	case "THB":
		return "฿"
	case "TRY":
		return "₺"
	case "ZAR":
		return "R"
	default:
		mapValue := currencySymbolFromMap(currencyString)
		if mapValue != "" {
			return mapValue
		}
		return currencyString
	}
}

func (c Currency) String() string {
	return strings.ToUpper(string(c))
}

// CurrencyMap maps currency codes to their symbols. The map is global and can be set with SetCurrencyMap.
type CurrencyMap map[string]string

var currencyMap CurrencyMap
var currencyMapMut = sync.RWMutex{}

func (cm CurrencyMap) Symbol(s string) string {
	symbol, ok := cm[s]
	if !ok {
		return s
	}
	return symbol
}

func currencySymbolFromMap(currency string) string {
	currencyMapMut.RLock()
	defer currencyMapMut.RUnlock()
	if currencyMap == nil {
		return ""
	}
	return currencyMap.Symbol(currency)
}

func CurrencySymbol(currency string) string {
	return Currency(currency).Symbol()
}

func FormatMoney(money float64, currency Currency) string {
	return fmt.Sprintf("%.2f %s", money, currency.Symbol())
}

func FormatMoneyBig(money *big.Float, currency Currency) string {
	return money.Text('f', 2) + " " + currency.Symbol()
}

func SetCurrencyMap(m CurrencyMap) {
	currencyMapMut.Lock()
	defer currencyMapMut.Unlock()
	currencyMap = m
}
