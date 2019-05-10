package models

// Columna
type Columna struct {
	Valor    interface{}
	Longitud int
}

// PlanillaTipoE Estructura para planilla de activos
type PlanillaTipoE struct {
	TipoRegistro                    Columna
	TipoDocumento                   Columna
	NumeroIdentificacion            Columna
	TipoCotizante                   Columna
	SubTipoCotizante                Columna
	ExtranjeroNoPension             Columna
	PrimerApellido                  Columna
	SegundoApellido                 Columna
	PrimerNombre                    Columna
	SegundoNombre                   Columna
	CodigoFondoPension              Columna
	TrasladoPension                 Columna
	CodigoEps                       Columna
	TrasladoEps                     Columna
	CodigoCCF                       Columna
	DiasLaborados                   Columna
	DiasPension                     Columna
	DiasSalud                       Columna
	DiasArl                         Columna
	DiasCaja                        Columna
	SalarioBase                     Columna
	SalarioIntegral                 Columna
	Ibcension                       Columna
	IbcSalud                        Columna
	IbcArl                          Columna
	IbcCcf                          Columna
	TarifaPension                   Columna
	PagoPension                     Columna
	AportePension                   Columna
	TotalPension                    Columna
	FondoSolidaridad                Columna
	FondoSubsistencia               Columna
	NoRetenidoAportesVolunarios     Columna
	TarifaSalud                     Columna
	PagoSalud                       Columna
	ValorUpc                        Columna
	AutorizacionEnfermedadGeneral   Columna
	ValorIncapacidadGeneral         Columna
	AutotizacionLicenciaMarternidad Columna
	ValorLicenciaMaternidad         Columna
	TarifaArl                       Columna
	CentroTrabajo                   Columna
	PagoArl                         Columna
	TarifaCaja                      Columna
	PagoCaja                        Columna
	TarifaSena                      Columna
	PagoSena                        Columna
	TarifaIcbf                      Columna
	PagoIcbf                        Columna
	TarifaEsap                      Columna
	PagoEsap                        Columna
	TarifaMen                       Columna
	PagoMen                         Columna
	TipoDocumentoCotizantePrincipal Columna
	DocumentoCotizantePrincipal     Columna
	ExoneradoPagoSalud              Columna
	CodigoArl                       Columna
	ClaseRiesgo                     Columna
	IndicadorTarifaEspecialPension  Columna
	FechasNovedades                 Columna
	IbcOtrosParaFiscales            Columna
	HorasLaboradas                  Columna
	EspacioBlanco                   Columna
}
