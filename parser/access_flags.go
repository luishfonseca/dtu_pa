package parser

import "fmt"

type accessFlags uint16

const (
	ACC_PUBLIC accessFlags = iota
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

func (af accessFlags) mask() accessFlags {
	return 1 << af
}

func (af accessFlags) decompose() []accessFlags {
	afs := []accessFlags{}
	for m := range ACCESS_FLAGS_COUNT {
		if af&m.mask() != 0 {
			afs = append(afs, m)
		}
	}
	return afs
}

func (af accessFlags) String() string {
	encoded := "<Flags:"
	for _, m := range af.decompose() {
		encoded += fmt.Sprintf(" 0x%04X", uint16(m.mask()))
	}
	return encoded + ">"
}
