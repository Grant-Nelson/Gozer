package typeTools

import "go/types"

// TODO: Implement a proper operation selections

func IsIndexable(t types.Type) bool {
	return IsSlice(t) || IsString(t) || IsMap(t) || IsArray(t) || IsArrayPtr(t)
}
