package models

type PagosSeguridadSocial struct {
	Persona      int64
	SaludUd      float64
	SaludTotal   float64
	PensionUd    float64
	PensionTotal float64
	Arl          float64
	UpcAdicional []UpcAdicional
}
