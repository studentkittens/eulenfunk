package ambilight

// Generated by this python session:
// >>> import math
// >>> import colorsys
// >>> h = [(-math.cos(i / 5) + 1) / 2 for i in range(1000)]
// >>> s = [(+math.cos(i / 3) + 1) / 2 for i in range(1000)]
// >>> v = [(+math.sin(i / 2) + 1) / 2 for i in range(1000)]
// >>> rgb = [tuple([int(v * 255) for v in  colorsys.hsv_to_rgb(h, s, v)]) for h,s,v in zip(h, s, v)]
// >>> print('\n'.join(['{%d, %d, %d, 0},' % (v[0], v[1], v[2]) for v in rgb]))

var DEFAULT_MOODBAR = []TimedColor{
	{127, 0, 0, 0},
	{188, 16, 5, 0},
	{234, 74, 25, 0},
	{254, 161, 58, 0},
	{243, 229, 93, 0},
	{168, 203, 111, 0},
	{106, 145, 103, 0},
	{69, 82, 76, 0},
	{29, 30, 31, 0},
	{2, 2, 2, 0},
	{5, 5, 5, 0},
	{36, 35, 37, 0},
	{91, 75, 88, 0},
	{154, 106, 127, 0},
	{211, 110, 127, 0},
	{247, 88, 93, 0},
	{253, 53, 54, 0},
	{229, 21, 41, 0},
	{180, 3, 58, 0},
	{117, 0, 73, 0},
	{55, 2, 58, 0},
	{8, 1, 15, 0},
	{0, 0, 0, 0},
	{6, 12, 15, 0},
	{33, 59, 52, 0},
	{86, 119, 91, 0},
	{166, 181, 156, 0},
	{229, 229, 219, 0},
	{253, 253, 253, 0},
	{246, 244, 243, 0},
	{210, 195, 193, 0},
	{153, 124, 124, 0},
	{90, 60, 60, 0},
	{36, 21, 18, 0},
	{4, 2, 1, 0},
	{3, 2, 0, 0},
	{26, 31, 2, 0},
	{27, 83, 1, 0},
	{0, 146, 36, 0},
	{9, 204, 173, 0},
	{34, 152, 243, 0},
	{72, 69, 254, 0},
	{175, 101, 234, 0},
	{187, 112, 184, 0},
	{126, 95, 112, 0},
	{65, 57, 59, 0},
	{19, 18, 18, 0},
	{0, 0, 0, 0},
	{12, 11, 11, 0},
	{52, 47, 48, 0},
	{110, 87, 98, 0},
	{173, 110, 164, 0},
	{187, 106, 224, 0},
	{105, 78, 252, 0},
	{42, 128, 249, 0},
	{14, 216, 213, 0},
	{1, 162, 64, 0},
	{18, 98, 0, 0},
	{30, 42, 2, 0},
	{7, 6, 1, 0},
	{1, 0, 0, 0},
	{26, 14, 12, 0},
	{75, 48, 47, 0},
	{138, 107, 106, 0},
	{197, 178, 177, 0},
	{240, 235, 233, 0},
	{254, 254, 254, 0},
	{238, 238, 231, 0},
	{185, 194, 173, 0},
	{103, 134, 103, 0},
	{44, 72, 61, 0},
	{10, 21, 24, 0},
	{0, 0, 1, 0},
	{3, 1, 8, 0},
	{38, 2, 45, 0},
	{102, 0, 73, 0},
	{165, 1, 63, 0},
	{219, 15, 44, 0},
	{250, 44, 47, 0},
	{251, 80, 82, 0},
	{222, 107, 121, 0},
	{170, 109, 131, 0},
	{107, 85, 100, 0},
	{48, 44, 49, 0},
	{10, 10, 10, 0},
	{0, 0, 0, 0},
	{20, 21, 21, 0},
	{59, 68, 65, 0},
	{96, 129, 98, 0},
	{150, 190, 112, 0},
	{233, 235, 99, 0},
	{254, 182, 67, 0},
	{242, 95, 32, 0},
	{201, 26, 8, 0},
	{143, 0, 0, 0},
	{80, 3, 1, 0},
	{29, 7, 2, 0},
	{2, 1, 0, 0},
	{5, 5, 2, 0},
	{34, 39, 19, 0},
	{69, 94, 62, 0},
	{127, 157, 138, 0},
	{196, 212, 211, 0},
	{244, 246, 247, 0},
	{252, 252, 253, 0},
	{224, 216, 227, 0},
	{177, 152, 175, 0},
	{115, 83, 99, 0},
	{56, 31, 37, 0},
	{14, 5, 6, 0},
	{0, 0, 0, 0},
	{16, 1, 2, 0},
	{61, 2, 16, 0},
	{121, 0, 65, 0},
	{183, 4, 170, 0},
	{147, 22, 231, 0},
	{67, 54, 254, 0},
	{90, 170, 245, 0},
	{110, 208, 197, 0},
	{105, 151, 118, 0},
	{77, 88, 73, 0},
	{34, 35, 32, 0},
	{4, 4, 4, 0},
	{3, 3, 3, 0},
	{33, 31, 31, 0},
	{85, 72, 72, 0},
	{148, 104, 104, 0},
	{206, 121, 111, 0},
	{244, 140, 91, 0},
	{254, 183, 56, 0},
	{220, 232, 23, 0},
	{85, 185, 4, 0},
	{0, 124, 12, 0},
	{1, 63, 44, 0},
	{2, 13, 18, 0},
	{0, 0, 0, 0},
	{8, 5, 13, 0},
	{52, 30, 53, 0},
	{112, 80, 101, 0},
	{175, 149, 157, 0},
	{226, 214, 215, 0},
	{252, 251, 251, 0},
	{248, 246, 246, 0},
	{215, 199, 202, 0},
	{159, 131, 142, 0},
	{96, 65, 89, 0},
	{37, 21, 41, 0},
	{3, 2, 6, 0},
	{0, 0, 1, 0},
	{2, 24, 27, 0},
	{1, 78, 42, 0},
	{6, 140, 0, 0},
	{119, 199, 7, 0},
	{241, 229, 31, 0},
	{254, 170, 66, 0},
	{237, 134, 98, 0},
	{193, 117, 112, 0},
	{132, 97, 97, 0},
	{70, 62, 61, 0},
	{23, 22, 22, 0},
	{0, 0, 0, 0},
	{9, 9, 9, 0},
	{45, 47, 43, 0},
	{86, 104, 83, 0},
	{109, 167, 135, 0},
	{107, 216, 220, 0},
	{81, 143, 250, 0},
	{87, 45, 251, 0},
	{165, 16, 220, 0},
	{167, 1, 138, 0},
	{105, 0, 47, 0},
	{47, 2, 10, 0},
	{9, 1, 1, 0},
	{0, 0, 0, 0},
	{22, 10, 11, 0},
	{70, 42, 50, 0},
	{132, 100, 119, 0},
	{192, 170, 192, 0},
	{233, 229, 237, 0},
	{254, 254, 254, 0},
	{235, 239, 241, 0},
	{180, 200, 195, 0},
	{109, 141, 116, 0},
	{60, 78, 49, 0},
	{26, 28, 13, 0},
	{1, 1, 0, 0},
	{6, 3, 1, 0},
	{40, 7, 2, 0},
	{96, 2, 0, 0},
	{159, 3, 1, 0},
	{214, 40, 13, 0},
	{248, 117, 40, 0},
	{252, 201, 76, 0},
	{209, 226, 105, 0},
	{134, 175, 110, 0},
	{88, 113, 93, 0},
	{49, 54, 53, 0},
	{12, 13, 13, 0},
	{0, 0, 0, 0},
	{17, 17, 18, 0},
	{62, 55, 62, 0},
	{123, 93, 111, 0},
	{185, 111, 132, 0},
	{232, 102, 113, 0},
	{254, 71, 71, 0},
	{245, 35, 43, 0},
	{206, 10, 47, 0},
	{149, 0, 68, 0},
	{86, 1, 71, 0},
	{24, 2, 33, 0},
	{1, 0, 3, 0},
	{1, 2, 4, 0},
	{17, 34, 34, 0},
	{57, 88, 70, 0},
	{125, 151, 121, 0},
	{202, 208, 190, 0},
	{245, 245, 241, 0},
	{254, 253, 253, 0},
	{231, 224, 222, 0},
	{183, 160, 159, 0},
	{121, 89, 89, 0},
	{61, 37, 35, 0},
	{17, 9, 7, 0},
	{0, 0, 0, 0},
	{14, 13, 1, 0},
	{33, 55, 2, 0},
	{4, 115, 0, 0},
	{3, 177, 97, 0},
	{19, 198, 227, 0},
	{51, 105, 253, 0},
	{134, 87, 248, 0},
	{193, 109, 213, 0},
	{157, 106, 145, 0},
	{94, 77, 84, 0},
	{39, 36, 37, 0},
	{6, 6, 6, 0},
	{2, 2, 2, 0},
	{29, 27, 27, 0},
	{80, 68, 72, 0},
	{142, 101, 128, 0},
	{194, 111, 201, 0},
	{157, 94, 242, 0},
	{60, 85, 254, 0},
	{26, 176, 236, 0},
	{5, 191, 132, 0},
	{0, 130, 12, 0},
	{32, 68, 1, 0},
	{20, 21, 2, 0},
	{0, 0, 0, 0},
	{10, 5, 3, 0},
	{48, 28, 26, 0},
	{106, 75, 74, 0},
	{169, 143, 142, 0},
	{222, 211, 209, 0},
	{251, 250, 249, 0},
	{250, 250, 248, 0},
	{216, 219, 205, 0},
	{145, 165, 138, 0},
	{71, 102, 80, 0},
	{24, 45, 43, 0},
	{3, 6, 8, 0},
	{0, 0, 0, 0},
	{15, 2, 24, 0},
	{72, 1, 67, 0},
	{134, 0, 71, 0},
	{194, 6, 52, 0},
	{238, 28, 41, 0},
	{254, 62, 62, 0},
	{240, 96, 104, 0},
	{198, 112, 131, 0},
	{138, 100, 120, 0},
	{76, 65, 75, 0},
	{26, 25, 26, 0},
	{1, 1, 1, 0},
	{7, 7, 7, 0},
	{39, 42, 42, 0},
	{80, 98, 86, 0},
	{120, 161, 107, 0},
	{189, 216, 109, 0},
	{249, 216, 85, 0},
	{252, 138, 49, 0},
	{225, 55, 18, 0},
	{173, 8, 2, 0},
	{111, 0, 0, 0},
	{52, 7, 2, 0},
	{12, 4, 1, 0},
	{0, 0, 0, 0},
	{19, 19, 8, 0},
	{51, 64, 38, 0},
	{93, 125, 95, 0},
	{163, 187, 178, 0},
	{225, 231, 233, 0},
	{254, 254, 254, 0},
	{241, 239, 244, 0},
	{203, 186, 205, 0},
	{147, 116, 137, 0},
	{84, 54, 64, 0},
	{32, 15, 17, 0},
	{3, 1, 1, 0},
	{4, 0, 0, 0},
	{36, 2, 7, 0},
	{90, 0, 34, 0},
	{153, 0, 111, 0},
	{178, 11, 209, 0},
	{107, 37, 246, 0},
	{73, 114, 253, 0},
	{103, 207, 230, 0},
	{111, 181, 152, 0},
	{91, 119, 91, 0},
	{56, 59, 53, 0},
	{16, 16, 15, 0},
	{0, 0, 0, 0},
	{15, 14, 14, 0},
	{57, 52, 51, 0},
	{117, 90, 90, 0},
	{179, 114, 111, 0},
	{228, 129, 104, 0},
	{253, 158, 74, 0},
	{247, 214, 38, 0},
	{151, 211, 12, 0},
	{27, 155, 0, 0},
	{0, 92, 37, 0},
	{2, 37, 37, 0},
	{0, 2, 5, 0},
	{1, 0, 2, 0},
	{25, 14, 30, 0},
	{82, 52, 77, 0},
	{144, 114, 128, 0},
	{203, 184, 188, 0},
	{243, 238, 238, 0},
	{254, 254, 254, 0},
	{235, 226, 227, 0},
	{189, 166, 172, 0},
	{128, 96, 114, 0},
	{66, 39, 65, 0},
	{15, 8, 20, 0},
	{0, 0, 0, 0},
	{1, 7, 11, 0},
	{2, 50, 42, 0},
	{0, 108, 26, 0},
	{56, 171, 2, 0},
	{188, 223, 17, 0},
	{252, 197, 47, 0},
	{249, 148, 83, 0},
	{217, 124, 108, 0},
	{163, 109, 108, 0},
	{100, 81, 81, 0},
	{44, 41, 40, 0},
	{8, 8, 7, 0},
	{1, 1, 1, 0},
	{25, 25, 24, 0},
	{68, 74, 64, 0},
	{99, 136, 105, 0},
	{112, 196, 174, 0},
	{97, 190, 239, 0},
	{63, 77, 254, 0},
	{129, 29, 239, 0},
	{189, 7, 196, 0},
	{136, 0, 85, 0},
	{74, 1, 23, 0},
	{25, 2, 4, 0},
	{1, 0, 0, 0},
	{8, 2, 3, 0},
	{44, 23, 26, 0},
	{100, 69, 82, 0},
	{163, 135, 157, 0},
	{214, 203, 217, 0},
	{248, 247, 249, 0},
	{250, 251, 252, 0},
	{211, 222, 223, 0},
	{144, 171, 157, 0},
	{79, 108, 76, 0},
	{42, 50, 27, 0},
	{11, 10, 4, 0},
	{0, 0, 0, 0},
	{20, 6, 2, 0},
	{66, 5, 1, 0},
	{128, 0, 0, 0},
	{189, 16, 5, 0},
	{235, 75, 25, 0},
	{254, 162, 58, 0},
	{243, 230, 93, 0},
	{168, 203, 111, 0},
	{106, 144, 102, 0},
	{69, 82, 75, 0},
	{28, 30, 30, 0},
	{2, 2, 2, 0},
	{5, 5, 5, 0},
	{37, 35, 37, 0},
	{92, 76, 88, 0},
	{155, 106, 127, 0},
	{211, 110, 127, 0},
	{247, 88, 92, 0},
	{253, 52, 53, 0},
	{228, 20, 41, 0},
	{179, 3, 58, 0},
	{117, 0, 74, 0},
	{55, 2, 57, 0},
	{8, 1, 15, 0},
	{0, 0, 0, 0},
	{6, 12, 16, 0},
	{34, 59, 52, 0},
	{87, 119, 92, 0},
	{167, 181, 156, 0},
	{229, 230, 220, 0},
	{253, 253, 253, 0},
	{246, 244, 242, 0},
	{209, 194, 192, 0},
	{153, 123, 123, 0},
	{90, 60, 59, 0},
	{36, 20, 18, 0},
	{4, 2, 1, 0},
	{3, 2, 0, 0},
	{26, 32, 2, 0},
	{26, 84, 1, 0},
	{0, 147, 37, 0},
	{9, 205, 174, 0},
	{34, 151, 244, 0},
	{74, 69, 254, 0},
	{175, 101, 233, 0},
	{187, 112, 184, 0},
	{125, 94, 112, 0},
	{64, 57, 59, 0},
	{19, 18, 18, 0},
	{0, 0, 0, 0},
	{12, 12, 12, 0},
	{52, 47, 48, 0},
	{111, 87, 98, 0},
	{173, 110, 165, 0},
	{186, 106, 225, 0},
	{104, 78, 252, 0},
	{42, 129, 249, 0},
	{13, 216, 212, 0},
	{1, 161, 63, 0},
	{18, 98, 0, 0},
	{30, 42, 2, 0},
	{7, 6, 1, 0},
	{1, 1, 0, 0},
	{26, 15, 12, 0},
	{76, 48, 47, 0},
	{138, 107, 107, 0},
	{198, 179, 177, 0},
	{240, 236, 234, 0},
	{254, 254, 254, 0},
	{238, 238, 231, 0},
	{184, 194, 172, 0},
	{102, 134, 103, 0},
	{44, 72, 61, 0},
	{10, 21, 24, 0},
	{0, 0, 0, 0},
	{3, 1, 8, 0},
	{39, 2, 45, 0},
	{102, 0, 73, 0},
	{165, 1, 63, 0},
	{219, 15, 43, 0},
	{250, 44, 47, 0},
	{251, 80, 82, 0},
	{222, 107, 122, 0},
	{169, 109, 131, 0},
	{106, 84, 100, 0},
	{48, 44, 48, 0},
	{10, 10, 10, 0},
	{0, 0, 0, 0},
	{20, 21, 21, 0},
	{60, 68, 65, 0},
	{96, 130, 98, 0},
	{151, 191, 112, 0},
	{234, 236, 99, 0},
	{254, 181, 67, 0},
	{242, 94, 32, 0},
	{201, 26, 8, 0},
	{142, 0, 0, 0},
	{80, 4, 1, 0},
	{29, 7, 2, 0},
	{2, 1, 0, 0},
	{6, 5, 2, 0},
	{34, 39, 20, 0},
	{70, 94, 63, 0},
	{128, 157, 138, 0},
	{197, 213, 212, 0},
	{245, 246, 248, 0},
	{252, 252, 253, 0},
	{223, 216, 227, 0},
	{177, 151, 174, 0},
	{115, 83, 99, 0},
	{55, 31, 36, 0},
	{14, 5, 6, 0},
	{0, 0, 0, 0},
	{17, 2, 3, 0},
	{61, 2, 16, 0},
	{121, 0, 66, 0},
	{183, 4, 171, 0},
	{146, 22, 231, 0},
	{66, 55, 254, 0},
	{90, 171, 245, 0},
	{111, 208, 196, 0},
	{104, 151, 118, 0},
	{77, 88, 73, 0},
	{34, 34, 32, 0},
	{4, 4, 4, 0},
	{3, 3, 3, 0},
	{33, 31, 31, 0},
	{86, 72, 72, 0},
	{149, 104, 104, 0},
	{206, 121, 111, 0},
	{245, 141, 91, 0},
	{254, 184, 56, 0},
	{219, 232, 23, 0},
	{84, 185, 4, 0},
	{0, 123, 13, 0},
	{1, 62, 44, 0},
	{2, 13, 18, 0},
	{0, 0, 0, 0},
	{8, 5, 13, 0},
	{52, 30, 54, 0},
	{113, 81, 102, 0},
	{175, 150, 158, 0},
	{226, 215, 216, 0},
	{252, 251, 251, 0},
	{248, 245, 245, 0},
	{214, 199, 201, 0},
	{159, 130, 142, 0},
	{96, 64, 88, 0},
	{36, 20, 40, 0},
	{3, 2, 6, 0},
	{0, 0, 1, 0},
	{2, 24, 28, 0},
	{1, 78, 42, 0},
	{6, 141, 0, 0},
	{120, 200, 8, 0},
	{241, 228, 31, 0},
	{254, 169, 66, 0},
	{237, 133, 98, 0},
	{192, 117, 112, 0},
	{132, 97, 97, 0},
	{70, 61, 61, 0},
	{22, 22, 21, 0},
	{0, 0, 0, 0},
	{9, 9, 9, 0},
	{46, 47, 43, 0},
	{86, 105, 83, 0},
	{109, 167, 135, 0},
	{107, 215, 220, 0},
	{81, 142, 251, 0},
	{88, 45, 250, 0},
	{165, 16, 220, 0},
	{167, 1, 137, 0},
	{104, 0, 47, 0},
	{47, 2, 10, 0},
	{9, 1, 1, 0},
	{0, 0, 0, 0},
	{23, 10, 11, 0},
	{70, 43, 51, 0},
	{132, 100, 119, 0},
	{192, 171, 193, 0},
	{233, 230, 237, 0},
	{254, 254, 254, 0},
	{235, 239, 241, 0},
	{179, 199, 195, 0},
	{109, 140, 115, 0},
	{59, 78, 48, 0},
	{25, 27, 12, 0},
	{1, 1, 0, 0},
	{6, 3, 1, 0},
	{41, 7, 2, 0},
	{96, 2, 0, 0},
	{159, 3, 1, 0},
	{215, 40, 13, 0},
	{248, 118, 41, 0},
	{252, 201, 77, 0},
	{209, 226, 105, 0},
	{133, 175, 110, 0},
	{88, 112, 93, 0},
	{48, 53, 52, 0},
	{12, 12, 13, 0},
	{0, 0, 0, 0},
	{18, 17, 18, 0},
	{63, 56, 63, 0},
	{124, 94, 111, 0},
	{185, 111, 132, 0},
	{232, 101, 113, 0},
	{254, 70, 71, 0},
	{244, 35, 43, 0},
	{206, 10, 48, 0},
	{148, 0, 68, 0},
	{85, 1, 71, 0},
	{24, 2, 33, 0},
	{1, 0, 3, 0},
	{1, 2, 4, 0},
	{17, 34, 35, 0},
	{58, 88, 71, 0},
	{125, 151, 121, 0},
	{202, 208, 191, 0},
	{245, 245, 241, 0},
	{254, 253, 253, 0},
	{231, 223, 221, 0},
	{183, 159, 158, 0},
	{121, 89, 89, 0},
	{60, 37, 35, 0},
	{16, 9, 7, 0},
	{0, 0, 0, 0},
	{14, 13, 1, 0},
	{33, 56, 2, 0},
	{3, 115, 0, 0},
	{3, 177, 98, 0},
	{20, 197, 227, 0},
	{51, 104, 253, 0},
	{135, 87, 247, 0},
	{193, 110, 212, 0},
	{157, 106, 144, 0},
	{94, 77, 84, 0},
	{39, 36, 36, 0},
	{5, 5, 5, 0},
	{2, 2, 2, 0},
	{29, 27, 28, 0},
	{80, 68, 72, 0},
	{143, 102, 129, 0},
	{194, 111, 201, 0},
	{156, 94, 242, 0},
	{59, 85, 254, 0},
	{26, 177, 235, 0},
	{5, 190, 131, 0},
	{0, 129, 12, 0},
	{32, 68, 1, 0},
	{20, 21, 2, 0},
	{0, 0, 0, 0},
	{10, 6, 4, 0},
	{49, 29, 26, 0},
	{107, 75, 75, 0},
	{170, 143, 143, 0},
	{222, 211, 209, 0},
	{251, 250, 250, 0},
	{250, 249, 248, 0},
	{215, 219, 204, 0},
	{144, 165, 137, 0},
	{70, 102, 80, 0},
	{24, 45, 43, 0},
	{3, 5, 8, 0},
	{0, 0, 1, 0},
	{15, 2, 24, 0},
	{72, 1, 67, 0},
	{134, 0, 71, 0},
	{194, 6, 52, 0},
	{238, 28, 41, 0},
	{254, 62, 62, 0},
	{240, 96, 104, 0},
	{197, 112, 131, 0},
	{138, 100, 119, 0},
	{75, 65, 74, 0},
	{25, 25, 26, 0},
	{1, 1, 1, 0},
	{7, 7, 7, 0},
	{39, 42, 42, 0},
	{80, 98, 86, 0},
	{120, 162, 108, 0},
	{190, 216, 109, 0},
	{249, 216, 84, 0},
	{252, 137, 48, 0},
	{224, 55, 18, 0},
	{173, 8, 2, 0},
	{110, 0, 0, 0},
	{52, 7, 2, 0},
	{12, 4, 1, 0},
	{0, 0, 0, 0},
	{19, 19, 8, 0},
	{51, 65, 38, 0},
	{94, 126, 96, 0},
	{164, 187, 179, 0},
	{225, 232, 234, 0},
	{254, 254, 254, 0},
	{241, 239, 243, 0},
	{202, 186, 204, 0},
	{146, 116, 136, 0},
	{83, 53, 64, 0},
	{31, 15, 17, 0},
	{3, 0, 1, 0},
	{4, 0, 0, 0},
	{36, 2, 7, 0},
	{90, 0, 35, 0},
	{153, 0, 112, 0},
	{178, 11, 210, 0},
	{106, 37, 246, 0},
	{73, 115, 253, 0},
	{103, 207, 229, 0},
	{111, 181, 151, 0},
	{91, 119, 91, 0},
	{56, 59, 52, 0},
	{15, 15, 15, 0},
	{0, 0, 0, 0},
	{15, 15, 14, 0},
	{58, 52, 51, 0},
	{117, 90, 90, 0},
	{180, 114, 111, 0},
	{229, 129, 104, 0},
	{253, 159, 74, 0},
	{247, 214, 38, 0},
	{150, 211, 11, 0},
	{27, 154, 0, 0},
	{0, 91, 37, 0},
	{2, 37, 37, 0},
	{0, 2, 5, 0},
	{1, 0, 2, 0},
	{25, 14, 31, 0},
	{82, 53, 78, 0},
	{145, 114, 129, 0},
	{203, 184, 188, 0},
	{243, 238, 238, 0},
	{254, 254, 254, 0},
	{234, 226, 227, 0},
	{188, 165, 171, 0},
	{127, 95, 113, 0},
	{66, 39, 65, 0},
	{15, 8, 20, 0},
	{0, 0, 0, 0},
	{1, 7, 11, 0},
	{2, 51, 42, 0},
	{0, 109, 26, 0},
	{57, 172, 2, 0},
	{189, 223, 17, 0},
	{252, 197, 48, 0},
	{249, 148, 84, 0},
	{217, 124, 108, 0},
	{163, 109, 108, 0},
	{100, 81, 81, 0},
	{43, 40, 40, 0},
	{7, 7, 7, 0},
	{1, 1, 1, 0},
	{25, 25, 24, 0},
	{68, 74, 64, 0},
	{99, 137, 105, 0},
	{112, 196, 175, 0},
	{96, 190, 239, 0},
	{63, 76, 254, 0},
	{129, 28, 239, 0},
	{190, 6, 195, 0},
	{135, 0, 84, 0},
	{73, 1, 23, 0},
	{25, 2, 4, 0},
	{1, 0, 0, 0},
	{8, 2, 3, 0},
	{44, 23, 27, 0},
	{101, 69, 83, 0},
	{164, 136, 158, 0},
	{214, 203, 218, 0},
	{248, 247, 250, 0},
	{250, 250, 251, 0},
	{210, 222, 223, 0},
	{144, 171, 157, 0},
	{79, 108, 76, 0},
	{41, 50, 27, 0},
	{11, 10, 4, 0},
	{0, 0, 0, 0},
	{20, 6, 2, 0},
	{67, 5, 1, 0},
	{128, 0, 0, 0},
	{189, 16, 5, 0},
	{235, 76, 25, 0},
	{254, 162, 59, 0},
	{242, 230, 93, 0},
	{167, 202, 111, 0},
	{105, 144, 102, 0},
	{69, 81, 75, 0},
	{28, 30, 30, 0},
	{2, 2, 2, 0},
	{5, 5, 5, 0},
	{37, 35, 38, 0},
	{92, 76, 89, 0},
	{156, 106, 127, 0},
	{212, 110, 127, 0},
	{247, 87, 92, 0},
	{253, 52, 53, 0},
	{228, 20, 41, 0},
	{179, 3, 58, 0},
	{116, 0, 74, 0},
	{54, 2, 57, 0},
	{7, 1, 14, 0},
	{0, 0, 0, 0},
	{6, 13, 16, 0},
	{34, 60, 53, 0},
	{88, 120, 92, 0},
	{167, 182, 157, 0},
	{229, 230, 220, 0},
	{253, 253, 253, 0},
	{246, 243, 242, 0},
	{209, 194, 192, 0},
	{152, 123, 122, 0},
	{89, 59, 59, 0},
	{35, 20, 17, 0},
	{4, 2, 1, 0},
	{3, 2, 0, 0},
	{26, 32, 2, 0},
	{26, 84, 1, 0},
	{0, 147, 38, 0},
	{9, 205, 175, 0},
	{34, 150, 244, 0},
	{75, 70, 254, 0},
	{176, 101, 233, 0},
	{186, 112, 183, 0},
	{125, 94, 111, 0},
	{64, 56, 58, 0},
	{19, 18, 18, 0},
	{0, 0, 0, 0},
	{12, 12, 12, 0},
	{53, 47, 48, 0},
	{111, 87, 99, 0},
	{174, 110, 165, 0},
	{186, 105, 225, 0},
	{103, 77, 252, 0},
	{41, 130, 249, 0},
	{13, 215, 211, 0},
	{1, 160, 62, 0},
	{19, 97, 0, 0},
	{30, 42, 2, 0},
	{7, 6, 1, 0},
	{1, 1, 0, 0},
	{27, 15, 12, 0},
	{77, 49, 48, 0},
	{139, 108, 108, 0},
	{198, 180, 178, 0},
	{240, 236, 234, 0},
	{254, 254, 254, 0},
	{238, 237, 230, 0},
	{183, 193, 172, 0},
	{102, 133, 102, 0},
	{43, 71, 60, 0},
	{10, 21, 23, 0},
	{0, 0, 0, 0},
	{4, 1, 9, 0},
	{39, 2, 46, 0},
	{103, 0, 73, 0},
	{166, 1, 63, 0},
	{219, 15, 43, 0},
	{250, 44, 48, 0},
	{251, 80, 83, 0},
	{221, 107, 122, 0},
	{169, 109, 131, 0},
	{106, 84, 99, 0},
	{47, 44, 48, 0},
	{10, 10, 10, 0},
	{0, 0, 0, 0},
	{21, 21, 22, 0},
	{60, 69, 66, 0},
	{97, 130, 98, 0},
	{152, 191, 112, 0},
	{235, 236, 99, 0},
	{254, 180, 66, 0},
	{241, 93, 31, 0},
	{201, 25, 8, 0},
	{142, 0, 0, 0},
	{79, 4, 1, 0},
	{28, 7, 2, 0},
	{2, 1, 0, 0},
	{6, 5, 2, 0},
	{34, 39, 20, 0},
	{70, 95, 63, 0},
	{129, 158, 139, 0},
	{197, 213, 212, 0},
	{245, 246, 248, 0},
	{252, 252, 253, 0},
	{223, 215, 227, 0},
	{176, 151, 173, 0},
	{114, 82, 98, 0},
	{55, 31, 36, 0},
	{13, 5, 5, 0},
	{0, 0, 0, 0},
	{17, 2, 3, 0},
	{61, 1, 17, 0},
	{122, 0, 66, 0},
	{184, 4, 172, 0},
	{146, 22, 231, 0},
	{66, 55, 254, 0},
	{90, 172, 245, 0},
	{111, 207, 195, 0},
	{104, 150, 117, 0},
	{77, 87, 73, 0},
	{33, 34, 32, 0},
	{4, 4, 4, 0},
	{3, 3, 3, 0},
	{34, 32, 31, 0},
	{87, 73, 72, 0},
	{149, 104, 104, 0},
	{207, 121, 111, 0},
	{245, 141, 90, 0},
	{254, 184, 55, 0},
	{218, 232, 23, 0},
	{83, 184, 4, 0},
	{0, 122, 13, 0},
	{1, 62, 44, 0},
	{2, 13, 17, 0},
	{0, 0, 0, 0},
	{8, 5, 13, 0},
	{53, 30, 54, 0},
	{113, 81, 102, 0},
	{176, 150, 159, 0},
	{226, 215, 216, 0},
	{253, 252, 252, 0},
	{248, 245, 245, 0},
	{214, 198, 201, 0},
	{158, 129, 141, 0},
	{95, 64, 88, 0},
	{36, 20, 40, 0},
	{3, 2, 6, 0},
	{0, 0, 2, 0},
	{2, 25, 28, 0},
	{1, 79, 42, 0},
	{7, 141, 0, 0},
	{121, 200, 8, 0},
	{241, 228, 31, 0},
	{254, 169, 66, 0},
	{236, 133, 99, 0},
	{192, 117, 112, 0},
	{131, 97, 97, 0},
	{69, 61, 60, 0},
	{22, 21, 21, 0},
	{0, 0, 0, 0},
	{9, 9, 9, 0},
	{46, 48, 43, 0},
	{86, 105, 84, 0},
	{109, 168, 136, 0},
	{107, 215, 221, 0},
	{81, 141, 251, 0},
	{89, 45, 250, 0},
	{166, 15, 220, 0},
	{166, 1, 136, 0},
	{103, 0, 46, 0},
	{46, 2, 10, 0},
	{9, 1, 1, 0},
	{0, 0, 0, 0},
	{23, 10, 11, 0},
	{71, 43, 51, 0},
	{133, 101, 120, 0},
	{193, 171, 193, 0},
	{234, 230, 237, 0},
	{254, 254, 254, 0},
	{235, 238, 240, 0},
	{178, 199, 194, 0},
	{108, 139, 114, 0},
	{59, 77, 48, 0},
	{25, 27, 12, 0},
	{1, 1, 0, 0},
	{6, 3, 1, 0},
	{41, 7, 2, 0},
	{97, 1, 0, 0},
	{160, 3, 1, 0},
	{215, 41, 13, 0},
	{248, 118, 41, 0},
	{252, 202, 77, 0},
	{208, 225, 105, 0},
	{133, 174, 110, 0},
	{88, 112, 93, 0},
	{48, 53, 52, 0},
	{12, 12, 12, 0},
	{0, 0, 0, 0},
	{18, 18, 18, 0},
	{63, 56, 63, 0},
	{124, 94, 112, 0},
	{186, 112, 132, 0},
	{233, 101, 112, 0},
	{254, 70, 70, 0},
	{244, 34, 43, 0},
	{206, 9, 48, 0},
	{148, 0, 68, 0},
	{85, 1, 71, 0},
	{24, 2, 32, 0},
	{1, 0, 3, 0},
	{1, 2, 4, 0},
	{17, 35, 35, 0},
	{58, 89, 71, 0},
	{126, 152, 122, 0},
	{203, 209, 191, 0},
	{246, 245, 242, 0},
	{254, 253, 253, 0},
	{230, 223, 221, 0},
	{182, 159, 158, 0},
	{120, 88, 88, 0},
	{60, 36, 34, 0},
	{16, 9, 6, 0},
	{0, 0, 0, 0},
	{14, 13, 1, 0},
	{33, 56, 2, 0},
	{3, 116, 0, 0},
	{3, 178, 100, 0},
	{20, 196, 228, 0},
	{52, 103, 253, 0},
	{136, 87, 247, 0},
	{193, 110, 212, 0},
	{156, 106, 143, 0},
	{93, 76, 83, 0},
	{38, 36, 36, 0},
	{5, 5, 5, 0},
	{2, 2, 2, 0},
	{29, 28, 28, 0},
	{81, 68, 73, 0},
	{143, 102, 130, 0},
	{194, 111, 202, 0},
	{155, 93, 242, 0},
	{59, 86, 254, 0},
	{25, 178, 235, 0},
	{5, 190, 129, 0},
	{0, 129, 11, 0},
}