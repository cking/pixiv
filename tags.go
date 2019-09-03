package pixiv

// Tags for objects
type Tags []Tag

// OriginalNames returns only original names
func (t Tags) OriginalNames() []string {
	l := make([]string, len(t))
	for i, v := range t {
		l[i] = v.Name
	}

	return l
}

// TranslatedNames returns only translated names, skipping empty ones
func (t Tags) TranslatedNames() []string {
	l := make([]string, 0, len(t))
	for _, v := range t {
		if v.TranslatedName == "" {
			continue
		}
		l = append(l, v.TranslatedName)
	}

	return l
}

// PreferTranslatedNames uses the translated name unless its empty
func (t Tags) PreferTranslatedNames() []string {
	l := make([]string, len(t))
	for i, v := range t {
		l[i] = v.TranslatedName
		if l[i] == "" {
			l[i] = v.Name
		}
	}

	return l
}

// Tag from the tags
type Tag struct {
	Name           string `json:"name,omitempty"`
	TranslatedName string `json:"translated_name,omitempty"`
}
