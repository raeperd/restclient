// Code generated by "stringer -type=ErrorKind"; DO NOT EDIT.

package restclient

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ErrURL-0]
	_ = x[ErrRequest-1]
	_ = x[ErrTransport-2]
	_ = x[ErrValidator-3]
	_ = x[ErrHandler-4]
}

const _ErrorKind_name = "ErrURLErrRequestErrTransportErrValidatorErrHandler"

var _ErrorKind_index = [...]uint8{0, 6, 16, 28, 40, 50}

func (i ErrorKind) String() string {
	if i < 0 || i >= ErrorKind(len(_ErrorKind_index)-1) {
		return "ErrorKind(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ErrorKind_name[_ErrorKind_index[i]:_ErrorKind_index[i+1]]
}
