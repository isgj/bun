package tagparser

import (
	"strings"
)

type Tag struct {
	Name    string
	Options map[string]string
}

func (t Tag) HasOption(name string) bool {
	_, ok := t.Options[name]
	return ok
}

func Parse(s string) Tag {
	p := parser{
		s: s,
	}
	p.parse()
	return p.tag
}

type parser struct {
	s string
	i int

	tag      Tag
	seenName bool // for empty names
}

func (p *parser) setName(name string) {
	if p.seenName {
		p.addOption(name, "")
	} else {
		p.seenName = true
		p.tag.Name = name
	}
}

func (p *parser) addOption(key, value string) {
	p.seenName = true
	if key == "" {
		return
	}
	if p.tag.Options == nil {
		p.tag.Options = make(map[string]string)
	}
	p.tag.Options[key] = value
}

func (p *parser) parse() {
	for p.valid() {
		p.parseKeyValue()
		if p.peek() == ',' {
			p.i++
		}
	}
}

func (p *parser) parseKeyValue() {
	start := p.i

	for p.valid() {
		switch c := p.read(); c {
		case ',':
			key := p.s[start : p.i-1]
			p.setName(key)
			return
		case ':':
			key := p.s[start : p.i-1]
			value := p.parseValue()
			p.addOption(key, value)
			return
		case '"':
			key := p.parseQuotedValue()
			p.setName(key)
			return
		}
	}

	key := p.s[start:p.i]
	p.setName(key)
}

func (p *parser) parseValue() string {
	start := p.i

	for p.valid() {
		switch c := p.read(); c {
		case '"':
			return p.parseQuotedValue()
		case ',':
			return p.s[start : p.i-1]
		}
	}

	if p.i == start {
		return ""
	}
	return p.s[start:p.i]
}

func (p *parser) parseQuotedValue() string {
	if i := strings.IndexByte(p.s[p.i:], '"'); i >= 0 && p.s[p.i+i-1] != '\\' {
		s := p.s[p.i : p.i+i]
		p.i += i + 1
		return s
	}

	b := make([]byte, 0, 16)

	for p.valid() {
		switch c := p.read(); c {
		case '\\':
			b = append(b, p.read())
		case '"':
			return string(b)
		default:
			b = append(b, c)
		}
	}

	return ""
}

func (p *parser) valid() bool {
	return p.i < len(p.s)
}

func (p *parser) read() byte {
	if !p.valid() {
		return 0
	}
	c := p.s[p.i]
	p.i++
	return c
}

func (p *parser) peek() byte {
	if !p.valid() {
		return 0
	}
	c := p.s[p.i]
	return c
}
