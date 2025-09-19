package parser

type accessFlags int

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
)

func (af accessFlags) mask() uint16 {
	return 1 << af
}

func maskFromAccessFlags(afs []accessFlags) uint16 {
	m := uint16(0)
	for _, af := range afs {
		m |= af.mask()
	}
	return m
}

func accessFlagsFromMask(m uint16) []accessFlags {
	afs := []accessFlags{}
	for af := range 16 {
		af := accessFlags(af)
		if m&af.mask() != 0 {
			afs = append(afs, af)
		}
	}
	return afs
}
