package search

import (
	"strings"

	"github.com/lbryio/lbry.go/v2/extras/errors"
	"github.com/lbryio/ozzo-validation/is"
)

var tayloredResults = map[string]string{
	"silvano":                "@SilvanoTrotta",
	"trotta":                 "@SilvanoTrotta",
	"silvano trotta":         "@SilvanoTrotta",
	"corbett":                "@CorbettReport",
	"linux gamer":            "thelinuxgamer",
	"linuxgamer":             "thelinuxgamer",
	"tim pool":               "timcast",
	"jordan peterson":        "jordanbpeterson",
	"quartering":             "thequartering",
	"bombards":               "Bombards_Body_Language",
	"bombard body language":  "Bombards_Body_Language",
	"bombards body language": "Bombards_Body_Language",
	"stefan molyneux":        "@freedomain",
	"crypto wendy":           "CRYPTOWENDYO",
	"Alex jones":             "alexjoneschannel",
	"styx":                   "Styxhexenhammer666",
	"styxx":                  "Styxhexenhammer666",
	"Radio Québec":           "Radio-Quebec",
	"The Alex Jones Channel": "Alex Jones Channel",
	"Alex Jones":             "Alex Jones Channel",
	"3Dto5DConsciousness":    "3D-to-5D-Consciousness",
	"PostMillennial":         "ThePostMillennial",
	"planètes":               "planetes360",
	"planetes":               "planetes360",
	"planètes 360":           "planetes360",
	"planetes 360":           "planetes360",
	"Louis Rossman":          "Louis Rossmann",
}

func checkForSpecialHandling(s string) string {
	sLower := strings.ToLower(s)
	if newSearch, ok := tayloredResults[sLower]; ok {
		return newSearch
	}
	return s
}

const limitForUsefulResults = 300

func truncate(s string) string {
	if len(s) > limitForUsefulResults {
		return s[:limitForUsefulResults]
	}
	return s
}

func collectFilters(s *string) ([]string, error) {
	if s == nil {
		return nil, nil
	}

	filters := strings.Split(*s, ",")
	for _, f := range filters {
		err := is.PrintableASCII.Validate(f)
		if err != nil {
			return nil, errors.Err(err)
		}
	}
	return filters, nil
}
