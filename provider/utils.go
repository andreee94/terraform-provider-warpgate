package provider

import (
	"terraform-provider-warpgate/warpgate"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

//	func SumIntsOrFloats[K comparable, V int64 | float64](m map[K]V) V {
//	    var s V
//	    for _, v := range m {
//	        s += v
//	    }
//	    return s
//	}W

func ArrayIntersection[K comparable](a []K, b []K) (inAAndB []K, inAButNotB []K, inBButNotA []K) {
	m := make(map[K]uint8)
	for _, k := range a {
		m[k] |= (1 << 0)
	}
	for _, k := range b {
		m[k] |= (1 << 1)
	}

	// var inAAndB, inAButNotB, inBButNotA []K

	for k, v := range m {
		a := v&(1<<0) != 0
		b := v&(1<<1) != 0
		switch {
		case a && b:
			inAAndB = append(inAAndB, k)
		case a && !b:
			inAButNotB = append(inAButNotB, k)
		case !a && b:
			inBButNotA = append(inBButNotA, k)
		}
	}
	return
}

func ArrayOfStringToTerraformSet(array []string) (result types.Set) {
	// sort.Strings(array)

	result.ElemType = types.StringType

	for _, v := range array {
		result.Elems = append(result.Elems, types.String{Value: v})
	}
	return
}

func ArrayOfRolesToTerraformSet(array []warpgate.Role) (result types.Set) {
	array_string := []string{}

	for _, v := range array {
		array_string = append(array_string, v.Id.String())
	}

	return ArrayOfStringToTerraformSet(array_string)
	// sort.Strings(array)

	// result.ElemType = types.StringType

	// for _, v := range array {
	// 	result.Elems = append(result.Elems, types.String{Value: v.Id.String()})
	// }
	// return
}

// func GetArraySortedToString(list types.List) (result []string) {

// 	array_string := []string{}

// 	for _, v := range list.Elems {
// 		array_string = append(array_string, v.String())
// 	}
// 	sort.Strings(array_string)

// 	return
// }

func If[T any](cond bool, vtrue, vfalse T) T {
	if cond {
		return vtrue
	}
	return vfalse
}
