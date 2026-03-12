package enum

import "github.com/zxc7563598/oneadmin/internal/i18n"

type YesNo int

const (
	No YesNo = iota
	Yes
)

func (y YesNo) Key() string {
	switch y {
	case No:
		return "no"
	case Yes:
		return "yes"
	default:
		return "unknown"
	}
}

func (y YesNo) Text(lang string) string {
	return i18n.T(lang, y.Key())
}
