package markdown

const (
	titleKey   = "Title"
	stylesKey  = "Stylesheets"
	scriptsKey = "Scripts"
)

// ------------------------------------------------------------
// Exported

type metadata map[string]interface{}

func (m metadata) Title() string {
	return m.getString(titleKey)
}

func (m metadata) Styles() []*stylesheet {
	var styles []*stylesheet
	for _, src := range m.getStrings(stylesKey) {
		styles = append(styles, &stylesheet{src})
	}
	return styles
}

func (m metadata) Scripts() []*script {
	var scripts []*script
	for _, md := range m.getMetadatas(scriptsKey) {
		if s, ok := scriptFrom(md); ok {
			scripts = append(scripts, s)
		}
	}
	for _, src := range m.getStrings(scriptsKey) {
		scripts = append(scripts, &script{src: src})
	}

	return scripts
}

// ------------------------------------------------------------
// Unexported

func (m metadata) getString(k string) string {
	v, ok := m[k]
	if !ok {
		return ""
	}
	s, ok := v.(string)
	return s
}

func (m metadata) getMetadata(k string) (metadata, bool) {
	v, ok := m[k]
	if !ok {
		return nil, false
	}
	m1, ok := v.(map[string]interface{})
	return m1, ok
}

// getStrings given an array value for `k`, returns any strings in that
// array. The strings are guaranteed to appear in the same order as in the
// original array, but any non-string values are ignored. If the given key
// does not exist or is not an array value, an empty slice is returned.
func (m metadata) getStrings(k string) (strings []string) {
	for _, v := range m.getArray(k) {
		if s, ok := v.(string); ok {
			strings = append(strings, s)
		}
	}
	return strings
}

// getMetadatas given an array value for `k`, returns any metadata objects
// in that array. The objects are guaranteed to appear in the same order
// as in the original array, but any non-object values are ignored. If the
// given key does not exist or is not an array value, an empty slice is returned.
func (m metadata) getMetadatas(k string) (metadatas []metadata) {
	for _, v := range m.getArray(k) {
		if m, ok := asMetadata(v); ok {
			metadatas = append(metadatas, m)
		}
	}
	return metadatas
}

func asMetadata(v interface{}) (m metadata, ok bool) {
	var h map[interface{}]interface{}
	h, ok = v.(map[interface{}]interface{})
	if ok {
		m = make(metadata)
		for k, v := range h {
			if ks, ok := k.(string); ok {
				m[ks] = v
			}
		}
	} else {
		m, ok = v.(metadata)
	}
	return m, ok
}

func (m metadata) getArray(k string) []interface{} {
	v, ok := m[k]
	if !ok {
		return nil
	}
	a, _ := v.([]interface{})
	return a
}
