package data

type AttributeCode struct {
	baseData
}

func (a *AttributeCode) AttributeCode() *AttributeCode {
	return a
}

func (*AttributeCode) Tag() Tag {
	return ATTR_CODE
}

type AttributeSourceFile struct {
	baseData
}

func (a *AttributeSourceFile) AttributeSourceFile() *AttributeSourceFile {
	return a
}

func (*AttributeSourceFile) Tag() Tag {
	return ATTR_SOURCE_FILE
}

type AttributeRuntimeVisibleAnnotations struct {
	baseData
}

func (a *AttributeRuntimeVisibleAnnotations) AttributeRuntimeVisibleAnnotations() *AttributeRuntimeVisibleAnnotations {
	return a
}

func (*AttributeRuntimeVisibleAnnotations) Tag() Tag {
	return ATTR_RUNTIME_VISIBLE_ANNOTATIONS
}

type AttributeInnerClasses struct {
	baseData
}

func (a *AttributeInnerClasses) AttributeInnerClasses() *AttributeInnerClasses {
	return a
}

func (*AttributeInnerClasses) Tag() Tag {
	return ATTR_INNER_CLASSES
}
