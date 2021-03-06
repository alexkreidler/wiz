// Code generated by go-enum
// DO NOT EDIT!

package packages

import (
	"fmt"
)

const (
	// PackageTypeDATA is a PackageType of type DATA
	PackageTypeDATA PackageType = iota + 1
	// PackageTypeCODE is a PackageType of type CODE
	PackageTypeCODE
)

const _PackageTypeName = "DATACODE"

var _PackageTypeMap = map[PackageType]string{
	1: _PackageTypeName[0:4],
	2: _PackageTypeName[4:8],
}

// String implements the Stringer interface.
func (x PackageType) String() string {
	if str, ok := _PackageTypeMap[x]; ok {
		return str
	}
	return fmt.Sprintf("PackageType(%d)", x)
}

var _PackageTypeValue = map[string]PackageType{
	_PackageTypeName[0:4]: 1,
	_PackageTypeName[4:8]: 2,
}

// ParsePackageType attempts to convert a string to a PackageType
func ParsePackageType(name string) (PackageType, error) {
	if x, ok := _PackageTypeValue[name]; ok {
		return x, nil
	}
	return PackageType(0), fmt.Errorf("%s is not a valid PackageType", name)
}

// MarshalText implements the text marshaller method
func (x *PackageType) MarshalText() ([]byte, error) {
	return []byte(x.String()), nil
}

// UnmarshalText implements the text unmarshaller method
func (x *PackageType) UnmarshalText(text []byte) error {
	name := string(text)
	tmp, err := ParsePackageType(name)
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}
