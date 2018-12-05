package models

type ConceptosIbc struct {
	Id               int
	Nombre           string
	Descripcion      string
	Estado           bool
	DescripcionHecho string
	Dominio          Dominio
	TipoPredicado    TipoPredicado
}
