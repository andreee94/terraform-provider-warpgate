package provider

import (
	"fmt"
	"terraform-provider-warpgate/warpgate"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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

func ArrayOfUint8ToTerraformList(array []uint8) (result types.List) {
	// sort.Strings(array)

	result.ElemType = types.Int64Type

	for _, v := range array {
		result.Elems = append(result.Elems, types.Int64{Value: int64(v)})
	}
	return
}

func ArrayOfUint16ToTerraformList(array []uint16) (result types.List) {
	// sort.Strings(array)

	result.ElemType = types.Int64Type

	for _, v := range array {
		result.Elems = append(result.Elems, types.Int64{Value: int64(v)})
	}
	return
}

func TerraformListToArrayOfUint8(list types.List) (result []uint8) {
	for _, v := range list.Elems {
		result = append(result, uint8(v.(types.Int64).Value))
	}
	return
}

func TerraformListToArrayOfUint16(list types.List) (result []uint16) {
	for _, v := range list.Elems {
		result = append(result, uint16(v.(types.Int64).Value))
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

func testCheckFuncValidUUID(name string, key string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ms := s.RootModule()
		is, err := modulePrimaryInstanceState(ms, name)

		if err != nil {
			return nil
		}

		v, ok := is.Attributes[key]

		if !ok {
			return fmt.Errorf("%s: Attribute '%s' not found", name, key)
		}

		_, err = uuid.Parse(v)

		return err
	}
}

func modulePrimaryInstanceState(ms *terraform.ModuleState, name string) (*terraform.InstanceState, error) {
	rs, ok := ms.Resources[name]
	if !ok {
		return nil, fmt.Errorf("not found: %s in %s", name, ms.Path)
	}

	is := rs.Primary
	if is == nil {
		return nil, fmt.Errorf("no primary instance: %s in %s", name, ms.Path)
	}

	return is, nil
}
