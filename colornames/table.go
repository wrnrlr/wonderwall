package colornames

import "image/color"

// Map contains named colors defined in the SVG 1.1 spec.
var Map = map[string]color.NRGBA{
	"aliceblue":            color.NRGBA{0xf0, 0xf8, 0xff, 0xff}, // rgb(240, 248, 255)
	"antiquewhite":         color.NRGBA{0xfa, 0xeb, 0xd7, 0xff}, // rgb(250, 235, 215)
	"aqua":                 color.NRGBA{0x00, 0xff, 0xff, 0xff}, // rgb(0, 255, 255)
	"aquamarine":           color.NRGBA{0x7f, 0xff, 0xd4, 0xff}, // rgb(127, 255, 212)
	"azure":                color.NRGBA{0xf0, 0xff, 0xff, 0xff}, // rgb(240, 255, 255)
	"beige":                color.NRGBA{0xf5, 0xf5, 0xdc, 0xff}, // rgb(245, 245, 220)
	"bisque":               color.NRGBA{0xff, 0xe4, 0xc4, 0xff}, // rgb(255, 228, 196)
	"black":                color.NRGBA{0x00, 0x00, 0x00, 0xff}, // rgb(0, 0, 0)
	"blanchedalmond":       color.NRGBA{0xff, 0xeb, 0xcd, 0xff}, // rgb(255, 235, 205)
	"blue":                 color.NRGBA{0x00, 0x00, 0xff, 0xff}, // rgb(0, 0, 255)
	"blueviolet":           color.NRGBA{0x8a, 0x2b, 0xe2, 0xff}, // rgb(138, 43, 226)
	"brown":                color.NRGBA{0xa5, 0x2a, 0x2a, 0xff}, // rgb(165, 42, 42)
	"burlywood":            color.NRGBA{0xde, 0xb8, 0x87, 0xff}, // rgb(222, 184, 135)
	"cadetblue":            color.NRGBA{0x5f, 0x9e, 0xa0, 0xff}, // rgb(95, 158, 160)
	"chartreuse":           color.NRGBA{0x7f, 0xff, 0x00, 0xff}, // rgb(127, 255, 0)
	"chocolate":            color.NRGBA{0xd2, 0x69, 0x1e, 0xff}, // rgb(210, 105, 30)
	"coral":                color.NRGBA{0xff, 0x7f, 0x50, 0xff}, // rgb(255, 127, 80)
	"cornflowerblue":       color.NRGBA{0x64, 0x95, 0xed, 0xff}, // rgb(100, 149, 237)
	"cornsilk":             color.NRGBA{0xff, 0xf8, 0xdc, 0xff}, // rgb(255, 248, 220)
	"crimson":              color.NRGBA{0xdc, 0x14, 0x3c, 0xff}, // rgb(220, 20, 60)
	"cyan":                 color.NRGBA{0x00, 0xff, 0xff, 0xff}, // rgb(0, 255, 255)
	"darkblue":             color.NRGBA{0x00, 0x00, 0x8b, 0xff}, // rgb(0, 0, 139)
	"darkcyan":             color.NRGBA{0x00, 0x8b, 0x8b, 0xff}, // rgb(0, 139, 139)
	"darkgoldenrod":        color.NRGBA{0xb8, 0x86, 0x0b, 0xff}, // rgb(184, 134, 11)
	"darkgray":             color.NRGBA{0xa9, 0xa9, 0xa9, 0xff}, // rgb(169, 169, 169)
	"darkgreen":            color.NRGBA{0x00, 0x64, 0x00, 0xff}, // rgb(0, 100, 0)
	"darkgrey":             color.NRGBA{0xa9, 0xa9, 0xa9, 0xff}, // rgb(169, 169, 169)
	"darkkhaki":            color.NRGBA{0xbd, 0xb7, 0x6b, 0xff}, // rgb(189, 183, 107)
	"darkmagenta":          color.NRGBA{0x8b, 0x00, 0x8b, 0xff}, // rgb(139, 0, 139)
	"darkolivegreen":       color.NRGBA{0x55, 0x6b, 0x2f, 0xff}, // rgb(85, 107, 47)
	"darkorange":           color.NRGBA{0xff, 0x8c, 0x00, 0xff}, // rgb(255, 140, 0)
	"darkorchid":           color.NRGBA{0x99, 0x32, 0xcc, 0xff}, // rgb(153, 50, 204)
	"darkred":              color.NRGBA{0x8b, 0x00, 0x00, 0xff}, // rgb(139, 0, 0)
	"darksalmon":           color.NRGBA{0xe9, 0x96, 0x7a, 0xff}, // rgb(233, 150, 122)
	"darkseagreen":         color.NRGBA{0x8f, 0xbc, 0x8f, 0xff}, // rgb(143, 188, 143)
	"darkslateblue":        color.NRGBA{0x48, 0x3d, 0x8b, 0xff}, // rgb(72, 61, 139)
	"darkslategray":        color.NRGBA{0x2f, 0x4f, 0x4f, 0xff}, // rgb(47, 79, 79)
	"darkslategrey":        color.NRGBA{0x2f, 0x4f, 0x4f, 0xff}, // rgb(47, 79, 79)
	"darkturquoise":        color.NRGBA{0x00, 0xce, 0xd1, 0xff}, // rgb(0, 206, 209)
	"darkviolet":           color.NRGBA{0x94, 0x00, 0xd3, 0xff}, // rgb(148, 0, 211)
	"deeppink":             color.NRGBA{0xff, 0x14, 0x93, 0xff}, // rgb(255, 20, 147)
	"deepskyblue":          color.NRGBA{0x00, 0xbf, 0xff, 0xff}, // rgb(0, 191, 255)
	"dimgray":              color.NRGBA{0x69, 0x69, 0x69, 0xff}, // rgb(105, 105, 105)
	"dimgrey":              color.NRGBA{0x69, 0x69, 0x69, 0xff}, // rgb(105, 105, 105)
	"dodgerblue":           color.NRGBA{0x1e, 0x90, 0xff, 0xff}, // rgb(30, 144, 255)
	"firebrick":            color.NRGBA{0xb2, 0x22, 0x22, 0xff}, // rgb(178, 34, 34)
	"floralwhite":          color.NRGBA{0xff, 0xfa, 0xf0, 0xff}, // rgb(255, 250, 240)
	"forestgreen":          color.NRGBA{0x22, 0x8b, 0x22, 0xff}, // rgb(34, 139, 34)
	"fuchsia":              color.NRGBA{0xff, 0x00, 0xff, 0xff}, // rgb(255, 0, 255)
	"gainsboro":            color.NRGBA{0xdc, 0xdc, 0xdc, 0xff}, // rgb(220, 220, 220)
	"ghostwhite":           color.NRGBA{0xf8, 0xf8, 0xff, 0xff}, // rgb(248, 248, 255)
	"gold":                 color.NRGBA{0xff, 0xd7, 0x00, 0xff}, // rgb(255, 215, 0)
	"goldenrod":            color.NRGBA{0xda, 0xa5, 0x20, 0xff}, // rgb(218, 165, 32)
	"gray":                 color.NRGBA{0x80, 0x80, 0x80, 0xff}, // rgb(128, 128, 128)
	"green":                color.NRGBA{0x00, 0x80, 0x00, 0xff}, // rgb(0, 128, 0)
	"greenyellow":          color.NRGBA{0xad, 0xff, 0x2f, 0xff}, // rgb(173, 255, 47)
	"grey":                 color.NRGBA{0x80, 0x80, 0x80, 0xff}, // rgb(128, 128, 128)
	"honeydew":             color.NRGBA{0xf0, 0xff, 0xf0, 0xff}, // rgb(240, 255, 240)
	"hotpink":              color.NRGBA{0xff, 0x69, 0xb4, 0xff}, // rgb(255, 105, 180)
	"indianred":            color.NRGBA{0xcd, 0x5c, 0x5c, 0xff}, // rgb(205, 92, 92)
	"indigo":               color.NRGBA{0x4b, 0x00, 0x82, 0xff}, // rgb(75, 0, 130)
	"ivory":                color.NRGBA{0xff, 0xff, 0xf0, 0xff}, // rgb(255, 255, 240)
	"khaki":                color.NRGBA{0xf0, 0xe6, 0x8c, 0xff}, // rgb(240, 230, 140)
	"lavender":             color.NRGBA{0xe6, 0xe6, 0xfa, 0xff}, // rgb(230, 230, 250)
	"lavenderblush":        color.NRGBA{0xff, 0xf0, 0xf5, 0xff}, // rgb(255, 240, 245)
	"lawngreen":            color.NRGBA{0x7c, 0xfc, 0x00, 0xff}, // rgb(124, 252, 0)
	"lemonchiffon":         color.NRGBA{0xff, 0xfa, 0xcd, 0xff}, // rgb(255, 250, 205)
	"lightblue":            color.NRGBA{0xad, 0xd8, 0xe6, 0xff}, // rgb(173, 216, 230)
	"lightcoral":           color.NRGBA{0xf0, 0x80, 0x80, 0xff}, // rgb(240, 128, 128)
	"lightcyan":            color.NRGBA{0xe0, 0xff, 0xff, 0xff}, // rgb(224, 255, 255)
	"lightgoldenrodyellow": color.NRGBA{0xfa, 0xfa, 0xd2, 0xff}, // rgb(250, 250, 210)
	"lightgray":            color.NRGBA{0xd3, 0xd3, 0xd3, 0xff}, // rgb(211, 211, 211)
	"lightgreen":           color.NRGBA{0x90, 0xee, 0x90, 0xff}, // rgb(144, 238, 144)
	"lightgrey":            color.NRGBA{0xd3, 0xd3, 0xd3, 0xff}, // rgb(211, 211, 211)
	"lightpink":            color.NRGBA{0xff, 0xb6, 0xc1, 0xff}, // rgb(255, 182, 193)
	"lightsalmon":          color.NRGBA{0xff, 0xa0, 0x7a, 0xff}, // rgb(255, 160, 122)
	"lightseagreen":        color.NRGBA{0x20, 0xb2, 0xaa, 0xff}, // rgb(32, 178, 170)
	"lightskyblue":         color.NRGBA{0x87, 0xce, 0xfa, 0xff}, // rgb(135, 206, 250)
	"lightslategray":       color.NRGBA{0x77, 0x88, 0x99, 0xff}, // rgb(119, 136, 153)
	"lightslategrey":       color.NRGBA{0x77, 0x88, 0x99, 0xff}, // rgb(119, 136, 153)
	"lightsteelblue":       color.NRGBA{0xb0, 0xc4, 0xde, 0xff}, // rgb(176, 196, 222)
	"lightyellow":          color.NRGBA{0xff, 0xff, 0xe0, 0xff}, // rgb(255, 255, 224)
	"lime":                 color.NRGBA{0x00, 0xff, 0x00, 0xff}, // rgb(0, 255, 0)
	"limegreen":            color.NRGBA{0x32, 0xcd, 0x32, 0xff}, // rgb(50, 205, 50)
	"linen":                color.NRGBA{0xfa, 0xf0, 0xe6, 0xff}, // rgb(250, 240, 230)
	"magenta":              color.NRGBA{0xff, 0x00, 0xff, 0xff}, // rgb(255, 0, 255)
	"maroon":               color.NRGBA{0x80, 0x00, 0x00, 0xff}, // rgb(128, 0, 0)
	"mediumaquamarine":     color.NRGBA{0x66, 0xcd, 0xaa, 0xff}, // rgb(102, 205, 170)
	"mediumblue":           color.NRGBA{0x00, 0x00, 0xcd, 0xff}, // rgb(0, 0, 205)
	"mediumorchid":         color.NRGBA{0xba, 0x55, 0xd3, 0xff}, // rgb(186, 85, 211)
	"mediumpurple":         color.NRGBA{0x93, 0x70, 0xdb, 0xff}, // rgb(147, 112, 219)
	"mediumseagreen":       color.NRGBA{0x3c, 0xb3, 0x71, 0xff}, // rgb(60, 179, 113)
	"mediumslateblue":      color.NRGBA{0x7b, 0x68, 0xee, 0xff}, // rgb(123, 104, 238)
	"mediumspringgreen":    color.NRGBA{0x00, 0xfa, 0x9a, 0xff}, // rgb(0, 250, 154)
	"mediumturquoise":      color.NRGBA{0x48, 0xd1, 0xcc, 0xff}, // rgb(72, 209, 204)
	"mediumvioletred":      color.NRGBA{0xc7, 0x15, 0x85, 0xff}, // rgb(199, 21, 133)
	"midnightblue":         color.NRGBA{0x19, 0x19, 0x70, 0xff}, // rgb(25, 25, 112)
	"mintcream":            color.NRGBA{0xf5, 0xff, 0xfa, 0xff}, // rgb(245, 255, 250)
	"mistyrose":            color.NRGBA{0xff, 0xe4, 0xe1, 0xff}, // rgb(255, 228, 225)
	"moccasin":             color.NRGBA{0xff, 0xe4, 0xb5, 0xff}, // rgb(255, 228, 181)
	"navajowhite":          color.NRGBA{0xff, 0xde, 0xad, 0xff}, // rgb(255, 222, 173)
	"navy":                 color.NRGBA{0x00, 0x00, 0x80, 0xff}, // rgb(0, 0, 128)
	"oldlace":              color.NRGBA{0xfd, 0xf5, 0xe6, 0xff}, // rgb(253, 245, 230)
	"olive":                color.NRGBA{0x80, 0x80, 0x00, 0xff}, // rgb(128, 128, 0)
	"olivedrab":            color.NRGBA{0x6b, 0x8e, 0x23, 0xff}, // rgb(107, 142, 35)
	"orange":               color.NRGBA{0xff, 0xa5, 0x00, 0xff}, // rgb(255, 165, 0)
	"orangered":            color.NRGBA{0xff, 0x45, 0x00, 0xff}, // rgb(255, 69, 0)
	"orchid":               color.NRGBA{0xda, 0x70, 0xd6, 0xff}, // rgb(218, 112, 214)
	"palegoldenrod":        color.NRGBA{0xee, 0xe8, 0xaa, 0xff}, // rgb(238, 232, 170)
	"palegreen":            color.NRGBA{0x98, 0xfb, 0x98, 0xff}, // rgb(152, 251, 152)
	"paleturquoise":        color.NRGBA{0xaf, 0xee, 0xee, 0xff}, // rgb(175, 238, 238)
	"palevioletred":        color.NRGBA{0xdb, 0x70, 0x93, 0xff}, // rgb(219, 112, 147)
	"papayawhip":           color.NRGBA{0xff, 0xef, 0xd5, 0xff}, // rgb(255, 239, 213)
	"peachpuff":            color.NRGBA{0xff, 0xda, 0xb9, 0xff}, // rgb(255, 218, 185)
	"peru":                 color.NRGBA{0xcd, 0x85, 0x3f, 0xff}, // rgb(205, 133, 63)
	"pink":                 color.NRGBA{0xff, 0xc0, 0xcb, 0xff}, // rgb(255, 192, 203)
	"plum":                 color.NRGBA{0xdd, 0xa0, 0xdd, 0xff}, // rgb(221, 160, 221)
	"powderblue":           color.NRGBA{0xb0, 0xe0, 0xe6, 0xff}, // rgb(176, 224, 230)
	"purple":               color.NRGBA{0x80, 0x00, 0x80, 0xff}, // rgb(128, 0, 128)
	"red":                  color.NRGBA{0xff, 0x00, 0x00, 0xff}, // rgb(255, 0, 0)
	"rosybrown":            color.NRGBA{0xbc, 0x8f, 0x8f, 0xff}, // rgb(188, 143, 143)
	"royalblue":            color.NRGBA{0x41, 0x69, 0xe1, 0xff}, // rgb(65, 105, 225)
	"saddlebrown":          color.NRGBA{0x8b, 0x45, 0x13, 0xff}, // rgb(139, 69, 19)
	"salmon":               color.NRGBA{0xfa, 0x80, 0x72, 0xff}, // rgb(250, 128, 114)
	"sandybrown":           color.NRGBA{0xf4, 0xa4, 0x60, 0xff}, // rgb(244, 164, 96)
	"seagreen":             color.NRGBA{0x2e, 0x8b, 0x57, 0xff}, // rgb(46, 139, 87)
	"seashell":             color.NRGBA{0xff, 0xf5, 0xee, 0xff}, // rgb(255, 245, 238)
	"sienna":               color.NRGBA{0xa0, 0x52, 0x2d, 0xff}, // rgb(160, 82, 45)
	"silver":               color.NRGBA{0xc0, 0xc0, 0xc0, 0xff}, // rgb(192, 192, 192)
	"skyblue":              color.NRGBA{0x87, 0xce, 0xeb, 0xff}, // rgb(135, 206, 235)
	"slateblue":            color.NRGBA{0x6a, 0x5a, 0xcd, 0xff}, // rgb(106, 90, 205)
	"slategray":            color.NRGBA{0x70, 0x80, 0x90, 0xff}, // rgb(112, 128, 144)
	"slategrey":            color.NRGBA{0x70, 0x80, 0x90, 0xff}, // rgb(112, 128, 144)
	"snow":                 color.NRGBA{0xff, 0xfa, 0xfa, 0xff}, // rgb(255, 250, 250)
	"springgreen":          color.NRGBA{0x00, 0xff, 0x7f, 0xff}, // rgb(0, 255, 127)
	"steelblue":            color.NRGBA{0x46, 0x82, 0xb4, 0xff}, // rgb(70, 130, 180)
	"tan":                  color.NRGBA{0xd2, 0xb4, 0x8c, 0xff}, // rgb(210, 180, 140)
	"teal":                 color.NRGBA{0x00, 0x80, 0x80, 0xff}, // rgb(0, 128, 128)
	"thistle":              color.NRGBA{0xd8, 0xbf, 0xd8, 0xff}, // rgb(216, 191, 216)
	"tomato":               color.NRGBA{0xff, 0x63, 0x47, 0xff}, // rgb(255, 99, 71)
	"turquoise":            color.NRGBA{0x40, 0xe0, 0xd0, 0xff}, // rgb(64, 224, 208)
	"violet":               color.NRGBA{0xee, 0x82, 0xee, 0xff}, // rgb(238, 130, 238)
	"wheat":                color.NRGBA{0xf5, 0xde, 0xb3, 0xff}, // rgb(245, 222, 179)
	"white":                color.NRGBA{0xff, 0xff, 0xff, 0xff}, // rgb(255, 255, 255)
	"whitesmoke":           color.NRGBA{0xf5, 0xf5, 0xf5, 0xff}, // rgb(245, 245, 245)
	"yellow":               color.NRGBA{0xff, 0xff, 0x00, 0xff}, // rgb(255, 255, 0)
	"yellowgreen":          color.NRGBA{0x9a, 0xcd, 0x32, 0xff}, // rgb(154, 205, 50)
}

// Names contains the color names defined in the SVG 1.1 spec.
var Names = []string{
	"aliceblue",
	"antiquewhite",
	"aqua",
	"aquamarine",
	"azure",
	"beige",
	"bisque",
	"black",
	"blanchedalmond",
	"blue",
	"blueviolet",
	"brown",
	"burlywood",
	"cadetblue",
	"chartreuse",
	"chocolate",
	"coral",
	"cornflowerblue",
	"cornsilk",
	"crimson",
	"cyan",
	"darkblue",
	"darkcyan",
	"darkgoldenrod",
	"darkgray",
	"darkgreen",
	"darkgrey",
	"darkkhaki",
	"darkmagenta",
	"darkolivegreen",
	"darkorange",
	"darkorchid",
	"darkred",
	"darksalmon",
	"darkseagreen",
	"darkslateblue",
	"darkslategray",
	"darkslategrey",
	"darkturquoise",
	"darkviolet",
	"deeppink",
	"deepskyblue",
	"dimgray",
	"dimgrey",
	"dodgerblue",
	"firebrick",
	"floralwhite",
	"forestgreen",
	"fuchsia",
	"gainsboro",
	"ghostwhite",
	"gold",
	"goldenrod",
	"gray",
	"green",
	"greenyellow",
	"grey",
	"honeydew",
	"hotpink",
	"indianred",
	"indigo",
	"ivory",
	"khaki",
	"lavender",
	"lavenderblush",
	"lawngreen",
	"lemonchiffon",
	"lightblue",
	"lightcoral",
	"lightcyan",
	"lightgoldenrodyellow",
	"lightgray",
	"lightgreen",
	"lightgrey",
	"lightpink",
	"lightsalmon",
	"lightseagreen",
	"lightskyblue",
	"lightslategray",
	"lightslategrey",
	"lightsteelblue",
	"lightyellow",
	"lime",
	"limegreen",
	"linen",
	"magenta",
	"maroon",
	"mediumaquamarine",
	"mediumblue",
	"mediumorchid",
	"mediumpurple",
	"mediumseagreen",
	"mediumslateblue",
	"mediumspringgreen",
	"mediumturquoise",
	"mediumvioletred",
	"midnightblue",
	"mintcream",
	"mistyrose",
	"moccasin",
	"navajowhite",
	"navy",
	"oldlace",
	"olive",
	"olivedrab",
	"orange",
	"orangered",
	"orchid",
	"palegoldenrod",
	"palegreen",
	"paleturquoise",
	"palevioletred",
	"papayawhip",
	"peachpuff",
	"peru",
	"pink",
	"plum",
	"powderblue",
	"purple",
	"red",
	"rosybrown",
	"royalblue",
	"saddlebrown",
	"salmon",
	"sandybrown",
	"seagreen",
	"seashell",
	"sienna",
	"silver",
	"skyblue",
	"slateblue",
	"slategray",
	"slategrey",
	"snow",
	"springgreen",
	"steelblue",
	"tan",
	"teal",
	"thistle",
	"tomato",
	"turquoise",
	"violet",
	"wheat",
	"white",
	"whitesmoke",
	"yellow",
	"yellowgreen",
}

var (
	Aliceblue            = color.NRGBA{0xf0, 0xf8, 0xff, 0xff} // rgb(240, 248, 255)
	Antiquewhite         = color.NRGBA{0xfa, 0xeb, 0xd7, 0xff} // rgb(250, 235, 215)
	Aqua                 = color.NRGBA{0x00, 0xff, 0xff, 0xff} // rgb(0, 255, 255)
	Aquamarine           = color.NRGBA{0x7f, 0xff, 0xd4, 0xff} // rgb(127, 255, 212)
	Azure                = color.NRGBA{0xf0, 0xff, 0xff, 0xff} // rgb(240, 255, 255)
	Beige                = color.NRGBA{0xf5, 0xf5, 0xdc, 0xff} // rgb(245, 245, 220)
	Bisque               = color.NRGBA{0xff, 0xe4, 0xc4, 0xff} // rgb(255, 228, 196)
	Black                = color.NRGBA{0x00, 0x00, 0x00, 0xff} // rgb(0, 0, 0)
	Blanchedalmond       = color.NRGBA{0xff, 0xeb, 0xcd, 0xff} // rgb(255, 235, 205)
	Blue                 = color.NRGBA{0x00, 0x00, 0xff, 0xff} // rgb(0, 0, 255)
	Blueviolet           = color.NRGBA{0x8a, 0x2b, 0xe2, 0xff} // rgb(138, 43, 226)
	Brown                = color.NRGBA{0xa5, 0x2a, 0x2a, 0xff} // rgb(165, 42, 42)
	Burlywood            = color.NRGBA{0xde, 0xb8, 0x87, 0xff} // rgb(222, 184, 135)
	Cadetblue            = color.NRGBA{0x5f, 0x9e, 0xa0, 0xff} // rgb(95, 158, 160)
	Chartreuse           = color.NRGBA{0x7f, 0xff, 0x00, 0xff} // rgb(127, 255, 0)
	Chocolate            = color.NRGBA{0xd2, 0x69, 0x1e, 0xff} // rgb(210, 105, 30)
	Coral                = color.NRGBA{0xff, 0x7f, 0x50, 0xff} // rgb(255, 127, 80)
	Cornflowerblue       = color.NRGBA{0x64, 0x95, 0xed, 0xff} // rgb(100, 149, 237)
	Cornsilk             = color.NRGBA{0xff, 0xf8, 0xdc, 0xff} // rgb(255, 248, 220)
	Crimson              = color.NRGBA{0xdc, 0x14, 0x3c, 0xff} // rgb(220, 20, 60)
	Cyan                 = color.NRGBA{0x00, 0xff, 0xff, 0xff} // rgb(0, 255, 255)
	Darkblue             = color.NRGBA{0x00, 0x00, 0x8b, 0xff} // rgb(0, 0, 139)
	Darkcyan             = color.NRGBA{0x00, 0x8b, 0x8b, 0xff} // rgb(0, 139, 139)
	Darkgoldenrod        = color.NRGBA{0xb8, 0x86, 0x0b, 0xff} // rgb(184, 134, 11)
	Darkgray             = color.NRGBA{0xa9, 0xa9, 0xa9, 0xff} // rgb(169, 169, 169)
	Darkgreen            = color.NRGBA{0x00, 0x64, 0x00, 0xff} // rgb(0, 100, 0)
	Darkgrey             = color.NRGBA{0xa9, 0xa9, 0xa9, 0xff} // rgb(169, 169, 169)
	Darkkhaki            = color.NRGBA{0xbd, 0xb7, 0x6b, 0xff} // rgb(189, 183, 107)
	Darkmagenta          = color.NRGBA{0x8b, 0x00, 0x8b, 0xff} // rgb(139, 0, 139)
	Darkolivegreen       = color.NRGBA{0x55, 0x6b, 0x2f, 0xff} // rgb(85, 107, 47)
	Darkorange           = color.NRGBA{0xff, 0x8c, 0x00, 0xff} // rgb(255, 140, 0)
	Darkorchid           = color.NRGBA{0x99, 0x32, 0xcc, 0xff} // rgb(153, 50, 204)
	Darkred              = color.NRGBA{0x8b, 0x00, 0x00, 0xff} // rgb(139, 0, 0)
	Darksalmon           = color.NRGBA{0xe9, 0x96, 0x7a, 0xff} // rgb(233, 150, 122)
	Darkseagreen         = color.NRGBA{0x8f, 0xbc, 0x8f, 0xff} // rgb(143, 188, 143)
	Darkslateblue        = color.NRGBA{0x48, 0x3d, 0x8b, 0xff} // rgb(72, 61, 139)
	Darkslategray        = color.NRGBA{0x2f, 0x4f, 0x4f, 0xff} // rgb(47, 79, 79)
	Darkslategrey        = color.NRGBA{0x2f, 0x4f, 0x4f, 0xff} // rgb(47, 79, 79)
	Darkturquoise        = color.NRGBA{0x00, 0xce, 0xd1, 0xff} // rgb(0, 206, 209)
	Darkviolet           = color.NRGBA{0x94, 0x00, 0xd3, 0xff} // rgb(148, 0, 211)
	Deeppink             = color.NRGBA{0xff, 0x14, 0x93, 0xff} // rgb(255, 20, 147)
	Deepskyblue          = color.NRGBA{0x00, 0xbf, 0xff, 0xff} // rgb(0, 191, 255)
	Dimgray              = color.NRGBA{0x69, 0x69, 0x69, 0xff} // rgb(105, 105, 105)
	Dimgrey              = color.NRGBA{0x69, 0x69, 0x69, 0xff} // rgb(105, 105, 105)
	Dodgerblue           = color.NRGBA{0x1e, 0x90, 0xff, 0xff} // rgb(30, 144, 255)
	Firebrick            = color.NRGBA{0xb2, 0x22, 0x22, 0xff} // rgb(178, 34, 34)
	Floralwhite          = color.NRGBA{0xff, 0xfa, 0xf0, 0xff} // rgb(255, 250, 240)
	Forestgreen          = color.NRGBA{0x22, 0x8b, 0x22, 0xff} // rgb(34, 139, 34)
	Fuchsia              = color.NRGBA{0xff, 0x00, 0xff, 0xff} // rgb(255, 0, 255)
	Gainsboro            = color.NRGBA{0xdc, 0xdc, 0xdc, 0xff} // rgb(220, 220, 220)
	Ghostwhite           = color.NRGBA{0xf8, 0xf8, 0xff, 0xff} // rgb(248, 248, 255)
	Gold                 = color.NRGBA{0xff, 0xd7, 0x00, 0xff} // rgb(255, 215, 0)
	Goldenrod            = color.NRGBA{0xda, 0xa5, 0x20, 0xff} // rgb(218, 165, 32)
	Gray                 = color.NRGBA{0x80, 0x80, 0x80, 0xff} // rgb(128, 128, 128)
	Green                = color.NRGBA{0x00, 0x80, 0x00, 0xff} // rgb(0, 128, 0)
	Greenyellow          = color.NRGBA{0xad, 0xff, 0x2f, 0xff} // rgb(173, 255, 47)
	Grey                 = color.NRGBA{0x80, 0x80, 0x80, 0xff} // rgb(128, 128, 128)
	Honeydew             = color.NRGBA{0xf0, 0xff, 0xf0, 0xff} // rgb(240, 255, 240)
	Hotpink              = color.NRGBA{0xff, 0x69, 0xb4, 0xff} // rgb(255, 105, 180)
	Indianred            = color.NRGBA{0xcd, 0x5c, 0x5c, 0xff} // rgb(205, 92, 92)
	Indigo               = color.NRGBA{0x4b, 0x00, 0x82, 0xff} // rgb(75, 0, 130)
	Ivory                = color.NRGBA{0xff, 0xff, 0xf0, 0xff} // rgb(255, 255, 240)
	Khaki                = color.NRGBA{0xf0, 0xe6, 0x8c, 0xff} // rgb(240, 230, 140)
	Lavender             = color.NRGBA{0xe6, 0xe6, 0xfa, 0xff} // rgb(230, 230, 250)
	Lavenderblush        = color.NRGBA{0xff, 0xf0, 0xf5, 0xff} // rgb(255, 240, 245)
	Lawngreen            = color.NRGBA{0x7c, 0xfc, 0x00, 0xff} // rgb(124, 252, 0)
	Lemonchiffon         = color.NRGBA{0xff, 0xfa, 0xcd, 0xff} // rgb(255, 250, 205)
	Lightblue            = color.NRGBA{0xad, 0xd8, 0xe6, 0xff} // rgb(173, 216, 230)
	Lightcoral           = color.NRGBA{0xf0, 0x80, 0x80, 0xff} // rgb(240, 128, 128)
	Lightcyan            = color.NRGBA{0xe0, 0xff, 0xff, 0xff} // rgb(224, 255, 255)
	Lightgoldenrodyellow = color.NRGBA{0xfa, 0xfa, 0xd2, 0xff} // rgb(250, 250, 210)
	Lightgray            = color.NRGBA{0xd3, 0xd3, 0xd3, 0xff} // rgb(211, 211, 211)
	Lightgreen           = color.NRGBA{0x90, 0xee, 0x90, 0xff} // rgb(144, 238, 144)
	Lightgrey            = color.NRGBA{0xd3, 0xd3, 0xd3, 0xff} // rgb(211, 211, 211)
	Lightpink            = color.NRGBA{0xff, 0xb6, 0xc1, 0xff} // rgb(255, 182, 193)
	Lightsalmon          = color.NRGBA{0xff, 0xa0, 0x7a, 0xff} // rgb(255, 160, 122)
	Lightseagreen        = color.NRGBA{0x20, 0xb2, 0xaa, 0xff} // rgb(32, 178, 170)
	Lightskyblue         = color.NRGBA{0x87, 0xce, 0xfa, 0xff} // rgb(135, 206, 250)
	Lightslategray       = color.NRGBA{0x77, 0x88, 0x99, 0xff} // rgb(119, 136, 153)
	Lightslategrey       = color.NRGBA{0x77, 0x88, 0x99, 0xff} // rgb(119, 136, 153)
	Lightsteelblue       = color.NRGBA{0xb0, 0xc4, 0xde, 0xff} // rgb(176, 196, 222)
	Lightyellow          = color.NRGBA{0xff, 0xff, 0xe0, 0xff} // rgb(255, 255, 224)
	Lime                 = color.NRGBA{0x00, 0xff, 0x00, 0xff} // rgb(0, 255, 0)
	Limegreen            = color.NRGBA{0x32, 0xcd, 0x32, 0xff} // rgb(50, 205, 50)
	Linen                = color.NRGBA{0xfa, 0xf0, 0xe6, 0xff} // rgb(250, 240, 230)
	Magenta              = color.NRGBA{0xff, 0x00, 0xff, 0xff} // rgb(255, 0, 255)
	Maroon               = color.NRGBA{0x80, 0x00, 0x00, 0xff} // rgb(128, 0, 0)
	Mediumaquamarine     = color.NRGBA{0x66, 0xcd, 0xaa, 0xff} // rgb(102, 205, 170)
	Mediumblue           = color.NRGBA{0x00, 0x00, 0xcd, 0xff} // rgb(0, 0, 205)
	Mediumorchid         = color.NRGBA{0xba, 0x55, 0xd3, 0xff} // rgb(186, 85, 211)
	Mediumpurple         = color.NRGBA{0x93, 0x70, 0xdb, 0xff} // rgb(147, 112, 219)
	Mediumseagreen       = color.NRGBA{0x3c, 0xb3, 0x71, 0xff} // rgb(60, 179, 113)
	Mediumslateblue      = color.NRGBA{0x7b, 0x68, 0xee, 0xff} // rgb(123, 104, 238)
	Mediumspringgreen    = color.NRGBA{0x00, 0xfa, 0x9a, 0xff} // rgb(0, 250, 154)
	Mediumturquoise      = color.NRGBA{0x48, 0xd1, 0xcc, 0xff} // rgb(72, 209, 204)
	Mediumvioletred      = color.NRGBA{0xc7, 0x15, 0x85, 0xff} // rgb(199, 21, 133)
	Midnightblue         = color.NRGBA{0x19, 0x19, 0x70, 0xff} // rgb(25, 25, 112)
	Mintcream            = color.NRGBA{0xf5, 0xff, 0xfa, 0xff} // rgb(245, 255, 250)
	Mistyrose            = color.NRGBA{0xff, 0xe4, 0xe1, 0xff} // rgb(255, 228, 225)
	Moccasin             = color.NRGBA{0xff, 0xe4, 0xb5, 0xff} // rgb(255, 228, 181)
	Navajowhite          = color.NRGBA{0xff, 0xde, 0xad, 0xff} // rgb(255, 222, 173)
	Navy                 = color.NRGBA{0x00, 0x00, 0x80, 0xff} // rgb(0, 0, 128)
	Oldlace              = color.NRGBA{0xfd, 0xf5, 0xe6, 0xff} // rgb(253, 245, 230)
	Olive                = color.NRGBA{0x80, 0x80, 0x00, 0xff} // rgb(128, 128, 0)
	Olivedrab            = color.NRGBA{0x6b, 0x8e, 0x23, 0xff} // rgb(107, 142, 35)
	Orange               = color.NRGBA{0xff, 0xa5, 0x00, 0xff} // rgb(255, 165, 0)
	Orangered            = color.NRGBA{0xff, 0x45, 0x00, 0xff} // rgb(255, 69, 0)
	Orchid               = color.NRGBA{0xda, 0x70, 0xd6, 0xff} // rgb(218, 112, 214)
	Palegoldenrod        = color.NRGBA{0xee, 0xe8, 0xaa, 0xff} // rgb(238, 232, 170)
	Palegreen            = color.NRGBA{0x98, 0xfb, 0x98, 0xff} // rgb(152, 251, 152)
	Paleturquoise        = color.NRGBA{0xaf, 0xee, 0xee, 0xff} // rgb(175, 238, 238)
	Palevioletred        = color.NRGBA{0xdb, 0x70, 0x93, 0xff} // rgb(219, 112, 147)
	Papayawhip           = color.NRGBA{0xff, 0xef, 0xd5, 0xff} // rgb(255, 239, 213)
	Peachpuff            = color.NRGBA{0xff, 0xda, 0xb9, 0xff} // rgb(255, 218, 185)
	Peru                 = color.NRGBA{0xcd, 0x85, 0x3f, 0xff} // rgb(205, 133, 63)
	Pink                 = color.NRGBA{0xff, 0xc0, 0xcb, 0xff} // rgb(255, 192, 203)
	Plum                 = color.NRGBA{0xdd, 0xa0, 0xdd, 0xff} // rgb(221, 160, 221)
	Powderblue           = color.NRGBA{0xb0, 0xe0, 0xe6, 0xff} // rgb(176, 224, 230)
	Purple               = color.NRGBA{0x80, 0x00, 0x80, 0xff} // rgb(128, 0, 128)
	Red                  = color.NRGBA{0xff, 0x00, 0x00, 0xff} // rgb(255, 0, 0)
	Rosybrown            = color.NRGBA{0xbc, 0x8f, 0x8f, 0xff} // rgb(188, 143, 143)
	Royalblue            = color.NRGBA{0x41, 0x69, 0xe1, 0xff} // rgb(65, 105, 225)
	Saddlebrown          = color.NRGBA{0x8b, 0x45, 0x13, 0xff} // rgb(139, 69, 19)
	Salmon               = color.NRGBA{0xfa, 0x80, 0x72, 0xff} // rgb(250, 128, 114)
	Sandybrown           = color.NRGBA{0xf4, 0xa4, 0x60, 0xff} // rgb(244, 164, 96)
	Seagreen             = color.NRGBA{0x2e, 0x8b, 0x57, 0xff} // rgb(46, 139, 87)
	Seashell             = color.NRGBA{0xff, 0xf5, 0xee, 0xff} // rgb(255, 245, 238)
	Sienna               = color.NRGBA{0xa0, 0x52, 0x2d, 0xff} // rgb(160, 82, 45)
	Silver               = color.NRGBA{0xc0, 0xc0, 0xc0, 0xff} // rgb(192, 192, 192)
	Skyblue              = color.NRGBA{0x87, 0xce, 0xeb, 0xff} // rgb(135, 206, 235)
	Slateblue            = color.NRGBA{0x6a, 0x5a, 0xcd, 0xff} // rgb(106, 90, 205)
	Slategray            = color.NRGBA{0x70, 0x80, 0x90, 0xff} // rgb(112, 128, 144)
	Slategrey            = color.NRGBA{0x70, 0x80, 0x90, 0xff} // rgb(112, 128, 144)
	Snow                 = color.NRGBA{0xff, 0xfa, 0xfa, 0xff} // rgb(255, 250, 250)
	Springgreen          = color.NRGBA{0x00, 0xff, 0x7f, 0xff} // rgb(0, 255, 127)
	Steelblue            = color.NRGBA{0x46, 0x82, 0xb4, 0xff} // rgb(70, 130, 180)
	Tan                  = color.NRGBA{0xd2, 0xb4, 0x8c, 0xff} // rgb(210, 180, 140)
	Teal                 = color.NRGBA{0x00, 0x80, 0x80, 0xff} // rgb(0, 128, 128)
	Thistle              = color.NRGBA{0xd8, 0xbf, 0xd8, 0xff} // rgb(216, 191, 216)
	Tomato               = color.NRGBA{0xff, 0x63, 0x47, 0xff} // rgb(255, 99, 71)
	Turquoise            = color.NRGBA{0x40, 0xe0, 0xd0, 0xff} // rgb(64, 224, 208)
	Violet               = color.NRGBA{0xee, 0x82, 0xee, 0xff} // rgb(238, 130, 238)
	Wheat                = color.NRGBA{0xf5, 0xde, 0xb3, 0xff} // rgb(245, 222, 179)
	White                = color.NRGBA{0xff, 0xff, 0xff, 0xff} // rgb(255, 255, 255)
	Whitesmoke           = color.NRGBA{0xf5, 0xf5, 0xf5, 0xff} // rgb(245, 245, 245)
	Yellow               = color.NRGBA{0xff, 0xff, 0x00, 0xff} // rgb(255, 255, 0)
	Yellowgreen          = color.NRGBA{0x9a, 0xcd, 0x32, 0xff} // rgb(154, 205, 50)
)
