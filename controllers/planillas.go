package controllers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/udistrital/ss_mid_api/models"
)

// PlanillasController operations for Planillas
type PlanillasController struct {
	beego.Controller
}

// URLMapping ...
func (c *PlanillasController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
	c.Mapping("GenerarPlanillaActivos", c.GenerarPlanillaActivos)
	//c.Mapping("GenerarPlanillaPensionados", c.GenerarPlanillaPensionados)
	//c.Mapping("GenerarPlanillaN", c.GenerarPlanillaN)
}

var formatoFecha = "2006-01-02"
var fila = ""

//Variables para cada una de las novedades y sus días validos
var ingreso = false
var fechaIngreso = ""

var retiro = false
var fechaRetiro = ""

var trasladoPensiones = false
var trasladoEps = false
var exterior = false
var suspencionContrato = false
var fechaInicioSuspencion = ""
var fechaFinSuspencion = ""

var licenciaNoRem = false
var comisionServicios = false
var incapacidadGeneral = false
var fechaInicioIge = ""
var fechaFinIge = ""

var licenciaMaternidad = false
var fechaInicioLma = ""
var fechaFinLma = ""

var vacaciones = false
var licenciaRem = false
var fechaInicioVac = ""
var fechaFinVac = ""

var aporteVoluntario = false
var variacionCentroTrabajo = false
var fechaInicioVct = ""
var fechaFinVct = ""

var diasIncapcidadLaboral = 0
var fechaInicioIrl = ""
var fechaFinIrl = ""

var fechaInicioVsp = ""

func (c *PlanillasController) GenerarPlanillaActivos() {
	idStr := c.Ctx.Input.Param(":id")
	idDescSegSocial, _ := strconv.Atoi(idStr)
	var personaNatural []models.InformacionPersonaNatural
	var proveedor []models.InformacionProveedor
	var contratosGeneral []models.ContratoGeneral
	var contratoGeneral []models.ContratoGeneral
	var upc []models.UpcAdicional
	var detallePreliquidacion []models.DetallePreliquidacion
	var detLiq []models.DetallePreliquidacion
	var conceptos []models.Concepto
	var conceptosSeguridadSocial []models.Concepto
	var errStrings []string
	tipoRegistro := "02"

	// Se obtienen todos los conceptos de seguridad social en tabla conceptos de titan
	errConceptosSs := getJson("http://"+beego.AppConfig.String("titanServicio")+
		"/concepto?limit=0&query=Naturaleza:seguridad_social", &conceptosSeguridadSocial)
	if errConceptosSs != nil {
		fmt.Println("errConceptosSs: ", errConceptosSs)
	}

	errLiquidacion := getJson("http://"+beego.AppConfig.String("titanServicio")+
		"/detalle_liquidacion?limit=-1", &detallePreliquidacion)
	if errLiquidacion != nil {
		errStrings = append(errStrings, errLiquidacion.Error())
	}

	errContratosGeneral := getJson("http://"+beego.AppConfig.String("agoraServicio")+
		"/contrato_general?limit=0", &contratosGeneral)
	if errContratosGeneral != nil {
		errStrings = append(errStrings, errContratosGeneral.Error())
	}

	errConceptos := getJson("http://"+beego.AppConfig.String("titanServicio")+
		"/concepto?limit=0", &conceptos)
	if errConceptos != nil {
		errStrings = append(errStrings, errConceptos.Error())
	}

	errUpc := getJson("http://"+beego.AppConfig.String("seguridadSocialService")+
		"/upc_adicional?limit=0", &upc)
	if errUpc != nil {
		errStrings = append(errStrings, errUpc.Error())
	}

	fmt.Println("errStrings: ", errStrings)
	if errStrings == nil {
		secuencia := 1
		for i := 0; i < len(detallePreliquidacion); i++ {
			for j := 0; j < len(contratosGeneral); j++ {
				if detallePreliquidacion[i].NumeroContrato == contratosGeneral[j].Id {
					if strings.Contains(fila, strconv.Itoa(int(contratosGeneral[j].Contratista))) {
						break
					} else {

						var ibcLiquidado = 0
						var pagoSalud = 0
						var pagoPension = 0
						var pagoArl = 0
						var pagoCaja = 0
						var pagoIcbf = 0

						fila += formatoDato(tipoRegistro, 2)                     //Tipo Registro
						fila += formatoDato(completarSecuencia(secuencia, 5), 5) //Secuencia

						fila += formatoDato("CC", 2)                                                //Tip de documento del cotizante
						fila += formatoDato(strconv.Itoa(int(contratosGeneral[j].Contratista)), 16) //Número de identificación del cotizante
						fila += formatoDato(completarSecuencia(1, 2), 2)                            //Tipo Cotizante
						fila += formatoDato(completarSecuencia(1, 2), 2)                            //Subtipo de Cotizante
						fila += formatoDato("", 1)                                                  //Extranjero no obligado a cotizar pensión

						errPersonaNatural := getJson("http://"+beego.AppConfig.String("agoraServicio")+
							"/informacion_persona_natural?limit=1&query=Id:"+strconv.Itoa(contratosGeneral[j].Contratista), &personaNatural)

						/*errPersonaNatural := getJson("http://"+beego.AppConfig.String("agoraServicio")+
						"/informacion_persona_natural"+
						"?limit=1"+
						"&query=Id:"+strconv.FormatFloat(proveedores[j].NumDocumento, 'E', -1, 64), &personaNatural)*/

						if errPersonaNatural != nil {
							fmt.Println("errPersonaNatural: ", errPersonaNatural)
						}

						if exterior {
							fila += formatoDato("X", 1) //Colombiano en el exterior
							fila += formatoDato(" ", 2) //Código del departamento de la ubicación laboral
							fila += formatoDato(" ", 3) //Código del municipio de ubicación laboral
						} else {
							fila += formatoDato("", 1)    //Colombiano en el exterior
							fila += formatoDato("11", 2)  //Código del departamento de la ubicación laboral
							fila += formatoDato("001", 3) //Código del municipio de ubicación laboral
						}

						fila += formatoDato(personaNatural[0].PrimerApellido, 20)  //Primer apellido
						fila += formatoDato(personaNatural[0].SegundoApellido, 30) //Segundo apellido
						fila += formatoDato(personaNatural[0].PrimerNombre, 20)    //Primer nombre
						fila += formatoDato(personaNatural[0].SegundoNombre, 30)   //Segundo nombre

						// --AQUÍ VA LA FUNCIÓN DE LAS NOVEDADES!--  //
						errNumContrato := getJson("http://"+beego.AppConfig.String("argoServicio")+"/contrato_general"+
							"?query=Id:"+detallePreliquidacion[i].NumeroContrato, &contratoGeneral)

						if errNumContrato != nil {
							fmt.Println("errNumContrato: ", errNumContrato)
						}

						errProveedor := getJson("http://"+beego.AppConfig.String("agoraServicio")+
							"/informacion_proveedor?limit=1&query:NumDocumento:"+strconv.Itoa(contratosGeneral[j].Contratista), &proveedor)

						if errProveedor != nil {
							fmt.Println("errProveedor: ", errProveedor)
						}

						/*/errPersonaNatural := getJson("http://"+beego.AppConfig.String("agoraServicio")+
							"/informacion_persona_natural?limit=1&query=Id:"+strconv.Itoa(contratosGeneral[j].Contratista), &personaNatural)

						if errPersonaNatural != nil {
							fmt.Println("errPersonaNatural: ", errPersonaNatural)
						}*/

						establecerNovedades(strconv.Itoa(proveedor[0].Id))

						//Código de la administradora de fondo de pensiones a la cual pertenece el afiliado
						fila += formatoDato("231001", 6)

						//Código de la admnistradora de pensiones a la cual se traslada el afiliado
						// Si hay un translado, debe aparecer el nuevo código, de lo contrario será un campo vació
						if trasladoPensiones {
							fila += formatoDato("230301", 6)
						} else {
							fila += formatoDato(" ", 6)
						}

						//Código EPS o EOC a la cual pertenece el afiliado
						fila += formatoDato("EPS010", 6)

						//Código EPS o EOC a la cual se traslada el afiliado
						// Si hay un translado, debe aparecer el nuevo código, de lo contrario será un campo vació
						if trasladoEps {
							fila += formatoDato("EPS012", 6)
						} else {
							fila += formatoDato(" ", 6)
						}

						//Código CCF a la cual pertenece el afiliado
						fila += formatoDato("CCF04", 6)

						errDiasLiquidados := getJson("http://"+beego.AppConfig.String("titanServicio")+
							"/detalle_preliquidacion?limit=0"+
							"&query=Concepto.NombreConcepto:ibc_novedad,NumeroContrato:"+detallePreliquidacion[i].NumeroContrato, &detLiq)
						if errDiasLiquidados != nil {
							fmt.Println("errDiasLiquidados: ", errDiasLiquidados)
						}
						diasCotizados := int(detLiq[0].DiasLiquidados)

						if ingreso || retiro {
							fila += formatoDato(completarSecuencia(diasCotizados, 2), 2) //Número de días cotizados a pensión
							fila += formatoDato(completarSecuencia(diasCotizados, 2), 2) //Número de días cotizados a salud
						} else {
							fila += formatoDato("30", 2) //Número de días cotizados a pensión
							fila += formatoDato("30", 2) //Número de días cotizados a salud
						}

						if incapacidadGeneral || licenciaMaternidad || vacaciones || diasIncapcidadLaboral > 0 {
							fila += formatoDato(completarSecuencia(diasCotizados, 2), 2) //Número de días cotizados a ARL
							fila += formatoDato(completarSecuencia(diasCotizados, 2), 2) //Número de días cotizados a CCF
						} else {
							fila += formatoDato("30", 2) //Número de días cotizados a ARL
							fila += formatoDato("30", 2) //Número de días cotizados a CCF
						}

						errSalarioBasico := getJson("http://"+beego.AppConfig.String("titanServicio")+
							"/detalle_preliquidacion?limit=0"+
							"&query=Concepto.NombreConcepto:salarioBase,NumeroContrato:"+detallePreliquidacion[i].NumeroContrato, &detLiq)
						if errSalarioBasico != nil {
							fmt.Println("errSalarioBasico: ", errSalarioBasico)
						} else {
							salarioBase := strconv.FormatFloat(detLiq[0].ValorCalculado, 'E', -1, 64)
							fila += formatoDato(salarioBase, 9) //Salario básico
						}

						fila += formatoDato("", 1) //Salario integral

						errSoloLiquidado := getJson("http://"+beego.AppConfig.String("titanServicio")+
							"/detalle_preliquidacion?limit=0"+
							"&query=Concepto.NombreConcepto:ibc_liquidado,NumeroContrato:"+detallePreliquidacion[i].NumeroContrato, &detLiq)
						if errSoloLiquidado != nil {
							fmt.Println("errSoloLiquidado: ", errSoloLiquidado)
						} else {
							ibcLiquidado = int(detLiq[0].ValorCalculado)
							fila += formatoDato(completarSecuencia(ibcLiquidado, 9), 9) //IBC pensión
							fila += formatoDato(completarSecuencia(ibcLiquidado, 9), 9) //IBC salud
							fila += formatoDato(completarSecuencia(ibcLiquidado, 9), 9) //IBC ARL
							fila += formatoDato(completarSecuencia(ibcLiquidado, 9), 9) //IBC CCF
						}

						fila += formatoDato(completarSecuencia(16, 7), 7) //Tarifa de aportes pensiones

						//Cotización obligatoria a pensiones
						for _, pago := range conceptosSeguridadSocial {
							switch pago.NombreConcepto {
							case "pensionTotal":
								pagoPension, _ = strconv.Atoi(obtenerPago(strconv.Itoa(idDescSegSocial), strconv.Itoa(detLiq[0].Id), strconv.Itoa(pago.Id)))
							case "saludTotal":
								pagoSalud, _ = strconv.Atoi(obtenerPago(strconv.Itoa(idDescSegSocial), strconv.Itoa(detLiq[0].Id), strconv.Itoa(pago.Id)))
							case "icbf":
								pagoIcbf, _ = strconv.Atoi(obtenerPago(strconv.Itoa(idDescSegSocial), strconv.Itoa(detLiq[0].Id), strconv.Itoa(pago.Id)))
							case "caja":
								pagoCaja, _ = strconv.Atoi(obtenerPago(strconv.Itoa(idDescSegSocial), strconv.Itoa(detLiq[0].Id), strconv.Itoa(pago.Id)))
							case "arl":
								pagoArl, _ = strconv.Atoi(obtenerPago(strconv.Itoa(idDescSegSocial), strconv.Itoa(detLiq[0].Id), strconv.Itoa(pago.Id)))
							}
						}

						fila += formatoDato(completarSecuencia(pagoPension, 9), 9) // Cotización obligatoria a pensiones

						fila += formatoDato(completarSecuencia(0, 9), 9) //Aporte voluntario del afiliado al fondo de pensiones obligatorias

						//Aporte voluntario del aportante al fondo de pensiones obligatoria
						errAporteVoluntario := getJson("http://"+beego.AppConfig.String("titanServicio")+
							"/detalle_preliquidacion"+
							"?query=NumeroContrato:"+detallePreliquidacion[i].NumeroContrato, &detLiq)

						if errAporteVoluntario != nil {
							fmt.Println("errAporteVoluntario", errAporteVoluntario)
							fila += formatoDato(completarSecuencia(0, 9), 9)
						} else {
							for _, liquidado := range detLiq {
								if liquidado.Concepto.NombreConcepto == "nombreRegla2176" {
									fila += formatoDato(strconv.FormatFloat(liquidado.ValorCalculado, 'E', -1, 64), 9)
								} else if liquidado.Concepto.NombreConcepto == "nombreRegla2178" {
									fila += formatoDato(strconv.FormatFloat(liquidado.ValorCalculado, 'E', -1, 64), 9)
								} else if liquidado.Concepto.NombreConcepto == "nombreRegla2173" {
									fila += formatoDato(strconv.FormatFloat(liquidado.ValorCalculado, 'E', -1, 64), 9)
								} else {
									break
								}
							}
						}

						fila += formatoDato(completarSecuencia(0, 9), 9)         // Total cotización Sistema General de Pensiones
						fila += formatoDato(completarSecuencia(0, 9), 9)         // Aportes a fondo de solidaridad pensional subcuenta de solidaridad
						fila += formatoDato(completarSecuencia(0, 9), 9)         // Aportes a fondo de solidaridad pensional subcuenta de subsistencia
						fila += formatoDato(completarSecuencia(0, 9), 9)         // Valor no retenido por aportes voluntarios
						fila += formatoDato("12.5", 7)                           // Tarifa de aportes salud
						fila += formatoDato(completarSecuencia(pagoSalud, 9), 9) // Cotización obligatoria a salud

						fila += formatoDato(completarSecuencia(0, 9), 9) //Valor UPC Adicional
						fila += formatoDato("", 15)                      //Nº de autorización de la incapacidad por enfermedad general
						fila += formatoDato(completarSecuencia(0, 9), 9) //Valor de la incapacidad por enfermedad general
						fila += formatoDato("", 15)                      //Nº de autorización de la licencia de maternidad o paternidad
						fila += formatoDato(completarSecuencia(0, 9), 9) //Valor de la licencia de maternidad

						fila += formatoDato(completarSecuenciaString("0.000522", 9), 9) //Tarifa de aportes a Riegos Laborales

						fila += formatoDato(completarSecuenciaString("0", 9), 9) //Centro de trabajo CT
						fila += formatoDato(completarSecuencia(pagoArl, 9), 9)   // Cotización obligatoria a salud

						fila += formatoDato(completarSecuenciaString("4", 7), 7) //Tarifa de aportes CCF
						fila += formatoDato(completarSecuencia(pagoCaja, 9), 9)  // Cotización obligatoria a salud

						fila += formatoDato(completarSecuencia(0, 7), 7) //Tarifa de aportes SENA
						fila += formatoDato(completarSecuencia(0, 9), 9) //Valor Aportes SENA

						fila += formatoDato(completarSecuencia(3, 7), 7)        //Tarifa de aportes ICBF
						fila += formatoDato(completarSecuencia(pagoIcbf, 9), 9) // Cotización obligatoria a salud

						fila += formatoDato(completarSecuencia(0, 7), 7) //Tarifa de aportes ESAP
						fila += formatoDato(completarSecuencia(0, 9), 9) //Valor de aporte ESAP
						fila += formatoDato(completarSecuencia(0, 7), 7) //Tarifa de aportes MEN
						fila += formatoDato(completarSecuencia(0, 9), 9) //Valor de aporte MEN

						//Para los registros de las UPC
						/*for _, upcAdicional := range upc {
							if upcAdicional.PersonaAsociada == detallePreliquidacion[i].Persona {
								fila += formatoDato(texto, longitud)
							}
						}*/

						// Estos campos están vacios porque solo aplican a los registros que son upc
						fila += formatoDato(" ", 2)  // Tipo de documento del cotizante principal
						fila += formatoDato(" ", 16) // Número de identificación del cotizante principal

						fila += formatoDato("N", 1)     // Cotizante exonerado de pago de aporte salud, SENA e ICBF - Ley 1607 de 2012
						fila += formatoDato("14-23", 6) // Código de la administradora de Riesgos Laborales a la cual pertenece el afiliado
						fila += formatoDato("1", 1)     // Clase de Riesgo en la que se encuentra el afiliado
						fila += formatoDato("", 1)      // Indicador tarifa especial pensiones (Actividades de alto riesgo, Senadores, CTI y Aviadores aplican)

						//Fechas de novedades (AAAA-MM-DD)
						fila += formatoDato(fechaIngreso, 10)          //Fecha ingreso
						fila += formatoDato(fechaRetiro, 10)           //Fecha retiro
						fila += formatoDato(fechaInicioVsp, 10)        //Fecha inicio VSP
						fila += formatoDato(fechaInicioSuspencion, 10) //Fecha inicio SLN
						fila += formatoDato(fechaFinSuspencion, 10)    //Fecha fin SLN
						fila += formatoDato(fechaInicioIge, 10)        //Fecha inicio IGE
						fila += formatoDato(fechaFinIge, 10)           //Fecha fin IGE
						fila += formatoDato(fechaInicioLma, 10)        //Fecha inicio LMA
						fila += formatoDato(fechaFinLma, 10)           //Fecha fin LMA
						fila += formatoDato(fechaInicioVac, 10)        //Fecha inicio VAC-LR
						fila += formatoDato(fechaFinVac, 10)           //Fecha fin VAC-LR
						fila += formatoDato(fechaInicioVct, 10)        //Fecha inicio VCT
						fila += formatoDato(fechaFinVct, 10)           //Fecha fin VCT
						fila += formatoDato(fechaInicioIrl, 10)        //Fecha inicio IRL
						fila += formatoDato(fechaFinIrl, 10)           //Fecha fin IRL

						fila += formatoDato(completarSecuencia(ibcLiquidado, 9), 9) //IBC otros parafiscales difenrentes a CCF
						fila += formatoDato("240", 3)
						fila += "\n"
						secuencia++
						fmt.Println("aqui va uno")
					}
				}
			}
		}
		c.Data["json"] = fila
	}
	c.ServeJSON()
}

func AgregarUpc(idPersonaAsociada string) {
	/*var upcAdicional []models.TipoUpc

	errUpcAdicional := getJson("http://"+beego.AppConfig.String("seguridadSocialService")+
		"/upc_adicional?limit=0,", target)

	errDiasLiquidados := getJson("http://"+beego.AppConfig.String("titanServicio")+
		"/detalle_liquidacion?limit=0"+
		"&query=Concepto.NombreConcepto:ibc_novedad,Persona:"+
		strconv.Itoa(detallePreliquidacion[i].Persona), &detLiq)*/
}

// Función para generalizar las novedades
func establecerNovedades(idPersona string) {

	//Variables para cada una de las novedades y sus días validos
	ingreso = false
	fechaIngreso = ""

	retiro = false
	fechaRetiro = ""

	trasladoPensiones = false
	trasladoEps = false
	exterior = false
	suspencionContrato = false
	fechaInicioSuspencion = ""
	fechaFinSuspencion = ""

	licenciaNoRem = false
	comisionServicios = false
	incapacidadGeneral = false
	fechaInicioIge = ""
	fechaFinIge = ""

	licenciaMaternidad = false
	fechaInicioLma = ""
	fechaFinLma = ""

	vacaciones = false
	licenciaRem = false
	fechaInicioVac = ""
	fechaFinVac = ""

	aporteVoluntario = false
	variacionCentroTrabajo = false
	fechaInicioVct = ""
	fechaFinVct = ""

	diasIncapcidadLaboral = 0
	fechaInicioIrl = ""
	fechaFinIrl = ""

	fechaInicioVsp = ""

	// Array para tener todas las novedades
	var conceptoPersona []models.ConceptoPorPersona
	var ConceptoNominaPersona []models.ConceptoNominaPorPersona

	var detallePreliquidacion []models.DetallePreliquidacion

	errConceptoPersona := getJson("http://"+beego.AppConfig.String("titanServicio")+
		"/concepto_nomina_por_persona"+
		"?limit=0"+
		"&query=activo:true,Persona:"+idPersona+
		",Concepto.NaturalezaConcepto.Nombre:seguridad_social", &ConceptoNominaPersona)

	if errConceptoPersona != nil {
		fmt.Println("errConceptoPersona: ", errConceptoPersona)
	}

	for h := 0; h < len(conceptoPersona); h++ {
		switch conceptoPersona[h].Concepto.NombreConcepto {
		case "retiro":
			retiro = true
			fechaRetiro = conceptoPersona[h].FechaDesde.Format(formatoFecha)
		case "ingreso":
			ingreso = true
			fechaIngreso = conceptoPersona[h].FechaDesde.Format(formatoFecha)
		case "exterior_familia":
			exterior = true
			//novedad = true
		case "suspencionContrato":
			suspencionContrato = true
			fechaInicioSuspencion = conceptoPersona[h].FechaDesde.Format(formatoFecha)
			fechaFinSuspencion = conceptoPersona[h].FechaDesde.Format(formatoFecha)
			//novedad = true
		case "licenciaNoRem":
			licenciaNoRem = true
			fechaInicioSuspencion = conceptoPersona[h].FechaDesde.Format(formatoFecha)
			fechaFinSuspencion = conceptoPersona[h].FechaDesde.Format(formatoFecha)
			//novedad = true
		case "comision_norem":
			comisionServicios = true
			fechaInicioSuspencion = conceptoPersona[h].FechaDesde.Format(formatoFecha)
			fechaFinSuspencion = conceptoPersona[h].FechaDesde.Format(formatoFecha)
			//novedad = true
		case "incapacidad_general":
			incapacidadGeneral = true
			fechaInicioIge = conceptoPersona[h].FechaDesde.Format(formatoFecha)
			fechaFinIge = conceptoPersona[h].FechaHasta.Format(formatoFecha)
			//novedad = true
		case "licenciaMaternidad":
			licenciaMaternidad = true
			fechaInicioLma = conceptoPersona[h].FechaDesde.Format(formatoFecha)
			fechaFinLma = conceptoPersona[h].FechaHasta.Format(formatoFecha)
			//novedad = true
		case "vacaciones":
			vacaciones = true
			fechaInicioVac = conceptoPersona[h].FechaDesde.Format(formatoFecha)
			fechaFinVac = conceptoPersona[h].FechaHasta.Format(formatoFecha)
			//novedad = true
		case "licencia_rem":
			licenciaRem = true
			fechaInicioVac = conceptoPersona[h].FechaDesde.Format(formatoFecha)
			fechaFinVac = conceptoPersona[h].FechaHasta.Format(formatoFecha)
			//novedad = true
		case "aporteVoluntario":
			aporteVoluntario = true
			//novedad = true
		case "variacionCentroTrabajo":
			variacionCentroTrabajo = true
			//novedad = true
		case "incapacidad_laboral":
			errEnfermedadLaboral := getJson("http://"+beego.AppConfig.String("titanServicio")+
				"/detalle_liquidacion?limit=-1", &detallePreliquidacion)
			if errEnfermedadLaboral != nil {
				fmt.Println("errEnfermedadLaboral: ", errEnfermedadLaboral)
			}
			diasIncapcidadLaboral = int(detallePreliquidacion[0].DiasLiquidados)
			fechaInicioIrl = conceptoPersona[h].FechaDesde.Format(formatoFecha)
			fechaFinIrl = conceptoPersona[h].FechaHasta.Format(formatoFecha)
			//valorIncapacidadLaboral = int(detalleIncapcidadLaboral[0].ValorCalculado)
			//novedad = true
		}
	}

	fila += formatoDato("", 1) //ING:Ingreso
	fila += formatoDato("", 1) //RET: retiro
	fila += formatoDato("", 1) //TDE: Traslado desde otra EPS o EOC
	fila += formatoDato("", 1) //TAE: Traslado a otra EPS o EOC
	//TDP: Traslado desde otra administradora de pensiones
	if trasladoPensiones {
		fila += formatoDato("X", 1)
	} else {
		fila += formatoDato("", 1)
	}

	fila += formatoDato("", 1) //TAP: Traslado a otra administradora de pensiones
	fila += formatoDato("", 1) //VSP: Variación permanente de salario
	fila += formatoDato("", 1) //Correcciones
	fila += formatoDato("", 1) //VST: Variación transitoria de salario

	//SLN: Suspención temporal del contrato de tabajo o licencia no remunerada o comisión de servicios
	if suspencionContrato {
		fila += formatoDato("X", 1)
	} else if licenciaNoRem {
		fila += formatoDato("X", 1)
	} else if comisionServicios {
		fila += formatoDato("C", 1)
	} else {
		fila += formatoDato("", 1)
	}

	//IGE: Incapacidad temporal por enfermedad general
	if incapacidadGeneral {
		fila += formatoDato("X", 1)
	} else {
		fila += formatoDato("", 1)
	}

	//LMA: Licencia de Maternidad o paternidad
	if licenciaMaternidad { //
		fila += formatoDato("X", 1)
	} else {
		fila += formatoDato("", 1)
	}

	//VAC: Vacaciones
	if vacaciones {
		fila += formatoDato("X", 1)
	} else if licenciaRem {
		fila += formatoDato("L", 1)
	} else {
		fila += formatoDato("", 1)
	}

	//AVP: Aporte voluntario
	if aporteVoluntario {
		fila += formatoDato("X", 1)
	} else {
		fila += formatoDato("", 1)
	}

	//VCT: Variación centros de trabajo
	if variacionCentroTrabajo {
		fila += formatoDato("X", 1)
	} else {
		fila += formatoDato("", 1)
	}

	//IRL: Días de incapacidad por accidente de trabajo o enfermedad laboral
	fila += formatoDato(completarSecuencia(diasIncapcidadLaboral, 2), 2)
}

func obtenerPago(idPeriodoPago, idDetalleLiqidacion, idTipoPago string) (valorPago string) {
	var pago []models.Pago
	/* Se obtiene un pago especefico con el periodo de pago, el detalle de la
	liquidacion y el tipo de pago */
	errPagosSalud := getJson("http://"+beego.AppConfig.String("seguridadSocialService")+
		"/pago?limit=1&query=PeriodoPago.Id:"+idPeriodoPago+",DetallePreliquidacion:"+idDetalleLiqidacion+
		",TipoPago:"+idTipoPago, &pago)
	if errPagosSalud != nil {
		fmt.Println("errPagosSalud: ", errPagosSalud)
	} else {
		valorPago = strconv.FormatFloat(pago[0].Valor, 'G', -1, 64)
	}
	return
}

func establecerNovedadesTranslado() {
}

/*

// Plantillas de pensionados
func (c *PlanillasController) GenerarPlanillaPensionados() {
	idStr := c.Ctx.Input.Param(":id")
	idDescSegSocial, _ := strconv.Atoi(idStr)
	var proveedores []models.InformacionProveedor
	var personasNatural []models.InformacionPersonaNatural
	var informacionPensionado []models.InformacionPensionado
	var conceptoPersona []models.ConceptoPorPersona
	var detallePreliquidacion []models.DetallePreliquidacion
	var personaNatural []models.InformacionPersonaNatural
	var pagosSalud []models.Pago
	var detallePreliquidacionConceptos []models.DetallePreliquidacion
	var errStrings []string
	tipoRegistro := "02"
	fila := ""

	errLiquidacion := getJson("http://"+beego.AppConfig.String("titanServicio")+"/detalle_preliquidacion"+
		"?limit=-1", &detallePreliquidacion)
	if errLiquidacion != nil {
		errStrings = append(errStrings, errLiquidacion.Error())
	}

	errProveedores := getJson("http://"+beego.AppConfig.String("titanServicio")+"/informacion_proveedor"+
		"?limit=-1", &proveedores)
	if errProveedores != nil {
		errStrings = append(errStrings, errProveedores.Error())
	}

	errPersonaNatural := getJson("http://"+beego.AppConfig.String("titanServicio")+"/informacion_persona_natural"+
		"?limit=0", &personasNatural)
	if errPersonaNatural != nil {
		errStrings = append(errStrings, errPersonaNatural.Error())
	}

	fmt.Println("**errStrings: ", errStrings)
	if errStrings == nil {
		secuencia := 1
		for i := 0; i < len(proveedores); i++ {
			for j := 0; j < len(detallePreliquidacion); j++ {
				for k := 0; k < len(personasNatural); k++ {
					if proveedores[i].Id == detallePreliquidacion[j].Persona {
						if int(proveedores[i].NumDocumento) == personasNatural[k].Id {
							var ibcLiquidado int = 0
							//Novedades
							var ingreso = false
							fechaIngreso := ""

							var retiro = false
							fechaRetiro := ""

							var trasladoPensiones = false
							var trasladoEps = false
							var exterior = false
							var suspencionContrato = false
							fechaInicioSuspencion := ""
							fechaFinSuspencion := ""

							var licenciaNoRem = false
							var comisionServicios = false
							var incapacidadGeneral = false
							fechaInicioIge := ""
							fechaFinIge := ""

							var licenciaMaternidad = false
							fechaInicioLma := ""
							fechaFinLma := ""

							var vacaciones = false
							var licenciaRem = false
							fechaInicioVac := ""
							fechaFinVac := ""

							var aporteVoluntario = false
							var variacionCentroTrabajo = false
							fechaInicioVct := ""
							fechaFinVct := ""

							var diasIncapacidadLaboral = 0
							fechaInicioIrl := ""
							fechaFinIrl := ""

							fechaInicioVsp := ""

							fila += formatoDato(tipoRegistro, 2)                              //Tipo registro
							fila += formatoDato(completarSecuencia(secuencia, 5), 5)          //Secuencia
							fila += formatoDato(personasNatural[k].PrimerApellido, 20)        //Primer apellido
							fila += formatoDato(personasNatural[k].SegundoApellido, 30)       //Segundo apellido
							fila += formatoDato(personasNatural[k].PrimerNombre, 20)          //Primer nombre
							fila += formatoDato(personasNatural[k].SegundoApellido, 30)       //Segundo nombre
							fila += formatoDato("CC", 2)                                      //Tipo de documento del pensionado
							fila += formatoDato(strconv.Itoa(int(personasNatural[k].Id)), 16) //Número de identificacion
							errPensionado := getJson("http://"+beego.AppConfig.String("titanServicio")+"/informacion_pensionado"+
								"?limit=1&query=InformacionProveedor:"+strconv.Itoa(proveedores[i].Id), &informacionPensionado)
							if errPensionado != nil {
								fmt.Println("errPensionado: ", errPensionado)
							} else {
								fila += formatoDato(completarSecuencia(informacionPensionado[0].TipoPension.Id, 2), 2) //Tipo de pensión
							}

							fila += formatoDato("N", 1) //Pensión compartida
							fila += formatoDato("", 20) //Primer apellido del causante
							fila += formatoDato("", 30) //Segundo apellido del causante
							fila += formatoDato("", 20) //Primer nombre del causante
							fila += formatoDato("", 30) //Segundo nombre del causante
							fila += formatoDato("", 2)  //Tipo de identificacion del causante de la pension
							fila += formatoDato("", 16) //Número de identificación del causante
							fila += formatoDato(" ", 1) //Tipo de pensionado

							errConceptoPersona := getJson("http://"+beego.AppConfig.String("titanServicio")+
								"/concepto_por_persona"+
								"?limit=0"+
								"&query=EstadoNovedad:Activo,Persona.Id:"+strconv.Itoa(detallePreliquidacion[i].Persona)+
								",Concepto.Naturaleza:seguridad_social", &conceptoPersona)

							if errConceptoPersona != nil {
								fmt.Println("errConceptoPersona: ", errConceptoPersona)
							}

							fmt.Println("Conceptos para el id: ", detallePreliquidacion[i].Persona)
							fmt.Println(conceptoPersona[0].Concepto.NombreConcepto)
							for h := 0; h < len(conceptoPersona); h++ {
								switch conceptoPersona[h].Concepto.NombreConcepto {
								case "retiro":
								case "ingreso":
								case "exterior_familia":
									exterior = true
									//novedad = true
								case "suspencionContrato":
									suspencionContrato = true
									//novedad = true
								case "licenciaNoRem":
									licenciaNoRem = true
									//novedad = true
								case "comision_norem":
									comisionServicios = true
									//novedad = true
								case "incapacidad_general":
									incapacidadGeneral = true
									//novedad = true
								case "licenciaMaternidad":
									licenciaMaternidad = true
									//novedad = true
								case "vacaciones":
									vacaciones = true
									//novedad = true
								case "licencia_rem":
									licenciaRem = true
									//novedad = true
								case "aporteVoluntario":
									aporteVoluntario = true
									//novedad = true
								case "variacionCentroTrabajo":
									variacionCentroTrabajo = true
									//novedad = true
								case "incapacidad_laboral":
									errEnfermedadLaboral := getJson("http://"+beego.AppConfig.String("titanServicio")+
										"/detalle_liquidacion?limit=-1", &detallePreliquidacionConceptos)
									if errEnfermedadLaboral != nil {
										fmt.Println("errEnfermedadLaboral: ", errEnfermedadLaboral)
									}
									diasIncapacidadLaboral, _ = strconv.Atoi(detallePreliquidacionConceptos[0].DiasLiquidados)
									//valorIncapacidadLaboral = int(detalleIncapcidadLaboral[0].ValorCalculado)
									//novedad = true
								}
							}

							if exterior {
								fila += formatoDato("X", 1) //Colombiano en el exterior
								fila += formatoDato(" ", 2) //Código del departamento de la ubicación laboral
								fila += formatoDato(" ", 3) //Código del municipio de ubicación laboral
							} else {
								fila += formatoDato("", 1)    //Colombiano en el exterior
								fila += formatoDato("11", 2)  //Código del departamento de la ubicación laboral
								fila += formatoDato("001", 3) //Código del municipio de ubicación laboral
							}

							errPersonaNatural := getJson("http://"+beego.AppConfig.String("titanServicio")+
								"/informacion_persona_natural"+
								"?limit=1"+
								"&query=Id:"+strconv.FormatFloat(proveedores[j].NumDocumento, 'E', -1, 64), &personaNatural)
							/*errPersonaNatural := getJson("http://"+beego.AppConfig.String("agoraServicio")+
							"/informacion_persona_natural"+
							"?limit=1"+
							"&query=Id:"+strconv.FormatFloat(proveedores[j].NumDocumento, 'E', -1, 64), &personaNatural)

							if errPersonaNatural != nil {
								fmt.Println("errPersonaNatural: ", errPersonaNatural)
							}

							fila += formatoDato("", 1) //ING:Ingreso
							fila += formatoDato("", 1) //RET: retiro
							fila += formatoDato("", 1) //TDE: Traslado desde otra EPS o EOC
							fila += formatoDato("", 1) //TAE: Traslado a otra EPS o EOC
							//TDP: Traslado desde otra administradora de pensiones
							if trasladoPensiones {
								fila += formatoDato("X", 1)
							} else {
								fila += formatoDato("", 1)
							}

							fila += formatoDato("", 1) //TAP: Traslado a otra administradora de pensiones
							fila += formatoDato("", 1) //VSP: Variación permanente de salario
							fila += formatoDato("", 1) //Correcciones
							fila += formatoDato("", 1) //VST: Variación transitoria de salario

							//SLN: Suspención temporal del contrato de tabajo o licencia no remunerada o comisión de servicios
							if suspencionContrato {
								fila += formatoDato("X", 1)
							} else if licenciaNoRem {
								fila += formatoDato("X", 1)
							} else if comisionServicios {
								fila += formatoDato("C", 1)
							} else {
								fila += formatoDato("", 1)
							}

							//IGE: Incapacidad temporal por enfermedad general
							if incapacidadGeneral {
								fila += formatoDato("X", 1)
							} else {
								fila += formatoDato("", 1)
							}

							//LMA: Licencia de Maternidad o paternidad
							if licenciaMaternidad { //
								fila += formatoDato("X", 1)
							} else {
								fila += formatoDato("", 1)
							}

							//VAC: Vacaciones
							if vacaciones {
								fila += formatoDato("X", 1)
							} else if licenciaRem {
								fila += formatoDato("L", 1)
							} else {
								fila += formatoDato("", 1)
							}

							//AVP: Aporte voluntario
							if aporteVoluntario {
								fila += formatoDato("X", 1)
							} else {
								fila += formatoDato("", 1)
							}

							//VCT: Variación centros de trabajo
							if variacionCentroTrabajo {
								fila += formatoDato("X", 1)
							} else {
								fila += formatoDato("", 1)
							}

							//IRL: Días de incapacidad por accidente de trabajo o enfermedad laboral
							fila += formatoDato(completarSecuencia(diasIncapacidadLaboral, 2), 2)

							//Código de la administradora de fondo de pensiones a la cual pertenece el afiliado
							fila += formatoDato("231001", 6)

							//Código de la admnistradora de pensiones a la cual se traslada el afiliado
							if trasladoPensiones {
								fila += formatoDato("230301", 6)
							} else {
								fila += formatoDato(" ", 6)
							}

							//Código EPS o EOC a la cual pertenece el afiliado
							fila += formatoDato("EPS010", 6)

							//Código EPS o EOC a la cual se traslada el afiliado
							if trasladoEps {
								fila += formatoDato("EPS012", 6)
							} else {
								fila += formatoDato(" ", 6)
							}

							//Código CCF a la cual pertenece el afiliado
							fila += formatoDato("CCF04", 6)

							errDiasLiquidados := getJson("http://"+beego.AppConfig.String("titanServicio")+
								"/detalle_liquidacion?limit=0"+
								"&query=Concepto.NombreConcepto:ibc_novedad,Persona:"+
								strconv.Itoa(detallePreliquidacion[i].Persona), &detallePreliquidacionConceptos)
							if errDiasLiquidados != nil {
								fmt.Println("errDiasLiquidados: ", errDiasLiquidados)
							}
							diasCotizados, _ := strconv.Atoi(detallePreliquidacionConceptos[0].DiasLiquidados)

							if ingreso || retiro {
								fila += formatoDato(completarSecuencia(diasCotizados, 2), 2) //Número de días cotizados a pensión
								fila += formatoDato(completarSecuencia(diasCotizados, 2), 2) //Número de días cotizados a salud
							} else {
								fila += formatoDato("30", 2) //Número de días cotizados a pensión
								fila += formatoDato("30", 2) //Número de días cotizados a salud
							}

							if incapacidadGeneral || licenciaMaternidad || vacaciones || diasIncapacidadLaboral > 0 {
								fila += formatoDato(completarSecuencia(diasCotizados, 2), 2) //Número de días cotizados a ARL
								fila += formatoDato(completarSecuencia(diasCotizados, 2), 2) //Número de días cotizados a CCF
							} else {
								fila += formatoDato("30", 2) //Número de días cotizados a ARL
								fila += formatoDato("30", 2) //Número de días cotizados a CCF
							}

							errSalarioBasico := getJson("http://"+beego.AppConfig.String("titanServicio")+
								"/detalle_liquidacion?limit=0"+
								"&query=Concepto.NombreConcepto:salarioBase,Persona:"+
								strconv.Itoa(detallePreliquidacion[i].Persona), &detallePreliquidacionConceptos)
							if errSalarioBasico != nil {
								fmt.Println("errSalarioBasico: ", errSalarioBasico)
							} else {
								salarioBase := strconv.FormatInt(detallePreliquidacionConceptos[0].ValorCalculado, 10)
								fila += formatoDato(salarioBase, 9) //Salario básico
							}

							fila += formatoDato("", 1) //Salario integral

							errSoloLiquidado := getJson("http://"+beego.AppConfig.String("titanServicio")+
								"/detalle_liquidacion?limit=0"+
								"&query=Concepto.NombreConcepto:ibc_liquidado,Persona:"+
								strconv.Itoa(detallePreliquidacion[i].Persona), &detallePreliquidacionConceptos)
							if errSoloLiquidado != nil {
								fmt.Println("errSoloLiquidado: ", errSoloLiquidado)
							} else {
								ibcLiquidado = int(detallePreliquidacionConceptos[0].ValorCalculado)
								fila += formatoDato(completarSecuencia(ibcLiquidado, 9), 9) //IBC pensión
								fila += formatoDato(completarSecuencia(ibcLiquidado, 9), 9) //IBC salud
								fila += formatoDato(completarSecuencia(ibcLiquidado, 9), 9) //IBC ARL
								fila += formatoDato(completarSecuencia(ibcLiquidado, 9), 9) //IBC CCF
							}

							fila += formatoDato(completarSecuencia(16, 7), 7) //Tarifa de aportes pensiones

							errPagosSalud := getJson("http://"+beego.AppConfig.String("seguridadSocialService")+
								"/pago"+
								"?query=PeriodoPago.Id:"+strconv.Itoa(idDescSegSocial)+
								",DetallePreliquidacion:"+strconv.Itoa(detallePreliquidacionConceptos[0].Id), &pagosSalud)
							if errPagosSalud != nil {
								fmt.Println("errPagosSalud: ", errPagosSalud)
							}

							fmt.Println("http://" + beego.AppConfig.String("seguridadSocialService") +
								"/pago" +
								"?query=PeriodoPago.Id:" + strconv.Itoa(idDescSegSocial) +
								",DetallePreliquidacion:" + strconv.Itoa(detallePreliquidacion[i].Id))

							//Cotización obligatoria a pensiones
							/*
								for _, pagoPension := range pagosSalud {
									if pagoPension.TipoPago.Nombre == "Pension" {
										fila += formatoDato(completarSecuencia(int(pagoPension.Valor), 9), 9)
										break
									}
								}

							fila += formatoDato(completarSecuencia(0, 9), 9) //Aporte voluntario del afiliado al fondo de pensiones obligatorias

							//Aporte voluntario del aportante al fondo de pensiones obligatoria
							errAporteVoluntario := getJson("http://"+beego.AppConfig.String("titanServicio")+
								"/detalle_liquidacion"+
								",Persona:"+strconv.Itoa(detallePreliquidacion[i].Persona), &detallePreliquidacionConceptos)

							if errAporteVoluntario != nil {
								fmt.Println("errAporteVoluntario", errAporteVoluntario)
								fila += formatoDato(completarSecuencia(0, 9), 9)
							} else {
								for _, liquidado := range detallePreliquidacionConceptos {
									if liquidado.Concepto.NombreConcepto == "nombreRegla2176" {
										fila += formatoDato(strconv.FormatInt(liquidado.ValorCalculado, 10), 9)
									} else if liquidado.Concepto.NombreConcepto == "nombreRegla2178" {
										fila += formatoDato(strconv.FormatInt(liquidado.ValorCalculado, 10), 9)
									} else if liquidado.Concepto.NombreConcepto == "nombreRegla2173" {
										fila += formatoDato(strconv.FormatInt(liquidado.ValorCalculado, 10), 9)
									} else {
										break
									}
								}
							}

							fila += formatoDato(completarSecuencia(0, 9), 9) //Total cotización Sistema General de Pensiones
							fila += formatoDato(completarSecuencia(0, 9), 9) //Aportes a fondo de solidaridad pensional subcuenta de solidaridad
							fila += formatoDato(completarSecuencia(0, 9), 9) //Aportes a fondo de solidaridad pensional subcuenta de subsistencia
							fila += formatoDato(completarSecuencia(0, 9), 9) //Valor no retenido por aportes voluntarios
							fila += formatoDato("12.5", 7)                   //Tarifa de aportes salud

							//Cotización obligatoria a salud
							/*
								for _, pagoSalud := range pagosSalud {
									if pagoSalud.TipoPago.Nombre == "Salud" {
										fila += formatoDato(completarSecuencia(int(pagoSalud.Valor), 9), 9)
										break
									}
								}

							fila += formatoDato(completarSecuencia(0, 9), 9) //Valor UPC Adicional
							fila += formatoDato("", 15)                      //Nº de autorización de la incapacidad por enfermedad general
							fila += formatoDato(completarSecuencia(0, 9), 9) //Valor de la incapacidad por enfermedad general
							fila += formatoDato("", 15)                      //Nº de autorización de la licencia de maternidad o paternidad
							fila += formatoDato(completarSecuencia(0, 9), 9) //Valor de la licencia de maternidad

							fila += formatoDato(completarSecuenciaString("0.000522", 9), 9) //Tarifa de aportes a Riegos Laborales

							fila += formatoDato(completarSecuenciaString("0", 9), 9) //Centro de trabajo CT

							//Cotización obligatoria a Sistema General de Riesgos Laborales
							/*
								for _, pagoArl := range pagosSalud {
									if pagoArl.TipoPago.Nombre == "ARL" {
										fila += formatoDato(completarSecuencia(int(pagoArl.Valor), 9), 9)
										break
									}
								}

							fila += formatoDato(completarSecuenciaString("4", 7), 7) //Tarifa de aportes CCF

							//Valor aporte CCF
							/*
								for _, pagoCaja := range pagosSalud {
									if pagoCaja.TipoPago.Nombre == "Caja" {
										fila += formatoDato(completarSecuencia(int(pagoCaja.Valor), 9), 9)
										break
									}
								}

							fila += formatoDato(completarSecuencia(0, 7), 7) //Tarifa de aportes SENA
							fila += formatoDato(completarSecuencia(0, 9), 9) //Valor Aportes SENA

							fila += formatoDato(completarSecuencia(3, 7), 7) //Tarifa de aportes ICBF

							//Valor aporte ICBF
							/*
								for _, pagoIcbf := range pagosSalud {
									if pagoIcbf.TipoPago.Nombre == "ICBF" {
										fila += formatoDato(completarSecuencia(int(pagoIcbf.Valor), 9), 9)
									}
								}

							fila += formatoDato(completarSecuencia(0, 7), 7) //Tarifa de aportes ESAP
							fila += formatoDato(completarSecuencia(0, 9), 9) //Valor de aporte ESAP
							fila += formatoDato(completarSecuencia(0, 7), 7) //Tarifa de aportes MEN
							fila += formatoDato(completarSecuencia(0, 9), 9) //Valor de aporte MEN

							//Para los registros de las UPC
							/*for _, upcAdicional := range upc {
								if upcAdicional.PersonaAsociada == detallePreliquidacion[i].Persona {
									fila += formatoDato(texto, longitud)
								}
							}

							// Estos campos están vacios porque solo aplican a los registros que osn upc
							fila += formatoDato(" ", 2)  //Tipo de documento del cotizante principal
							fila += formatoDato(" ", 16) //Número de identificación del cotizante principal

							fila += formatoDato("N", 1)     //Cotizante exonerado de pago de aporte salud, SENA e ICBF - Ley 1607 de 2012
							fila += formatoDato("14-23", 6) //Código de la administradora de Riesgos Laborales a la cual pertenece el afiliado
							fila += formatoDato("1", 1)     //Clase de Riesgo en la que se encuentra el afiliado
							fila += formatoDato("", 1)      //Indicador tarifa especial pensiones (Actividades de alto riesgo, Senadores, CTI y Aviadores aplican)

							//Fechas de novedades (AAAA-MM-DD)
							fila += formatoDato(fechaIngreso, 10)          //Fecha ingreso
							fila += formatoDato(fechaRetiro, 10)           //Fecha retiro
							fila += formatoDato(fechaInicioVsp, 10)        //Fecha inicio VSP
							fila += formatoDato(fechaInicioSuspencion, 10) //Fecha inicio SLN
							fila += formatoDato(fechaFinSuspencion, 10)    //Fecha fin SLN
							fila += formatoDato(fechaInicioIge, 10)        //Fecha inicio IGE
							fila += formatoDato(fechaFinIge, 10)           //Fecha fin IGE
							fila += formatoDato(fechaInicioLma, 10)        //Fecha inicio LMA
							fila += formatoDato(fechaFinLma, 10)           //Fecha fin LMA
							fila += formatoDato(fechaInicioVac, 10)        //Fecha inicio VAC-LR
							fila += formatoDato(fechaFinVac, 10)           //Fecha fin VAC-LR
							fila += formatoDato(fechaInicioVct, 10)        //Fecha inicio VCT
							fila += formatoDato(fechaFinVct, 10)           //Fecha fin VCT
							fila += formatoDato(fechaInicioIrl, 10)        //Fecha inicio IRL
							fila += formatoDato(fechaFinIrl, 10)           //Fecha fin IRL

							fila += formatoDato(completarSecuencia(ibcLiquidado, 9), 9) //IBC otros parafiscales difenrentes a CCF
							fila += formatoDato("240", 3)
							fila += "\n"
							secuencia++
						}
					}
				}
			}
			fmt.Println("Filas:\n", fila)
			c.Data["json"] = fila
		}
		c.ServeJSON()
	}
	c.ServeJSON()
}

func (c *PlanillasController) GenerarPlanillaN() {
	idStr := c.Ctx.Input.Param(":id")
	idDescSegSocial, _ := strconv.Atoi(idStr)
	var proveedores []models.InformacionProveedor
	var upc []models.UpcAdicional
	var detallePreliquidacion []models.DetallePreliquidacion
	var detLiq []models.DetallePreliquidacion
	var conceptos []models.Concepto
	var personaNatural []models.InformacionPersonaNatural
	var conceptosSeguridadSocial []models.Concepto
	var errStrings []string
	//formatoFecha := "2006-01-02"
	tipoRegistro := "02"

	// Se obtienen todos los conceptos de seguridad social en tabla conceptos de titan
	errConceptosSs := getJson("http://"+beego.AppConfig.String("titanServicio")+
		"/concepto?limit=0&query=Naturaleza:seguridad_social", &conceptosSeguridadSocial)
	if errConceptosSs != nil {
		fmt.Println("errConceptosSs: ", errConceptosSs)
	}

	errLiquidacion := getJson("http://"+beego.AppConfig.String("titanServicio")+
		"/detalle_liquidacion?limit=-1", &detallePreliquidacion)
	if errLiquidacion != nil {
		errStrings = append(errStrings, errLiquidacion.Error())
	}

	errProveedores := getJson("http://"+beego.AppConfig.String("titanServicio")+
		"/informacion_proveedor?limit=0", &proveedores)
	/*errProveedores := getJson("http://"+beego.AppConfig.String("agoraServicio")+
	"/informacion_proveedor?limit=0", &proveedores)
	if errProveedores != nil {
		errStrings = append(errStrings, errProveedores.Error())
	}

	errConceptos := getJson("http://"+beego.AppConfig.String("titanServicio")+
		"/concepto?limit=0", &conceptos)
	if errConceptos != nil {
		errStrings = append(errStrings, errConceptos.Error())
	}

	errUpc := getJson("http://"+beego.AppConfig.String("seguridadSocialService")+
		"/upc_adicional?limit=0", &upc)
	if errUpc != nil {
		errStrings = append(errStrings, errUpc.Error())
	}

	fmt.Println("errStrings: ", errStrings)
	if errStrings == nil {
		secuencia := 1
		for i := 0; i < len(detallePreliquidacion); i++ {
			for j := 0; j < len(proveedores); j++ {
				if detallePreliquidacion[i].Persona == proveedores[j].Id {
					if strings.Contains(fila, strconv.Itoa(int(proveedores[j].NumDocumento))) {
						break
					} else {
						var ibcLiquidado = 0
						var pagoSalud = 0
						var pagoPension = 0
						var pagoArl = 0
						var pagoCaja = 0
						var pagoIcbf = 0

						fila += formatoDato(tipoRegistro, 2)                     //Tipo Registro
						fila += formatoDato(completarSecuencia(secuencia, 5), 5) //Secuencia

						fila += formatoDato("CC", 2)                                            //Tip de documento del cotizante
						fila += formatoDato(strconv.Itoa(int(proveedores[j].NumDocumento)), 16) //Número de identificación del cotizante
						fila += formatoDato(completarSecuencia(1, 2), 2)                        //Tipo Cotizante
						fila += formatoDato(completarSecuencia(1, 2), 2)                        //Subtipo de Cotizante
						fila += formatoDato("", 1)                                              //Extranjero no obligado a cotizar pensión

						errPersonaNatural := getJson("http://"+beego.AppConfig.String("titanServicio")+
							"/informacion_persona_natural"+
							"?limit=1"+
							"&query=Id:"+strconv.FormatFloat(proveedores[j].NumDocumento, 'E', -1, 64), &personaNatural)

						/*errPersonaNatural := getJson("http://"+beego.AppConfig.String("agoraServicio")+
						"/informacion_persona_natural"+
						"?limit=1"+
						"&query=Id:"+strconv.FormatFloat(proveedores[j].NumDocumento, 'E', -1, 64), &personaNatural)

						if errPersonaNatural != nil {
							fmt.Println("errPersonaNatural: ", errPersonaNatural)
						}

						if exterior {
							fila += formatoDato("X", 1) //Colombiano en el exterior
							fila += formatoDato(" ", 2) //Código del departamento de la ubicación laboral
							fila += formatoDato(" ", 3) //Código del municipio de ubicación laboral
						} else {
							fila += formatoDato("", 1)    //Colombiano en el exterior
							fila += formatoDato("11", 2)  //Código del departamento de la ubicación laboral
							fila += formatoDato("001", 3) //Código del municipio de ubicación laboral
						}

						fila += formatoDato(personaNatural[0].PrimerApellido, 20)  //Primer apellido
						fila += formatoDato(personaNatural[0].SegundoApellido, 30) //Segundo apellido
						fila += formatoDato(personaNatural[0].PrimerNombre, 20)    //Primer nombre
						fila += formatoDato(personaNatural[0].SegundoNombre, 30)   //Segundo nombre

						// --AQUÍ VA LA FUNCIÓN DE LAS NOVEDADES!--  //
						establecerNovedades(strconv.Itoa(detallePreliquidacion[i].Persona))

						//Código de la administradora de fondo de pensiones a la cual pertenece el afiliado
						fila += formatoDato("231001", 6)

						//Código de la admnistradora de pensiones a la cual se traslada el afiliado
						// Si hay un translado, debe aparecer el nuevo código, de lo contrario será un campo vació
						if trasladoPensiones {
							fila += formatoDato("230301", 6)
						} else {
							fila += formatoDato(" ", 6)
						}

						//Código EPS o EOC a la cual pertenece el afiliado
						fila += formatoDato("EPS010", 6)

						//Código EPS o EOC a la cual se traslada el afiliado
						// Si hay un translado, debe aparecer el nuevo código, de lo contrario será un campo vació
						if trasladoEps {
							fila += formatoDato("EPS012", 6)
						} else {
							fila += formatoDato(" ", 6)
						}

						//Código CCF a la cual pertenece el afiliado
						fila += formatoDato("CCF04", 6)

						errDiasLiquidados := getJson("http://"+beego.AppConfig.String("titanServicio")+
							"/detalle_liquidacion?limit=0"+
							"&query=Concepto.NombreConcepto:ibc_novedad,Persona:"+
							strconv.Itoa(detallePreliquidacion[i].Persona), &detLiq)
						if errDiasLiquidados != nil {
							fmt.Println("errDiasLiquidados: ", errDiasLiquidados)
						}
						diasCotizados, _ := strconv.Atoi(detLiq[0].DiasLiquidados)

						if ingreso || retiro {
							fila += formatoDato(completarSecuencia(diasCotizados, 2), 2) //Número de días cotizados a pensión
							fila += formatoDato(completarSecuencia(diasCotizados, 2), 2) //Número de días cotizados a salud
						} else {
							fila += formatoDato("30", 2) //Número de días cotizados a pensión
							fila += formatoDato("30", 2) //Número de días cotizados a salud
						}

						if incapacidadGeneral || licenciaMaternidad || vacaciones || diasIncapcidadLaboral > 0 {
							fila += formatoDato(completarSecuencia(diasCotizados, 2), 2) //Número de días cotizados a ARL
							fila += formatoDato(completarSecuencia(diasCotizados, 2), 2) //Número de días cotizados a CCF
						} else {
							fila += formatoDato("30", 2) //Número de días cotizados a ARL
							fila += formatoDato("30", 2) //Número de días cotizados a CCF
						}

						fmt.Println(detallePreliquidacion[i].Persona)
						errSalarioBasico := getJson("http://"+beego.AppConfig.String("titanServicio")+
							"/detalle_liquidacion?limit=0"+
							"&query=Concepto.NombreConcepto:salarioBase,Persona:"+
							strconv.Itoa(detallePreliquidacion[i].Persona), &detLiq)
						if errSalarioBasico != nil {
							fmt.Println("errSalarioBasico: ", errSalarioBasico)
						} else {
							salarioBase := strconv.FormatInt(detLiq[0].ValorCalculado, 10)
							fila += formatoDato(salarioBase, 9) //Salario básico
						}

						fila += formatoDato("", 1) //Salario integral

						errSoloLiquidado := getJson("http://"+beego.AppConfig.String("titanServicio")+
							"/detalle_liquidacion?limit=0"+
							"&query=Concepto.NombreConcepto:ibc_liquidado,Persona:"+
							strconv.Itoa(detallePreliquidacion[i].Persona), &detLiq)
						if errSoloLiquidado != nil {
							fmt.Println("errSoloLiquidado: ", errSoloLiquidado)
						} else {
							ibcLiquidado = int(detLiq[0].ValorCalculado)
							fila += formatoDato(completarSecuencia(ibcLiquidado, 9), 9) //IBC pensión
							fila += formatoDato(completarSecuencia(ibcLiquidado, 9), 9) //IBC salud
							fila += formatoDato(completarSecuencia(ibcLiquidado, 9), 9) //IBC ARL
							fila += formatoDato(completarSecuencia(ibcLiquidado, 9), 9) //IBC CCF
						}

						fila += formatoDato(completarSecuencia(16, 7), 7) //Tarifa de aportes pensiones

						//Cotización obligatoria a pensiones
						for _, pago := range conceptosSeguridadSocial {
							switch pago.NombreConcepto {
							case "pensionTotal":
								pagoPension, _ = strconv.Atoi(obtenerPago(strconv.Itoa(idDescSegSocial), strconv.Itoa(detLiq[0].Id), strconv.Itoa(pago.Id)))
							case "saludTotal":
								pagoSalud, _ = strconv.Atoi(obtenerPago(strconv.Itoa(idDescSegSocial), strconv.Itoa(detLiq[0].Id), strconv.Itoa(pago.Id)))
							case "icbf":
								pagoIcbf, _ = strconv.Atoi(obtenerPago(strconv.Itoa(idDescSegSocial), strconv.Itoa(detLiq[0].Id), strconv.Itoa(pago.Id)))
							case "caja":
								pagoCaja, _ = strconv.Atoi(obtenerPago(strconv.Itoa(idDescSegSocial), strconv.Itoa(detLiq[0].Id), strconv.Itoa(pago.Id)))
							case "arl":
								pagoArl, _ = strconv.Atoi(obtenerPago(strconv.Itoa(idDescSegSocial), strconv.Itoa(detLiq[0].Id), strconv.Itoa(pago.Id)))
							}
						}

						fila += formatoDato(completarSecuencia(pagoPension, 9), 9) // Cotización obligatoria a pensiones

						fila += formatoDato(completarSecuencia(0, 9), 9) //Aporte voluntario del afiliado al fondo de pensiones obligatorias

						//Aporte voluntario del aportante al fondo de pensiones obligatoria
						errAporteVoluntario := getJson("http://"+beego.AppConfig.String("titanServicio")+
							"/detalle_liquidacion"+
							"?query=Persona:"+strconv.Itoa(detallePreliquidacion[i].Persona), &detLiq)

						if errAporteVoluntario != nil {
							fmt.Println("errAporteVoluntario", errAporteVoluntario)
							fila += formatoDato(completarSecuencia(0, 9), 9)
						} else {
							for _, liquidado := range detLiq {
								if liquidado.Concepto.NombreConcepto == "nombreRegla2176" {
									fila += formatoDato(strconv.FormatInt(liquidado.ValorCalculado, 10), 9)
								} else if liquidado.Concepto.NombreConcepto == "nombreRegla2178" {
									fila += formatoDato(strconv.FormatInt(liquidado.ValorCalculado, 10), 9)
								} else if liquidado.Concepto.NombreConcepto == "nombreRegla2173" {
									fila += formatoDato(strconv.FormatInt(liquidado.ValorCalculado, 10), 9)
								} else {
									break
								}
							}
						}

						fila += formatoDato(completarSecuencia(0, 9), 9)         // Total cotización Sistema General de Pensiones
						fila += formatoDato(completarSecuencia(0, 9), 9)         // Aportes a fondo de solidaridad pensional subcuenta de solidaridad
						fila += formatoDato(completarSecuencia(0, 9), 9)         // Aportes a fondo de solidaridad pensional subcuenta de subsistencia
						fila += formatoDato(completarSecuencia(0, 9), 9)         // Valor no retenido por aportes voluntarios
						fila += formatoDato("12.5", 7)                           // Tarifa de aportes salud
						fila += formatoDato(completarSecuencia(pagoSalud, 9), 9) // Cotización obligatoria a salud

						fila += formatoDato(completarSecuencia(0, 9), 9) //Valor UPC Adicional
						fila += formatoDato("", 15)                      //Nº de autorización de la incapacidad por enfermedad general
						fila += formatoDato(completarSecuencia(0, 9), 9) //Valor de la incapacidad por enfermedad general
						fila += formatoDato("", 15)                      //Nº de autorización de la licencia de maternidad o paternidad
						fila += formatoDato(completarSecuencia(0, 9), 9) //Valor de la licencia de maternidad

						fila += formatoDato(completarSecuenciaString("0.000522", 9), 9) //Tarifa de aportes a Riegos Laborales

						fila += formatoDato(completarSecuenciaString("0", 9), 9) //Centro de trabajo CT
						fila += formatoDato(completarSecuencia(pagoArl, 9), 9)   // Cotización obligatoria a salud

						fila += formatoDato(completarSecuenciaString("4", 7), 7) //Tarifa de aportes CCF
						fila += formatoDato(completarSecuencia(pagoCaja, 9), 9)  // Cotización obligatoria a salud

						fila += formatoDato(completarSecuencia(0, 7), 7) //Tarifa de aportes SENA
						fila += formatoDato(completarSecuencia(0, 9), 9) //Valor Aportes SENA

						fila += formatoDato(completarSecuencia(3, 7), 7)        //Tarifa de aportes ICBF
						fila += formatoDato(completarSecuencia(pagoIcbf, 9), 9) // Cotización obligatoria a salud

						fila += formatoDato(completarSecuencia(0, 7), 7) //Tarifa de aportes ESAP
						fila += formatoDato(completarSecuencia(0, 9), 9) //Valor de aporte ESAP
						fila += formatoDato(completarSecuencia(0, 7), 7) //Tarifa de aportes MEN
						fila += formatoDato(completarSecuencia(0, 9), 9) //Valor de aporte MEN

						//Para los registros de las UPC
						/*for _, upcAdicional := range upc {
							if upcAdicional.PersonaAsociada == detallePreliquidacion[i].Persona {
								fila += formatoDato(texto, longitud)
							}
						}

						// Estos campos están vacios porque solo aplican a los registros que son upc
						fila += formatoDato(" ", 2)  // Tipo de documento del cotizante principal
						fila += formatoDato(" ", 16) // Número de identificación del cotizante principal

						fila += formatoDato("N", 1)     // Cotizante exonerado de pago de aporte salud, SENA e ICBF - Ley 1607 de 2012
						fila += formatoDato("14-23", 6) // Código de la administradora de Riesgos Laborales a la cual pertenece el afiliado
						fila += formatoDato("1", 1)     // Clase de Riesgo en la que se encuentra el afiliado
						fila += formatoDato("", 1)      // Indicador tarifa especial pensiones (Actividades de alto riesgo, Senadores, CTI y Aviadores aplican)

						//Fechas de novedades (AAAA-MM-DD)
						fila += formatoDato(fechaIngreso, 10)          //Fecha ingreso
						fila += formatoDato(fechaRetiro, 10)           //Fecha retiro
						fila += formatoDato(fechaInicioVsp, 10)        //Fecha inicio VSP
						fila += formatoDato(fechaInicioSuspencion, 10) //Fecha inicio SLN
						fila += formatoDato(fechaFinSuspencion, 10)    //Fecha fin SLN
						fila += formatoDato(fechaInicioIge, 10)        //Fecha inicio IGE
						fila += formatoDato(fechaFinIge, 10)           //Fecha fin IGE
						fila += formatoDato(fechaInicioLma, 10)        //Fecha inicio LMA
						fila += formatoDato(fechaFinLma, 10)           //Fecha fin LMA
						fila += formatoDato(fechaInicioVac, 10)        //Fecha inicio VAC-LR
						fila += formatoDato(fechaFinVac, 10)           //Fecha fin VAC-LR
						fila += formatoDato(fechaInicioVct, 10)        //Fecha inicio VCT
						fila += formatoDato(fechaFinVct, 10)           //Fecha fin VCT
						fila += formatoDato(fechaInicioIrl, 10)        //Fecha inicio IRL
						fila += formatoDato(fechaFinIrl, 10)           //Fecha fin IRL

						fila += formatoDato(completarSecuencia(ibcLiquidado, 9), 9) //IBC otros parafiscales difenrentes a CCF
						fila += formatoDato("240", 3)
						fila += "\n"
						secuencia++
						fmt.Println("aqui va uno")
					}
				}
			}
		}
		c.Data["json"] = fila
	}
	c.ServeJSON()
}
*/

func completarSecuencia(num, cantSecuencia int) (secuencia string) {
	tamanioNum := len(strconv.Itoa(num))
	for i := 0; i < cantSecuencia-tamanioNum; i++ {
		secuencia += "0"
	}
	secuencia += strconv.Itoa(num)
	return
}

func completarSecuenciaString(num string, cantSecuencia int) (secuencia string) {
	tamanioNum := len(num)
	for i := 0; i < cantSecuencia-tamanioNum; i++ {
		secuencia += "0"
	}
	secuencia += num
	return
}

func formatoDato(texto string, longitud int) (textoEscribir string) {
	for _, r := range texto {
		textoEscribir += string(r)
	}
	for i := 0; i < longitud-len(texto); i++ {
		textoEscribir += " "
	}
	return
}

// Post ...
// @Title Create
// @Description create Planillas
// @Param	body		body 	models.Planillas	true		"body for Planillas content"
// @Success 201 {object} models.Planillas
// @Failure 403 body is empty
// @router / [post]
func (c *PlanillasController) Post() {

}

// GetOne ...
// @Title GetOne
// @Description get Planillas by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Planillas
// @Failure 403 :id is empty
// @router /:id [get]
func (c *PlanillasController) GetOne() {

}

// GetAll ...
// @Title GetAll
// @Description get Planillas
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Planillas
// @Failure 403
// @router / [get]
func (c *PlanillasController) GetAll() {

}

// Put ...
// @Title Put
// @Description update the Planillas
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Planillas	true		"body for Planillas content"
// @Success 200 {object} models.Planillas
// @Failure 403 :id is not int
// @router /:id [put]
func (c *PlanillasController) Put() {

}

// Delete ...
// @Title Delete
// @Description delete the Planillas
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *PlanillasController) Delete() {

}