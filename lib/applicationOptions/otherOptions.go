package applicationOptions

// -=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
//																		     	 ThemeTypes

// ThemeType identifier for the various color themes supported by the app
type ThemeType int64

const (
	// ThtUndefined Theme is not defined
	ThtUndefined ThemeType = iota
	// ThtUnknown Theme type provided but not currently supported
	ThtUnknown
	// ThtLight color theme for light colored backgrounds
	ThtLight
	// ThtDark color theme for dark colored backgrounds
	ThtDark
)

// ToString Converts the value from our enumeration to a string representation
func (t *ThemeType) ToString() string {
	switch *t {
	case ThtLight:
		return "light"
	case ThtDark:
		return "dark"
	case ThtUndefined:
		fallthrough
	case ThtUnknown:
		fallthrough
	default:
		return ""
	}
}

// FromString populates our enumeration from an arbitrary character string
func (t *ThemeType) FromString(value string) {
	switch value {
	case "light":
		*t = ThtLight
	case "dark":
		*t = ThtDark
	case "":
		*t = ThtUndefined
	default:
		*t = ThtUnknown
	}
}

// -=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
//																	           OtherOptions

type OtherOptions struct {
	// Theme color scheme to use when presenting colored output
	Theme ThemeType
}
