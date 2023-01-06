package rakuten

import (
	"encoding/xml"
	"github.com/syahnur197/rakuten/storage"
	"io"
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
