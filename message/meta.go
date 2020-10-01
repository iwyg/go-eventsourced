package message

type EventMetaData map[string]interface{}

func (em EventMetaData) merge(meta EventMetaData) EventMetaData {
	m := em
	if meta == nil {
		return m
	}

	for key, val := range meta {
		m[key] = val
	}

	return m
}
func (em EventMetaData) copy() EventMetaData {
	m := make(EventMetaData, len(em))
	for key, val := range em {
		m[key] = val
	}

	return m
}
