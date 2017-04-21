package models

type Respuesta struct {
	Id           		int
	Nombre  		 		string
	NumDocumento 		float64
	Tipo_Descuneto  string
	Valor_descuento string
}
type FormatoPreliqu struct {
	//Contrato   *ContratoGeneral
	Respuesta *Respuesta
}
