package gen

import (
	"errors"
	"strconv"
	"strings"
)

var modifiers = map[string]string{
	// aliases
	"c":   "control",
	"a":   "option",
	"o":   "option",
	"s":   "shift",
	"cmd": "command",

	// original
	"caps_lock":     "caps_lock",
	"left_command":  "left_command",
	"left_control":  "left_control",
	"left_option":   "left_option",
	"left_shift":    "left_shift",
	"right_command": "right_command",
	"right_control": "right_control",
	"right_option":  "right_option",
	"right_shift":   "right_shift",
	"fn":            "fn",
	"command":       "command",
	"control":       "control",
	"option":        "option",
	"shift":         "shift",
	"left_alt":      "left_alt",
	"left_gui":      "left_gui",
	"right_alt":     "right_alt",
	"right_gui":     "right_gui",
	"any":           "any",
}

var keys = map[string]string{
	"caps_lock":                          "caps_lock",
	"left_control":                       "left_control",
	"left_shift":                         "left_shift",
	"left_option":                        "left_option",
	"left_command":                       "left_command",
	"right_control":                      "right_control",
	"right_shift":                        "right_shift",
	"right_option":                       "right_option",
	"right_command":                      "right_command",
	"fn":                                 "fn",
	"return_or_enter":                    "return_or_enter",
	"escape":                             "escape",
	"delete_or_backspace":                "delete_or_backspace",
	"delete_forward":                     "delete_forward",
	"tab":                                "tab",
	"spacebar":                           "spacebar",
	"hyphen":                             "hyphen",
	"equal_sign":                         "equal_sign",
	"open_bracket":                       "open_bracket",
	"close_bracket":                      "close_bracket",
	"backslash":                          "backslash",
	"non_us_pound":                       "non_us_pound",
	"semicolon":                          "semicolon",
	"quote":                              "quote",
	"grave_accent_and_tilde":             "grave_accent_and_tilde",
	"comma":                              "comma",
	"period":                             "period",
	"slash":                              "slash",
	"non_us_backslash":                   "non_us_backslash",
	"up_arrow":                           "up_arrow",
	"down_arrow":                         "down_arrow",
	"left_arrow":                         "left_arrow",
	"right_arrow":                        "right_arrow",
	"page_up":                            "page_up",
	"page_down":                          "page_down",
	"home":                               "home",
	"end":                                "end",
	"a":                                  "a",
	"b":                                  "b",
	"c":                                  "c",
	"d":                                  "d",
	"e":                                  "e",
	"f":                                  "f",
	"g":                                  "g",
	"h":                                  "h",
	"i":                                  "i",
	"j":                                  "j",
	"k":                                  "k",
	"l":                                  "l",
	"m":                                  "m",
	"n":                                  "n",
	"o":                                  "o",
	"p":                                  "p",
	"q":                                  "q",
	"r":                                  "r",
	"s":                                  "s",
	"t":                                  "t",
	"u":                                  "u",
	"v":                                  "v",
	"w":                                  "w",
	"x":                                  "x",
	"y":                                  "y",
	"z":                                  "z",
	"1":                                  "1",
	"2":                                  "2",
	"3":                                  "3",
	"4":                                  "4",
	"5":                                  "5",
	"6":                                  "6",
	"7":                                  "7",
	"8":                                  "8",
	"9":                                  "9",
	"0":                                  "0",
	"f1":                                 "f1",
	"f2":                                 "f2",
	"f3":                                 "f3",
	"f4":                                 "f4",
	"f5":                                 "f5",
	"f6":                                 "f6",
	"f7":                                 "f7",
	"f8":                                 "f8",
	"f9":                                 "f9",
	"f10":                                "f10",
	"f11":                                "f11",
	"f12":                                "f12",
	"f13":                                "f13",
	"f14":                                "f14",
	"f15":                                "f15",
	"f16":                                "f16",
	"f17":                                "f17",
	"f18":                                "f18",
	"f19":                                "f19",
	"f20":                                "f20",
	"f21":                                "f21",
	"f22":                                "f22",
	"f23":                                "f23",
	"f24":                                "f24",
	"mission_control":                    "mission_control",
	"launchpad":                          "launchpad",
	"dashboard":                          "dashboard",
	"illumination_decrement":             "illumination_decrement",
	"illumination_increment":             "illumination_increment",
	"apple_display_brightness_decrement": "apple_display_brightness_decrement",
	"apple_display_brightness_increment": "apple_display_brightness_increment",
	"apple_top_case_display_brightness_decrement": "apple_top_case_display_brightness_decrement",
	"apple_top_case_display_brightness_increment": "apple_top_case_display_brightness_increment",
	"keypad_num_lock":               "keypad_num_lock",
	"keypad_slash":                  "keypad_slash",
	"keypad_asterisk":               "keypad_asterisk",
	"keypad_hyphen":                 "keypad_hyphen",
	"keypad_plus":                   "keypad_plus",
	"keypad_enter":                  "keypad_enter",
	"keypad_1":                      "keypad_1",
	"keypad_2":                      "keypad_2",
	"keypad_3":                      "keypad_3",
	"keypad_4":                      "keypad_4",
	"keypad_5":                      "keypad_5",
	"keypad_6":                      "keypad_6",
	"keypad_7":                      "keypad_7",
	"keypad_8":                      "keypad_8",
	"keypad_9":                      "keypad_9",
	"keypad_0":                      "keypad_0",
	"keypad_period":                 "keypad_period",
	"keypad_equal_sign":             "keypad_equal_sign",
	"keypad_comma":                  "keypad_comma",
	"vk_none":                       "vk_none",
	"print_screen":                  "print_screen",
	"scroll_lock":                   "scroll_lock",
	"pause":                         "pause",
	"insert":                        "insert",
	"application":                   "application",
	"help":                          "help",
	"power":                         "power",
	"execute":                       "execute",
	"menu":                          "menu",
	"select":                        "select",
	"stop":                          "stop",
	"again":                         "again",
	"undo":                          "undo",
	"cut":                           "cut",
	"copy":                          "copy",
	"paste":                         "paste",
	"find":                          "find",
	"international1":                "international1",
	"international2":                "international2",
	"international3":                "international3",
	"international4":                "international4",
	"international5":                "international5",
	"international6":                "international6",
	"international7":                "international7",
	"international8":                "international8",
	"international9":                "international9",
	"lang1":                         "lang1",
	"lang2":                         "lang2",
	"lang3":                         "lang3",
	"lang4":                         "lang4",
	"lang5":                         "lang5",
	"lang6":                         "lang6",
	"lang7":                         "lang7",
	"lang8":                         "lang8",
	"lang9":                         "lang9",
	"japanese_eisuu":                "japanese_eisuu",
	"japanese_kana":                 "japanese_kana",
	"japanese_pc_nfer":              "japanese_pc_nfer",
	"japanese_pc_xfer":              "japanese_pc_xfer",
	"japanese_pc_katakana":          "japanese_pc_katakana",
	"keypad_equal_sign_as400":       "keypad_equal_sign_as400",
	"locking_caps_lock":             "locking_caps_lock",
	"locking_num_lock":              "locking_num_lock",
	"locking_scroll_lock":           "locking_scroll_lock",
	"alternate_erase":               "alternate_erase",
	"sys_req_or_attention":          "sys_req_or_attention",
	"cancel":                        "cancel",
	"clear":                         "clear",
	"prior":                         "prior",
	"return":                        "return",
	"separator":                     "separator",
	"out":                           "out",
	"oper":                          "oper",
	"clear_or_again":                "clear_or_again",
	"cr_sel_or_props":               "cr_sel_or_props",
	"ex_sel":                        "ex_sel",
	"left_alt":                      "left_alt",
	"left_gui":                      "left_gui",
	"right_alt":                     "right_alt",
	"right_gui":                     "right_gui",
	"vk_consumer_brightness_down":   "vk_consumer_brightness_down",
	"vk_consumer_brightness_up":     "vk_consumer_brightness_up",
	"vk_mission_control":            "vk_mission_control",
	"vk_launchpad":                  "vk_launchpad",
	"vk_dashboard":                  "vk_dashboard",
	"vk_consumer_illumination_down": "vk_consumer_illumination_down",
	"vk_consumer_illumination_up":   "vk_consumer_illumination_up",
	"vk_consumer_previous":          "vk_consumer_previous",
	"vk_consumer_play":              "vk_consumer_play",
	"vk_consumer_next":              "vk_consumer_next",
	"volume_down":                   "volume_down",
	"volume_up":                     "volume_up",
	"display_brightness_decrement":  "display_brightness_decrement",
	"display_brightness_increment":  "display_brightness_increment",
	"rewind":                        "rewind",
	"play_or_pause":                 "play_or_pause",
	"fastforward":                   "fastforward",
	"mute":                          "mute",
	"volume_decrement":              "volume_decrement",
	"volume_increment":              "volume_increment",
}

func keyCodeOf(s string) (string, error) {
	s = strings.ToLower(s)
	_, err := strconv.ParseInt(s, 10, 32)
	if err == nil {
		return s, nil
	}
	key, found := keys[s]
	if !found {
		return "", errors.New("unsupported key: " + s)
	}
	return key, nil
}

func modifierOf(s string) (string, error) {
	s = strings.ToLower(s)
	modifier, found := modifiers[s]
	if !found {
		return "", errors.New("unsupported modifier: " + s)
	}
	return modifier, nil
}
