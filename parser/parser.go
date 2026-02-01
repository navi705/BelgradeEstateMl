package main

import (
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
)

type WhoCreated int

const (
	Unknown = iota
	Agent
	User
	Investor
)

type RealEstate struct {
	Price               int32
	Currency            string
	PricePerSquareMeter int32
	SquareMeter         int32
	City                string
	District            string
	Municipality        string
	Street              string
	FullLocation        string
	WhoCreated          WhoCreated
	QuantityRoom        float32
	Floor               float32
	FloorTotal          float32
	Link                string
	ParsingDate         time.Time
	Source              string
}

func (w WhoCreated) String() string {
	return [...]string{"Unknown", "Agent", "User", "Investor"}[w]
}

func setupParser() *colly.Collector {
	const userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
	const accept = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"
	const acceptLanguage = "en-US,en;q=0.9,sr;q=0.8,rs;q=0.7"
	const referer = "https://www.google.com/"
	const delay = 2 * time.Second
	const randomDelay = 3 * time.Second
	const parallelism = 1

	parser := colly.NewCollector()
	parser.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Delay:       delay,
		RandomDelay: randomDelay,
		Parallelism: parallelism,
	})

	parser.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", userAgent)
		r.Headers.Set("Accept", accept)
		r.Headers.Set("Accept-Language", acceptLanguage)
		r.Headers.Set("Referer", referer)
		r.Headers.Set("Connection", "keep-alive")
	})

	extensions.Referer(parser)

	slog.Debug("parser initialized")

	return parser
}

type EstateParser func(e *colly.HTMLElement) RealEstate

func parseWebSiteData(domen string, page int, firstPage string, secondPage string, goLangQuery string, callback EstateParser, paginationCallback func(*colly.HTMLElement) int) ([]RealEstate, int, error) {
	parser := setupParser()
	var parsingError error
	var estates []RealEstate
	var totalItems int

	parser.OnResponse(func(r *colly.Response) {
		if r.StatusCode != 200 {
			slog.Error("request failed for", "domen", domen, "status code", r.StatusCode, "url", r.Request.URL)
			parsingError = fmt.Errorf("request failed for %s: %d %s", domen, r.StatusCode, r.Request.URL)
		}
	})

	parser.OnHTML(goLangQuery, func(e *colly.HTMLElement) {
		estate := callback(e)
		if estate.Price > 0 {
			estates = append(estates, estate)
		}
	})

	// Add pagination callback if provided
	if paginationCallback != nil {
		// CityExpert specific selector for now, but could be parameterized
		if domen == "cityexpert.rs" {
			parser.OnHTML(".cx-pagination span", func(e *colly.HTMLElement) {
				count := paginationCallback(e)
				if count > 0 {
					totalItems = count
				}
			})
		}
	}

	if page <= 0 {
		slog.Error("page must be greater than 0", "page", page)
		return nil, 0, errors.New("page must be greater than 0")
	}

	visitError := error(nil)
	if page == 1 {
		visitError = parser.Visit(firstPage)
	} else {
		visitError = parser.Visit(fmt.Sprintf(secondPage, page))
	}

	if visitError != nil {
		slog.Error("visit failed for", "domen", domen, "url", firstPage, "error", visitError)
		return nil, 0, visitError
	}

	if parsingError != nil {
		slog.Error("parsing failed for", "domen", domen, "error", parsingError)
		return nil, 0, parsingError
	}

	count := len(estates)
	slog.Info("parsing successfully finished", "domain", domen, "count", count, "totalItems", totalItems)

	if count > 0 {
		slog.Info("first element", "data", estates[0])
		slog.Info("last element", "data", estates[count-1])
	}

	return estates, totalItems, nil
}

func FourZidaList(page int) ([]RealEstate, int, error) {
	estates, total, err := parseWebSiteData("4zida.rs", page, "https://www.4zida.rs/prodaja-stanova/beograd", "https://www.4zida.rs/prodaja-stanova/beograd?strana=%d", "[test-data='ad-search-card']", parse4ZidaCard, nil)
	return estates, total, err
}

func parse4ZidaCard(e *colly.HTMLElement) RealEstate {
	priceStr := e.ChildText("div.w-3\\/8 p:nth-child(1)")
	location := e.ChildText("p.line-clamp-2")

	estate := RealEstate{
		Source:       "4zida.rs",
		Street:       e.ChildText("p.truncate"),
		FullLocation: location,
		Link:         e.Request.AbsoluteURL(e.ChildAttr("a", "href")),
		Price:        parseNumeric(priceStr),
		Currency:     parseCurrency(priceStr),
		ParsingDate:  time.Now(),
	}

	estate.PricePerSquareMeter = parseNumeric(e.ChildText("div.w-3\\/8 p:nth-child(2)"))

	parseLocationParts4Zida(location, &estate)
	parseDetails(e.ChildText("a.px-3"), &estate)

	metaStr := e.ChildText("div:nth-child(3) div:nth-child(1) span")
	if strings.Contains(metaStr, "Agencija") {
		estate.WhoCreated = Agent
	} else if strings.Contains(metaStr, "Investitor") {
		estate.WhoCreated = Investor
	} else if strings.Contains(metaStr, "Vlasnik") {
		estate.WhoCreated = User
	}

	return estate
}

func HaloOglasiList(page int) ([]RealEstate, int, error) {
	estates, total, err := parseWebSiteData("halooglasi.com", page, "https://www.halooglasi.com/nekretnine/prodaja-stanova/beograd", "https://www.halooglasi.com/nekretnine/prodaja-stanova/beograd?page=%d", ".product-item", parseHaloOglasiCard, nil)
	return estates, total, err
}

func parseHaloOglasiCard(e *colly.HTMLElement) RealEstate {
	priceStr := e.ChildText(".central-feature span[data-value]")
	if priceStr == "" {
		priceStr = e.ChildText(".central-feature i")
	}

	estate := RealEstate{
		Source:      "halooglasi.com",
		Link:        e.Request.AbsoluteURL(e.ChildAttr(".product-title a", "href")),
		Price:       parseNumeric(priceStr),
		Currency:    parseCurrency(priceStr),
		ParsingDate: time.Now(),
	}

	estate.PricePerSquareMeter = parseNumeric(e.ChildText(".price-by-surface span"))

	var locParts []string
	e.ForEach(".subtitle-places li", func(_ int, el *colly.HTMLElement) {
		locParts = append(locParts, strings.TrimSpace(el.Text))
	})
	if len(locParts) > 0 {
		estate.FullLocation = strings.Join(locParts, ", ")
		if len(locParts) >= 2 {
			estate.City = locParts[0]
			estate.Municipality = locParts[1]
		}
		if len(locParts) >= 3 {
			estate.City = locParts[0]
			estate.Municipality = locParts[1]
			estate.District = locParts[2]
		}
		if len(locParts) >= 4 {
			estate.City = locParts[0]
			estate.Municipality = locParts[1]
			estate.District = locParts[2]
			estate.Street = locParts[3]
		}
	}

	e.ForEach(".product-features li", func(_ int, el *colly.HTMLElement) {
		val := el.ChildText(".value-wrapper")
		legend := el.ChildText(".legend")

		switch legend {
		case "Kvadratura":
			estate.SquareMeter = parseNumeric(val)
		case "Broj soba":
			roomsStr := strings.TrimSpace(strings.ReplaceAll(val, "Broj soba", ""))
			if r, err := strconv.ParseFloat(roomsStr, 32); err == nil {
				estate.QuantityRoom = float32(r)
			}
		case "Spratnost":
			floorVal := strings.ReplaceAll(val, "Spratnost", "")
			estate.Floor, estate.FloorTotal = parseFloor(floorVal)
		}
	})

	if strings.Contains(strings.ToLower(e.ChildText(".basic-info")), "agencija") {
		estate.WhoCreated = Agent
	} else if strings.Contains(strings.ToLower(e.ChildText(".basic-info")), "vlasnik") {
		estate.WhoCreated = User
	}

	return estate
}

func NekretnineList(page int) ([]RealEstate, int, error) {
	estates, total, err := parseWebSiteData("nekretnine.rs", page, "https://www.nekretnine.rs/stambeni-objekti/stanovi/izdavanje-prodaja/prodaja/grad/beograd/lista/", "https://www.nekretnine.rs/stambeni-objekti/stanovi/izdavanje-prodaja/prodaja/grad/beograd/lista/stranica/%d/", ".row.offer", parseNekretnineCard, nil)
	return estates, total, err
}

func parseNekretnineCard(e *colly.HTMLElement) RealEstate {
	priceStr := e.ChildText(".offer-price:not(.offer-price--invert) span")
	location := strings.TrimSpace(e.ChildText(".offer-location"))

	estate := RealEstate{
		Source:       "nekretnine.rs",
		FullLocation: location,
		Link:         e.Request.AbsoluteURL(e.ChildAttr("h2.offer-title a", "href")),
		Price:        parseNumeric(priceStr),
		Currency:     parseCurrency(priceStr),
		ParsingDate:  time.Now(),
	}

	estate.PricePerSquareMeter = parseNumeric(e.ChildText(".offer-price:not(.offer-price--invert) small"))
	estate.SquareMeter = parseNumeric(e.ChildText(".offer-price--invert span"))

	parseLocationPartsNekretnine(location, &estate)

	metaStr := strings.ToUpper(e.ChildText(".owner-box"))
	if strings.Contains(metaStr, "AGENCIJA") {
		estate.WhoCreated = Agent
	} else if strings.Contains(metaStr, "VLASNIK") {
		estate.WhoCreated = User
	} else if strings.Contains(metaStr, "INVESTITOR") {
		estate.WhoCreated = Investor
	}

	dateParts := strings.Split(strings.TrimSpace(e.ChildText(".offer-meta-info")), " | ")
	for _, part := range dateParts {
		part = strings.TrimSpace(part)
		if strings.Contains(part, "soban") || strings.Contains(part, "Garsonjera") {
			estate.QuantityRoom = parseSerbianRooms(part)
		}
	}

	return estate
}

func parseSerbianRooms(s string) float32 {
	s = strings.ToUpper(strings.TrimSpace(s))
	// Clean up commonly attached words
	s = strings.ReplaceAll(s, "STAN", "")
	s = strings.TrimSpace(s)

	switch {
	case strings.Contains(s, "GARSONJERA"):
		return 0.5
	case strings.Contains(s, "JEDNOSOBAN"):
		return 1.0
	case strings.Contains(s, "JEDNOIPOSOBAN"):
		return 1.5
	case strings.Contains(s, "DVOSOBAN"):
		return 2.0
	case strings.Contains(s, "DVOIPOSOBAN"):
		return 2.5
	case strings.Contains(s, "TROSOBAN"):
		return 3.0
	case strings.Contains(s, "TROIPOSOBAN"):
		return 3.5
	case strings.Contains(s, "ČETVOROSOBAN"), strings.Contains(s, "CETVOROSOBAN"): // Handle latin/cyrillic charset issues if any
		return 4.0
	case strings.Contains(s, "ČETVOIPOSOBAN"), strings.Contains(s, "CETVOIPOSOBAN"):
		return 4.5
	case strings.Contains(s, "PETOSOBAN"):
		return 5.0
	case strings.Contains(s, "PETIPOSOBAN"):
		return 5.5
	case strings.Contains(s, "ŠESTOSOBAN"), strings.Contains(s, "SESTOSOBAN"):
		return 6.0
	}
	return 0
}

func CityExpertList(page int) ([]RealEstate, int, error) {
	return parseWebSiteData("cityexpert.rs", page, "https://cityexpert.rs/prodaja-nekretnina/beograd?ptId=1", "https://cityexpert.rs/prodaja-nekretnina/beograd?ptId=1&currentPage=%d", ".prop-card", parseCityExpertCard, parseCityExpertPagination)
}

func parseCityExpertPagination(e *colly.HTMLElement) int {
	text := strings.TrimSpace(e.Text)
	return parseCityExpertTotalCount(text)
}

func parseCityExpertTotalCount(text string) int {
	// Example: "571-596 od 596 rezultata"
	parts := strings.Split(text, " od ")
	if len(parts) == 2 {
		totalStr := strings.TrimSpace(strings.Split(parts[1], " ")[0])
		return int(parseNumeric(totalStr))
	}
	return 0
}

func parseCityExpertCard(e *colly.HTMLElement) RealEstate {
	priceStr := e.ChildText(".property-card__price-value span")
	location := strings.TrimSpace(e.ChildText(".property-card__place:not(.property-card__place--break)"))

	estate := RealEstate{
		Source:       "cityexpert.rs",
		FullLocation: location,
		Link:         "https://cityexpert.rs" + e.ChildAttr("a", "href"),
		Price:        parseNumeric(priceStr),
		Currency:     parseCurrency(priceStr),
		ParsingDate:  time.Now(),
	}

	e.ForEach(".property-card__feature", func(_ int, el *colly.HTMLElement) {
		text := strings.TrimSpace(el.Text)
		if strings.Contains(text, "m²") {
			estate.SquareMeter = parseNumeric(text)
			estate.PricePerSquareMeter = estate.Price / estate.SquareMeter
		} else if strings.Contains(text, "Spavaćih soba") {
			numStr := strings.TrimSpace(strings.ReplaceAll(text, "Spavaćih soba", ""))
			if rooms, err := strconv.ParseFloat(numStr, 32); err == nil {
				estate.QuantityRoom = float32(rooms)
			}
		}
	})

	parseLocationPartsCE(location, &estate)

	return estate
}

func parseLocationPartsNekretnine(location string, estate *RealEstate) {
	parts := strings.Split(location, ",")
	if len(parts) >= 2 {
		estate.District = strings.TrimSpace(parts[0])
		estate.City = strings.TrimSpace(parts[1])
	}
}

func parseLocationPartsCE(location string, estate *RealEstate) {
	parts := strings.Split(location, ", ")
	if len(parts) >= 2 {
		estate.Street = parts[0]
		estate.Municipality = parts[1]
		estate.City = "Beograd"
	}
}

func parseNumeric(s string) int32 {
	start := -1
	for i, r := range s {
		if r >= '0' && r <= '9' {
			start = i
			break
		}
	}
	if start == -1 {
		return 0
	}

	end := start
	for i := start; i < len(s); i++ {
		r := rune(s[i])
		if (r >= '0' && r <= '9') || r == '.' || r == ',' || r == ' ' {
			end = i + 1
		} else {
			break
		}
	}

	numStr := s[start:end]

	numStr = strings.ReplaceAll(numStr, " ", "")

	numStr = strings.ReplaceAll(numStr, ".", "")
	numStr = strings.ReplaceAll(numStr, ",", ".")

	val, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0
	}

	return int32(val + 0.5)
}

func parseCurrency(s string) string {
	if strings.Contains(s, "€") || strings.Contains(s, "EUR") {
		return "EUR"
	} else if strings.Contains(s, "RSD") {
		return "RSD"
	}
	return ""
}

func parseLocationParts4Zida(location string, estate *RealEstate) {
	parts := strings.Split(location, ", ")
	if len(parts) >= 3 {
		estate.District = parts[0]
		estate.Municipality = parts[1]
		estate.City = parts[2]
	} else if len(parts) >= 2 {
		estate.Municipality = parts[0]
		estate.City = parts[1]
	}
}

func parseDetails(details string, estate *RealEstate) {
	parts := strings.Split(details, " | ")
	for _, part := range parts {
		if strings.Contains(part, "m²") {
			m2Str := strings.TrimSuffix(part, "m²")
			estate.SquareMeter = parseNumeric(m2Str)
		} else if strings.Contains(part, "sobe") || strings.Contains(part, "soba") {
			roomsStr := strings.Fields(part)[0]
			if r, err := strconv.ParseFloat(roomsStr, 32); err == nil {
				estate.QuantityRoom = float32(r)
			}
		} else if strings.Contains(part, "sprat") {
			floorStr := strings.Fields(part)[0]
			estate.Floor, estate.FloorTotal = parseFloor(floorStr)
		}
	}
}

func parseFloor(s string) (float32, float32) {
	var floorVal, totalVal float32 = -5, -5
	parts := strings.Split(s, "/")

	floorVal = parseSingleFloorPart(parts[0])
	if len(parts) > 1 {
		totalVal = parseSingleFloorPart(parts[1])
	}

	return floorVal, totalVal
}

func parseSingleFloorPart(s string) float32 {
	s = strings.ToUpper(strings.TrimSpace(s))
	if s == "" {
		return -5
	}

	cleanS := strings.ReplaceAll(s, ".", "")
	cleanS = strings.ReplaceAll(cleanS, ",", "")

	switch cleanS {
	case "SUT", "SU", "SUTEREN", "PODRUM":
		return -3.0
	case "PSUT":
		return -2.0
	case "NPR", "NISKOPRIZEMLJE", "NISKO PRIZEMLJE":
		return -0.5
	case "PR", "PRIZEMLJE":
		return 0.0
	case "VPR", "VISOKOPRIZEMLJE", "VISOKO PRIZEMLJE":
		return 0.5
	case "PTK", "POTKROVLJE":
		return 1000.0
	}

	numStr := strings.ReplaceAll(s, ".", "")
	numStr = strings.ReplaceAll(numStr, ",", ".")

	if val, err := strconv.ParseFloat(numStr, 32); err == nil {
		return float32(val)
	}

	if val := romanToInt(cleanS); val > 0 {
		return float32(val)
	}

	return -100
}

func romanToInt(s string) int {
	romanMap := map[byte]int{
		'I': 1,
		'V': 5,
		'X': 10,
		'L': 50,
		'C': 100,
		'D': 500,
		'M': 1000,
	}

	result := 0
	length := len(s)

	for i := 0; i < length; i++ {
		val := romanMap[s[i]]
		if val == 0 {
			return -100
		}

		if i+1 < length && romanMap[s[i+1]] > val {
			result -= val
		} else {
			result += val
		}
	}

	return result
}
