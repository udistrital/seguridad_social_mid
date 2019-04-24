package models

// PlanillaTipoE Estructura para planilla de activos
type PlanillaTipoE struct {
	TipoRegistro         string
	Secuencia            int
	TipoDocumento        string
	NumeroIdentificacion string
	TipoCotizante        int
	SubTipoCotizante     int
	ExtranjeroNoPension  string
	PrimerApellido       string
	SegundoApellido      string
	PrimerNombre         string
	SegundoNombre        string
	CodigoFondoPension   string
	TrasladoPension      string
	CodigoEps            string
	TrasladoEps          string
	CodigoCCF            string
	DiasLaborados        int
	HorasLaboradas       int
	DiasPension          string
	DiasSalud            string
	DiasArl              int
	DiasCaja             string
	PagoCaja             string
	SalarioBase          int
}
