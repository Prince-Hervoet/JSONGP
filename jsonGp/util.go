package jsongp

const (
	// 对象开始
	_OBJECT_BEGIN = 1
	// 对象结束
	_OBJECT_END = 1 << 1
	// 数组开始
	_ARRAY_BEGIN = 1 << 2
	// 数组结束
	_ARRAY_END = 1 << 3
	// string
	_STRING = 1 << 4
	// 冒号
	_COLON = 1 << 5
	// 逗号
	_COMMA = 1 << 6
	// boolean
	_BOOL   = 1 << 7
	_NULL   = 1 << 8
	_NUMBER = 1 << 9
	_KEY    = 1 << 10
	_VALUE  = 1 << 11
)

const (
	// 前导零
	_NUMBER_FIRST_ZERO = 1
	// 正常零
	_NUMBER_MULTI_ZERO = 1 << 1
	// 中划线
	_NUMBER_MINUS = 1 << 2
	// 小数点
	_NUMBER_POINT = 1 << 3
	// 正常数字，不包含0
	_NUMBER_COMMON      = 1 << 4
	_NUMBER_POINT_FALSE = 1 << 5
)
