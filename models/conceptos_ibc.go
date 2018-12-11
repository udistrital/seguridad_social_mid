package models

// ConceptosIbc estructura para crear el hecho concepto_ibc
type ConceptosIbc struct {
	Id               int
	Nombre           string
	Descripcion      string
	Estado           bool
	DescripcionHecho string
	Dominio          Dominio
	TipoPredicado    TipoPredicado
}

// ConceptoAporte estructura para crear el hecho de concepto_aporte
type ConceptoAporte struct {
	Id               int
	NombreAporte     string
	Porcentaje       string
	Nomina           string
	NombreResolucion string
	Resolucion       ResolucionAporte
}

// ResolucionAporte estructura para crear el hecho resolucion_aporte
type ResolucionAporte struct {
	Id           int
	NombreAporte string
	Resolucion   string
	Vigencia     string
	Estado       bool
}
