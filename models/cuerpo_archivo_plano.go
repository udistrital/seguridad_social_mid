package models

type CuerpoArchivoPlano struct {
	TipoRegistro            string
	Secuencia               string
	TipoDocumento           string
	NumeroIdentificacion    string
	TipoCotizante           string
	SubtipoCotizante        string
	ExtjNoPension           string
	ColombianoExterior      string
	CodDeptUbicacionLaboral string
	CodMunUbicacionLaboral  string
	PrimerApellido          string
	SegundoApellido         string
	PrimerNombre            string
	SegundoNombre           string
	Ing                     string //Ingreso
	Ret                     string //Retiro
	Tde                     string //Traslado desde otra EPS o EOC
	Tae                     string //Traslado a otra EPS o EOC
	Tdp                     string //Traslado desde otra administradora de pensiones
	Tap                     string //Trasllado a otra adminisradora de pensiones
	Vsp                     string //Variación permanente de salario
	Correciones             string
	Vst                     string //Variación transitoria de salario
	Sln                     string //Suspención temporal del contrato o licencia no remunerada o comisión de servicios
	Ige                     string //Incapacidad temporal o por incacidad general
	Lma                     string //Licencia de maternidad o paternidad
	Vac                     string //Vacaciones, Licencia remunerada
	Avp                     string //Aporte voluntario
	Vct                     string //Variación de centros de trabajo
	Irl                     string //Días de incapacidad por accidente de trabajo o enfermedad laboral
	CodPensiones            string //Código de administradora de fondos de pensiones a la cual pertenece el afiliado
	CodNewPensiones         string //Codigo de la administradoara de fondo de pensiones a la cual se traslada el afiliado
	CodEps                  string //Código EPS o EOC a la cual pertenece el afiliado
	CodNewEps               string //Código EPS o EOc a la cuál se traslada el afiliado
	CodCcf                  string //Código de caja de compensación familiar a la cuál pertenece el afiliado
	DíasLaborados           string //Para Número de días cotizados a pensión, salud, arl y ccf
	SalarioBasico           string
	SalarioIntegral         string
	IbcPension              string
	IbcSalud                string
	IbcArl                  string
	IbcCcf                  string
	TarifaAportesPension    string
	CotOblgPen              string //Cotizacion obligatoria a pensiones
	AprtVoluntFondoPen      string //Aporte voluntario a fondo de pensiones
	TotalPension            string //Total cotización Sistema General de pensiones
	AportFondoSol           string //Aportes a fondo de solidaridad pensional, subcuenta de solidaridad
	AportFondoSols          string //Aportes a fondo de SegundoApellido pensional - subcuenta de subsistencia
	ValorNoRetVol           string //Valor no retenido por aportes voluntarios.
	TarifaSalud             string //Tarifa de aportes de Salud
	CotOblgSalud            string //Cotización obligatoria a Salud
	ValUpcAdicional         string //Valor de la Upc Adicional
	NumAutoIncGeneral       string //Número de autorización por incapacidad General
	ValIncapidadGeneral     string //Valor de la incapcidad general por enfermedad general
	NumAutoLicMat           string //Número de autorización por licencia de maternidad o partenidad
	ValLicMat               string //Valor de la licencia de maternidad
	TarifaArl               string //Tarifa de aportes a riesgos laborales
	CentTrabCt              string //Centro de Trabajo C.T.
	CotOblgArl              string //Cotización obligatoria a Sistema General de Riesgos Laborales
	TarifaCcf               string //Tarifa de aportes a caja de compensacion familiar
	ValorCcf                string //Valor caja de compensacion familiar
	TarifaSena              string //Tarifa de aportes SENA
	ValorSena               string //Valor aportes SENA
	TarifaIcbf              string //Tarifa ICBF
	ValorIcbf               string //Valor ICBF
	TarifaEsap              string //Tarifa aportes ESAP
	ValorEsap               string //Valor aportes ESAP
	TarifaAportesMen        string //Valor aportes ESAP
	ValorMen                string //Tarifa de aportes menor
	TipoDocu                string
}
