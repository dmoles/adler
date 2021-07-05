package markdown

const (
	titleKey = "Title"
	stylesKey = "Styles"
	scriptsKey = "Scripts"
)

// ------------------------------------------------------------
// Exported

type Metadata map[string]interface{}

func (m Metadata) Title() string {
	return m.getString(titleKey)
}

func (m Metadata) Styles() []Stylesheet {
	var styles []Stylesheet
	for _, src := range m.getStrings(stylesKey) {
		styles = append(styles, &stylesheet{src})
	}
	return styles
}

func (m Metadata) Scripts() []Script {
	var scripts []Script
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

func (m Metadata) getString(k string) string {
	v, ok := m[k]
	if !ok {
		return ""
	}
	s, ok := v.(string)
	return s
}

func (m Metadata) getMetadata(k string) (Metadata, bool) {
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
func (m Metadata) getStrings(k string) (strings []string) {

	for _, v := range m.getArray(k) {
		if s, ok := v.(string); ok {
			strings = append(strings, s)
		}
	}
	return strings
}

// getMetadatas given an array value for `k`, returns any Metadata objects
// in that array. The objects are guaranteed to appear in the same order
// as in the original array, but any non-object values are ignored. If the
// given key does not exist or is not an array value, an empty slice is returned.
func (m Metadata) getMetadatas(k string) (metadatas []Metadata) {
	for _, v := range m.getArray(k) {
		if m, ok := asMetadata(v); ok {
			metadatas = append(metadatas, m)
		}
	}
	return metadatas
}

func asMetadata(v interface{}) (m Metadata, ok bool) {
	m, ok = v.(map[string]interface{})
	if !ok {
		m, ok = v.(Metadata)
	}
	return m, ok
}

func (m Metadata) getArray(k string) []interface{} {
	v, ok := m[k]
	if !ok {
		return nil
	}
	a, _ := v.([]interface{})
	return a
}

//
//func (m Metadata) FindString(path... string) (string, error) {
//	var pathLen int = len(path)
//	if pathLen == 0 {
//		return "", nil
//	}
//
//	p1 := path[0]
//	val := m[p1]
//	if val == nil {
//		return "", nil
//	}
//
//	if pathLen == 1 {
//		if valStr, ok := val.(string); ok {
//			return valStr, nil
//		} else {
//			return "", fmt.Errorf("bad m: expected string value for key %#s, got %T: %#v", p1, val, val)
//		}
//	}
//	if valMap, ok := val.(map[string]interface{}); ok {
//		return valMap.FindString(path[1:]...)
//	}
//	return "", fmt.Errorf("bad m: expected map value for key %#s, got %R: %#v", p1, val, val)
//}
