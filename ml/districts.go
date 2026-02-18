package main

import (
	"sort"
	"strings"
)

var districtMapping = map[string]string{
	"Beograd na vodi":           "Savski venac",
	"Savamala":                  "Savski venac",
	"Dedinje":                   "Savski venac",
	"Senjak":                    "Savski venac",
	"Palata pravde":             "Savski venac",
	"Klinički centar":           "Savski venac",
	"Savski trg":                "Savski venac",
	"Mostarska petlja":          "Savski venac",
	"Karaburma":                 "Palilula",
	"Višnjička banja":           "Palilula",
	"Višnjička Banja":           "Palilula",
	"Borča":                     "Palilula",
	"Borča I":                   "Palilula",
	"Borča II":                  "Palilula",
	"Borča III":                 "Palilula",
	"Borča IV":                  "Palilula",
	"Krnjača":                   "Palilula",
	"Kotež":                     "Palilula",
	"Ovča":                      "Palilula",
	"Padinska skela":            "Palilula",
	"Bogoslovija":               "Palilula",
	"Profesorska kolonija":      "Palilula",
	"Hadžipopovac":              "Palilula",
	"Kaluđerica":                "Grocka",
	"Leštane":                   "Grocka",
	"Boleč":                     "Grocka",
	"Vinča":                     "Grocka",
	"Grocka":                    "Grocka",
	"Mirijevo":                  "Zvezdara",
	"Mirijevo I":                "Zvezdara",
	"Mirijevo II":               "Zvezdara",
	"Mirijevo III":              "Zvezdara",
	"Mirijevo IV":               "Zvezdara",
	"Lion":                      "Zvezdara",
	"Konjarnik":                 "Zvezdara",
	"Đeram":                     "Zvezdara",
	"Đeram pijaca":              "Zvezdara",
	"Cvetkova pijaca":           "Zvezdara",
	"Cvetkova Pijaca":           "Zvezdara",
	"Učiteljsko naselje":        "Zvezdara",
	"Učiteljsko Naselje":        "Zvezdara",
	"Bulbulder":                 "Zvezdara",
	"Bulbuder":                  "Zvezdara",
	"Severni bulevar":           "Zvezdara",
	"Severni Bulevar":           "Zvezdara",
	"Olimp":                     "Zvezdara",
	"Denkova bašta":             "Zvezdara",
	"Denkova Bašta":             "Zvezdara",
	"Veliki mokri lug":          "Zvezdara",
	"Veliki Mokri Lug":          "Zvezdara",
	"Mali mokri lug":            "Zvezdara",
	"Mali Mokri Lug":            "Zvezdara",
	"Zvezdarska šuma":           "Zvezdara",
	"Zvezdarska Šuma":           "Zvezdara",
	"Zeleno brdo":               "Zvezdara",
	"Zeleno Brdo":               "Zvezdara",
	"Lipov lad":                 "Zvezdara",
	"Lipov Lad":                 "Zvezdara",
	"Banovo brdo":               "Čukarica",
	"Banovo Brdo":               "Čukarica",
	"Žarkovo":                   "Čukarica",
	"Cerak":                     "Čukarica",
	"Cerak vinogradi":           "Čukarica",
	"Cerak Vinogradi":           "Čukarica",
	"Vidikovac":                 "Čukarica",
	"Vidikovac (Centar)":        "Čukarica",
	"Vidikovačka padina":        "Čukarica",
	"Vidikovačka Padina":        "Čukarica",
	"Vidikovački venac":         "Čukarica",
	"Bele vode":                 "Čukarica",
	"Bele Vode":                 "Čukarica",
	"Železnik":                  "Čukarica",
	"Sremčica":                  "Čukarica",
	"Umka":                      "Čukarica",
	"Ostružnica":                "Čukarica",
	"Julino brdo":               "Čukarica",
	"Julino Brdo":               "Čukarica",
	"Čukarička padina":          "Čukarica",
	"Čukarička Padina":          "Čukarica",
	"Golf naselje":              "Čukarica",
	"Golf Naselje":              "Čukarica",
	"Filmski grad":              "Čukarica",
	"Filmski Grad":              "Čukarica",
	"Košutnjak":                 "Čukarica",
	"Stari Košutnjak":           "Čukarica",
	"Braće Jerković":            "Voždovac",
	"Medaković":                 "Voždovac",
	"Medaković 1":               "Voždovac",
	"Medaković 2":               "Voždovac",
	"Medaković 3":               "Voždovac",
	"Medaković I":               "Voždovac",
	"Medaković II":              "Voždovac",
	"Medaković III":             "Voždovac",
	"Banjica":                   "Voždovac",
	"Autokomanda":               "Voždovac",
	"Trošarina":                 "Voždovac",
	"Dušanovac":                 "Voždovac",
	"Stepa Stepanović":          "Voždovac",
	"Kumodraž":                  "Voždovac",
	"Kumodraž I":                "Voždovac",
	"Kumodraž II":               "Voždovac",
	"Lekino brdo":               "Voždovac",
	"Lekino Brdo":               "Voždovac",
	"Jajinci":                   "Voždovac",
	"Vojvode Stepe":             "Voždovac",
	"Šumice":                    "Voždovac",
	"Voždovačka crkva":          "Voždovac",
	"Saobraćajni fakultet":      "Voždovac",
	"Altina":                    "Zemun",
	"Batajnica":                 "Zemun",
	"Zemun Polje":               "Zemun",
	"Zemun polje":               "Zemun",
	"Meandri":                   "Zemun",
	"Galenika":                  "Zemun",
	"Gornji grad":               "Zemun",
	"Gornji Grad":               "Zemun",
	"Cara Dušana":               "Zemun",
	"Kalvarija":                 "Zemun",
	"Zemunske kapije":           "Zemun",
	"Donji grad":                "Zemun",
	"Donji Grad":                "Zemun",
	"Retenzija":                 "Zemun",
	"Novi grad":                 "Zemun",
	"Novi Grad":                 "Zemun",
	"Zemun (Kej)":               "Zemun",
	"Zemunski Kej":              "Zemun",
	"Zemun (Marije Bursać)":     "Zemun",
	"Kalenić pijaca":            "Vračar",
	"Kalenić":                   "Vračar",
	"Čubura":                    "Vračar",
	"Neimar":                    "Vračar",
	"Hram svetog Save":          "Vračar",
	"Hram Svetog Save":          "Vračar",
	"Crveni krst":               "Vračar",
	"Crveni Krst":               "Vračar",
	"Crveni Krst Vračar":        "Vračar",
	"Slavija":                   "Vračar",
	"Cvetni trg":                "Vračar",
	"Cvetni Trg":                "Vračar",
	"Južni bulevar":             "Vračar",
	"Južni Bulevar":             "Vračar",
	"Vukov spomenik":            "Vračar",
	"Vukov Spomenik":            "Vračar",
	"Krunski venac":             "Vračar",
	"Krunski Venac":             "Vračar",
	"Bulevar kralja Aleksandra": "Vračar",
	"Bulevar Kr. Aleksandra":    "Vračar",
	"Dorćol":                    "Stari Grad",
	"Gornji Dorćol":             "Stari Grad",
	"Donji Dorćol":              "Stari Grad",
	"Skadarlija":                "Stari Grad",
	"Terazije":                  "Stari Grad",
	"Knez Mihailova":            "Stari Grad",
	"Kosančićev venac":          "Stari Grad",
	"Kosančićev Venac":          "Stari Grad",
	"Kopitareva gradina":        "Stari Grad",
	"Kopitareva Gradina":        "Stari Grad",
	"Kalemegdan":                "Stari Grad",
	"Skupština":                 "Stari Grad",
	"Bajlonijeva pijaca":        "Stari Grad",
	"Bajlonijeva Pijaca":        "Stari Grad",
	"Gundulićev venac":          "Stari Grad",
	"Gundulićev Venac":          "Stari Grad",
	"Obilićev venac":            "Stari Grad",
	"Obilićev Venac":            "Stari Grad",
	"Trg Republike":             "Stari Grad",
	"Miljakovac":                "Rakovica",
	"Miljakovac I":              "Rakovica",
	"Miljakovac II":             "Rakovica",
	"Miljakovac III":            "Rakovica",
	"Kanarevo brdo":             "Rakovica",
	"Kanarevo Brdo":             "Rakovica",
	"Labudovo brdo":             "Rakovica",
	"Labudovo Brdo":             "Rakovica",
	"Petlovo brdo":              "Rakovica",
	"Petlovo Brdo":              "Rakovica",
	"Resnik":                    "Rakovica",
	"Kneževac":                  "Rakovica",
	"Skojevsko naselje":         "Rakovica",
	"Skojevsko Naselje":         "Rakovica",
}

func GetAllStandardizedDistricts() []string {
	unique := make(map[string]struct{})
	for _, val := range districtMapping {
		unique[val] = struct{}{}
	}
	districts := []string{
		"Novi Beograd", "Zemun", "Vračar", "Savski venac", "Stari Grad",
		"Palilula", "Zvezdara", "Voždovac", "Čukarica", "Rakovica", "Grocka",
	}
	for _, d := range districts {
		unique[d] = struct{}{}
	}
	res := make([]string, 0, len(unique))
	for k := range unique {
		res = append(res, k)
	}
	sort.Strings(res)
	return res
}

func StandardizeDistrict(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "Unknown"
	}

	if strings.Contains(strings.ToLower(raw), "blok") {
		return "Novi Beograd"
	}

	if strings.HasPrefix(strings.ToLower(raw), "novi beograd") {
		return "Novi Beograd"
	}

	if strings.HasPrefix(strings.ToLower(raw), "bežanijska kosa") {
		return "Novi Beograd"
	}

	if val, ok := districtMapping[raw]; ok {
		return val
	}

	for key, val := range districtMapping {
		if strings.Contains(strings.ToLower(raw), strings.ToLower(key)) {
			return val
		}
	}

	return raw
}
