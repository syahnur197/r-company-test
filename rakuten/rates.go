package rakuten

import (
	"encoding/xml"
	"github.com/syahnur197/rakuten/storage"
	"io"
)

type Currency string

const (
	USD Currency = "USD"
	JPY Currency = "JPY"
	BGN Currency = "BGN"
	CZK Currency = "CZK"
	DKK Currency = "DKK"
	GBP Currency = "GBP"
	HUF Currency = "HUF"
	PLN Currency = "PLN"
	RON Currency = "RON"
	SEK Currency = "SEK"
	CHF Currency = "CHF"
	ISK Currency = "ISK"
	NOK Currency = "NOK"
	TRY Currency = "TRY"
	AUD Currency = "AUD"
	BRL Currency = "BRL"
	CAD Currency = "CAD"
	CNY Currency = "CNY"
	HKD Currency = "HKD"
	IDR Currency = "IDR"
	ILS Currency = "ILS"
	INR Currency = "INR"
	KRW Currency = "KRW"
	MXN Currency = "MXN"
	MYR Currency = "MYR"
	NZD Currency = "NZD"
	PHP Currency = "PHP"
	SGD Currency = "SGD"
	THB Currency = "THB"
	ZAR Currency = "ZAR"
)

type Rate struct {
	Base  string `db:"base"`
	Quote string `xml:"currency,attr" db:"quote"`
	Rate  string `xml:"rate,attr" db:"rate"`
	Date  string `xml:"time,attr" db:"published_date"`
}

type Rates struct {
	Rates RateList `xml:"Cube>Cube"`
}

type RateList []Rate

func (ls *RateList) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	date := start.Attr[0].Value

	for {
		tok, err := d.Token()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		if se, ok := tok.(xml.StartElement); ok {
			rate := Rate{Date: date}
			if err := d.DecodeElement(&rate, &se); err != nil {
				return err
			}

			*ls = append(*ls, rate)
		}
	}
}

func ConvertToStoreRate(rate Rate) storage.Rate {
	return storage.Rate{
		Base:  rate.Base,
		Quote: rate.Quote,
		Rate:  rate.Rate,
		Date:  rate.Date,
	}
}
