package models

// DetallePreliquidacion modelo de titán
type DetallePreliquidacion struct {
	Id                   int                   `orm:"column(id);pk"`
	ValorCalculado       float64               `orm:"column(valor_calculado)"`
	NumeroContrato       string                `orm:"column(numero_contrato);null"`
	VigenciaContrato     int                   `orm:"column(vigencia_contrato);null"`
	Persona              int                   `orm:"column(persona)"`
	DiasLiquidados       float64               `orm:"column(dias_liquidados);null"`
	TipoPreliquidacion   *TipoPreliquidacion   `orm:"column(tipo_preliquidacion);rel(fk)"`
	Preliquidacion       *Preliquidacion       `orm:"column(preliquidacion);rel(fk)"`
	Concepto             *ConceptoNomina       `orm:"column(concepto);rel(fk)"`
	EstadoDisponibilidad *EstadoDisponibilidad `orm:"column(estado_disponibilidad);rel(fk)"`
	NombreCompleto       string
	Documento            string
}

// EstadoDisponibilidad modelo de titán
type EstadoDisponibilidad struct {
	Id                int     `orm:"column(id);pk"`
	Nombre            string  `orm:"column(nombre)"`
	Descripcion       string  `orm:"column(descripcion);null"`
	CodigoAbreviacion string  `orm:"column(codigo_abreviacion);null"`
	Activo            bool    `orm:"column(activo)"`
	NumeroOrden       float64 `orm:"column(numero_orden);null"`
}
