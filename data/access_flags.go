package data

import "fmt"

type AccessFlags uint16

const (
	ACC_PUBLIC AccessFlags = iota
	ACC_PRIVATE
	ACC_PROTECTED
	ACC_STATIC
	ACC_FINAL
	ACC_SUPER_OR_SYNCHRONIZED
	ACC_VOLATILE_OR_BRIDGE
	ACC_TRANSIENT_OR_VARARGS
	ACC_NATIVE
	ACC_INTERFACE
	ACC_ABSTRACT
	ACC_STRICT
	ACC_SYNTHETIC
	ACC_ANNOTATION
	ACC_ENUM
	ACC_MANDATED
	ACCESS_FLAGS_COUNT
)

func (af AccessFlags) mask() AccessFlags {
	return 1 << af
}

func (af AccessFlags) decompose() []AccessFlags {
	afs := []AccessFlags{}
	for m := range ACCESS_FLAGS_COUNT {
		if af&m.mask() != 0 {
			afs = append(afs, m)
		}
	}
	return afs
}

func (af AccessFlags) String() string {
	encoded := "<Flags: ["
	for i, m := range af.decompose() {
		if i > 0 {
			encoded += " "
		}
		encoded += fmt.Sprintf("0x%04X", uint16(m.mask()))
	}
	return encoded + "]>"
}
