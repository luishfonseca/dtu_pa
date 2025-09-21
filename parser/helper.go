package parser

import (
	"fmt"
	"io"

	"github.com/luishfonseca/dtu_pa/data"
)

func parseMember(p *Parser, m data.MemberType) (*data.MemberInfo, error) {
	info := &data.MemberInfo{
		MemberType: m,
	}

	if err := p.readDecode(&info.AccessFlags); err != nil {
		return nil, err
	}

	var cpIndex uint16
	if err := p.readDecode(&cpIndex); err != nil {
		return nil, err
	}

	info.Name = *p.class.ConstantPool[cpIndex-1].ConstantUtf8()

	if err := p.readDecode(&cpIndex); err != nil {
		return nil, err
	}

	info.Descriptor = *p.class.ConstantPool[cpIndex-1].ConstantUtf8()

	var n uint16
	if err := p.readDecode(&n); err != nil {
		return nil, err
	}

	info.Attributes = make(map[data.Tag]*data.AttributeHandle)
	for range n {
		if attr, err := parseAttribute(p); err != nil {
			return nil, err
		} else {
			info.Attributes[attr.AttributeTag] = attr
		}
	}

	return info, nil
}

func parseAttribute(p *Parser) (*data.AttributeHandle, error) {
	var cpIndex uint16
	if err := p.readDecode(&cpIndex); err != nil {
		return nil, err
	}

	name := p.class.ConstantPool[cpIndex-1].ConstantUtf8().Value

	var tag data.Tag
	switch name {
	case "Code":
		tag = data.ATTR_CODE
	case "RuntimeVisibleAnnotations":
		tag = data.ATTR_RUNTIME_VISIBLE_ANNOTATIONS
	case "SourceFile":
		tag = data.ATTR_SOURCE_FILE
	case "InnerClasses":
		tag = data.ATTR_INNER_CLASSES
	case "LineNumberTable":
		tag = data.ATTR_LINE_NUMBER_TABLE
	case "LocalVariableTable":
		tag = data.ATTR_LOCAL_VARIABLE_TABLE
	case "StackMapTable":
		tag = data.ATTR_STACK_MAP_TABLE
	default:
		return nil, fmt.Errorf("unknown attribute name: %s", name)
	}

	var size uint32
	if err := p.readDecode(&size); err != nil {
		return nil, err
	}

	begin, err := p.input.Seek(0, io.SeekCurrent) // Mark current position
	if err != nil {
		return nil, err
	}

	_, err = p.input.Seek(int64(size), io.SeekCurrent) // Skip attribute content
	if err != nil {
		return nil, err
	}

	return &data.AttributeHandle{
		AttributeTag: tag,
		Begin:        begin,
	}, nil
}
