package enum

import "github.com/zxc7563598/oneadmin/internal/i18n"

type Enable int

const (
	EnableDisable Enable = iota
	EnableEnable
)

func (e Enable) Key() string {
	switch e {
	case EnableDisable:
		return "enable.disable"
	case EnableEnable:
		return "enable.enable"
	default:
		return "unknown"
	}
}

func (e Enable) Text(lang string) string {
	return i18n.T(lang, e.Key())
}
