package sctp

// errorCauseUnrecognizedChunkType represents an SCTP error cause
type errorCauseUnrecognizedChunkType struct {
	errorCauseHeader
}

func (e *errorCauseUnrecognizedChunkType) marshal() ([]byte, error) {
	return e.errorCauseHeader.marshal()
}

func (e *errorCauseUnrecognizedChunkType) unmarshal(raw []byte) error {
	return e.errorCauseHeader.unmarshal(raw)
}
