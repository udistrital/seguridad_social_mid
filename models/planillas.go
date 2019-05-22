package models

// Columna
type Columna struct {
	Valor    interface{}
	Longitud int
}

// CabeceraPlanilla Estructura para la cabecera de todas las planillas
type CabeceraPlanilla struct {
	Codigo           Columna
	NombreProveedor  Columna
	NitProveedor     Columna
	TipoPlanilla     Columna
	Sucursal         Columna
	CodigoArl        Columna
	PeriodoSalud     Columna
	PeriodoPension   Columna
	CantidadPersonas Columna
	TotalNomina      Columna
	CodigoProveedor  Columna
	CodigoOperador   Columna
}

// PlanillaTipoE Estructura para planilla de activos
type PlanillaTipoE struct {
	TipoRegistro                    Columna
	TipoDocumento                   Columna
	NumeroIdentificacion            Columna
	TipoCotizante                   Columna
	SubTipoCotizante                Columna
	ExtranjeroNoPension             Columna
	ColombianoExterior              Columna
	CodigoDepartamento              Columna
	CodigoMunicipio                 Columna
	PrimerApellido                  Columna
	SegundoApellido                 Columna
	PrimerNombre                    Columna
	SegundoNombre                   Columna
	NovIng                          Columna // Ingreso
	NovRet                          Columna // Retiro
	NovTde                          Columna // Traslado desde otra EPS o EOC
	NovTae                          Columna // Traslado a otra EPS o EOC
	NovTdp                          Columna // Traslado desde otra administadora de pensiones
	NovTap                          Columna // Traslado a otra administradora de pensiones
	NovVsp                          Columna // Variación permanente de salario
	NovCorrecciones                 Columna // Corecciones
	NovVst                          Columna // Varaición transitoria de salario
	NovSln                          Columna // Suspención temporal del contrato de trabajo
	NovIge                          Columna // Inapacidad temporal por enfermedad general
	NovLma                          Columna // Licencia de maternidad o de paternidad
	NovVac                          Columna // Vacaciones, Licencia remunerada
	NovAvp                          Columna // Aporte voluntario
	NovVct                          Columna // Variación centros de trabajo
	NavIrl                          Columna // Días de incapacidad por accidente de trabajo o enfermedad laboral
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
	IbcPension                      Columna
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
	FechaIngreso                    Columna
	FechaRetiro                     Columna
	FechaInicioVsp                  Columna
	FechaInicioSuspencion           Columna
	FechaFinSuspencion              Columna
	FechaInicioIge                  Columna
	FechaFinIge                     Columna
	FechaInicioLma                  Columna
	FechaFinLma                     Columna
	FechaInicioVac                  Columna
	FechaFinVac                     Columna
	FechaInicioVct                  Columna
	FechaFinVct                     Columna
	FechaInicioIrl                  Columna
	FechaFinIrl                     Columna
	IbcOtrosParaFiscales            Columna
	HorasLaboradas                  Columna
	EspacioBlanco                   Columna
}
