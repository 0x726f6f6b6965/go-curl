// Code generated by "stringer -type=FileType ./filetype/file_type.go"; DO NOT EDIT.

package filetype

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[PEM-0]
	_ = x[DER-1]
	_ = x[ENG-2]
}

const _FileType_name = "PEMDERENG"

var _FileType_index = [...]uint8{0, 3, 6, 9}

func (i FileType) String() string {
	if i < 0 || i >= FileType(len(_FileType_index)-1) {
		return "FileType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _FileType_name[_FileType_index[i]:_FileType_index[i+1]]
}
