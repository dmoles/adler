package markdown

type Script struct {
	src string
	typ string
}

const scriptsKey = "Scripts"

func scriptFrom(md Metadata) (*Script, bool) {
	src, ok := md.GetString("src")
	if !ok {
		return nil, ok
	}
	typ, ok := md.GetString("type")
	if !ok {
		return nil, ok
	}
	return &Script{src: src, typ: typ}, true
}

func scriptsFrom(metadata Metadata) (scripts []*Script) {
	for _, md := range metadata.GetMetadatas(scriptsKey) {
		if s, ok := scriptFrom(md); ok {
			scripts = append(scripts, s)
		}
	}
	for _, src := range metadata.GetStrings(scriptsKey) {
		scripts = append(scripts, &Script{src: src})
	}

	return scripts
}
