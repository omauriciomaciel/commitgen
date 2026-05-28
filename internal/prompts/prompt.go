package prompts

import _ "embed"

//go:embed en.md
var English string

//go:embed pt_br.md
var PtBr string

//go:embed pr_en.md
var PREnglish string

//go:embed pr_pt_br.md
var PRPtBr string

var registry = map[string]string{
	"en":    English,
	"pt-BR": PtBr,
}

var prRegistry = map[string]string{
	"en":    PREnglish,
	"pt-BR": PRPtBr,
}

func Get(lang string) (string, bool) {
	p, ok := registry[lang]
	return p, ok
}

func GetPR(lang string) (string, bool) {
	p, ok := prRegistry[lang]
	return p, ok
}

func IsSupported(lang string) bool {
	_, ok := registry[lang]
	return ok
}
