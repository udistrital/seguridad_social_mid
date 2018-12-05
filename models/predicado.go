package models

type Predicado struct {
	Id               int    `orm:"column(id);pk;auto"`
	Nombre           string `orm:"column(nombre)"`
	Descripcion      string `orm:"column(descripcion)"`
	Estado           bool
	DescripcionHecho string
	Dominio          Dominio
	TipoPredicado    TipoPredicado
}

type Dominio struct {
	Id          int    `orm:"column(id);pk;auto"`
	Nombre      string `orm:"column(nombre)"`
	Descripcion string `orm:"column(descripcion);null"`
}

type TipoPredicado struct {
	Id     int    `orm:"column(id);pk;auto"`
	Nombre string `orm:"column(nombre);null"`
}
