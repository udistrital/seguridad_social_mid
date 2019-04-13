//Comienza la corrección...

package controllers

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/udistrital/ss_mid_api/models"
)

// PlanillasController operations for Planillas
type PlanillasController struct {
	beego.Controller
}

// URLMapping ...
func (c *PlanillasController) URLMapping() {
	c.Mapping("GenerarPlanillaActivos", c.GenerarPlanillaActivos)
}

var (
	contratistas = false

	formatoFecha = "2006-01-02"
	fila         string
	filaAux      string
	filas        string

	//Variables para cada una de las novedades y sus días validos
	ingreso      = false
	fechaIngreso = ""

	retiro      = false
	fechaRetiro = ""

	trasladoDesdeEps            = false
	trasladoDesdePensiones      = false
	trasladoAPensiones          = false
	trasladoAEps                = false
	variacionPermanteSalario    = false
	corecciones                 = false
	variacionTransitoriaSalario = false
	suspencionTemporalContrato  = false

	exterior = false

	fechaInicioSuspencion = ""
	fechaFinSuspencion    = ""

	licenciaNoRemunerada = false
	comisionServicios    = false
	incapacidadGeneral   = false
	fechaInicioIge       = ""
	fechaFinIge          = ""

	licenciaMaternidad = false
	fechaInicioLma     = ""
	fechaFinLma        = ""

	vacaciones         = false
	licenciaRemunerada = false
	fechaInicioVac     = ""
	fechaFinVac        = ""

	aporteVoluntario       = false
	variacionCentroTrabajo = false
	fechaInicioVct         = ""
	fechaFinVct            = ""

	diasIncapcidadLaboral = 0
	fechaInicioIrl        = ""
	fechaFinIrl           = ""

	fechaInicioVsp = ""

	diasArl       = 30
	tarifaIcbf    = "0.03000"
	tarifaCaja    = "0.04000"
	tarifaArl     = "0.0052200"
	tarifaPension = "0.16000"
	tarifaSalud   = "0.12500"

	tipoPreliquidacion string
	codigoFondoPension string
	codigoEps          string
	codigoCaja         string

	diasSuspencionContrato int
	diasLicenciaNoRem      int
	diasComisionServicios  int
	diasLicenciaRem        int
	diasLicenciaMaternidad int
	diasVacaciones         int
	diasIncapacidad        int

	valorUpc       string
	salarioBase    int
	mesPeriodo     int
	anioPeriodo    int
	novedadPersona = false
	diasIrl        = 0

	tipoRegistro = "02"
	secuencia    = 1
)

// GenerarPlanillaActivos ...
// @Title Generar planilla de activos
// @Description Recibe un periodo pago y devuelve un arreglo de json con la información de la planilla
// @Param	body body PeriodoPago true	"body for PeriodoPago"
// @Success 200 {string} string
// @Failure 403 body is empty
// @router /GenerarPlanillaActivos [post]
func (c *PlanillasController) GenerarPlanillaActivos() {
	var (
		periodoPago           *models.PeriodoPago
		detallePreliquidacion []models.DetallePreliquidacion
		personas              []string
	)

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &periodoPago); err == nil {

		tipoPreliquidacion = periodoPago.TipoLiquidacion
		mesPeriodo = int(periodoPago.Mes)
		anioPeriodo = int(periodoPago.Anio)

		if err = getJson("http://"+beego.AppConfig.String("titanServicio")+"/detalle_preliquidacion?"+
			"limit=-1"+
			"&query=Preliquidacion.Id:"+strconv.Itoa(periodoPago.Liquidacion)+
			",Concepto.NombreConcepto:ibc_liquidado", &detallePreliquidacion); err == nil {

			if periodoPago.TipoLiquidacion == "CT" {
				contratistas = true
			}
			filas = ""

			for i := range detallePreliquidacion {
				personas = append(personas, fmt.Sprint(detallePreliquidacion[i].Persona))
			}

			mapProveedores, _ := GetInfoProveedor(personas)
			mapPersonas, _ := GetInfoPersona(mapProveedores)
			for key, value := range mapPersonas {
				var (
					preliquidacion []models.DetallePreliquidacion
					ibcLiquidado,
					pagoSalud,
					pagoPension,
					pagoArl,
					pagoCaja,
					pagoIcbf,
					horasLaboradas string
					totalPagoPension int
				)
				idPersona := key
				idPreliquidacion := strconv.Itoa(detallePreliquidacion[0].Preliquidacion.Id)
				cedulaPersona := fmt.Sprint(value.Id)
				fila = ""
				filaAux = ""
				fila += formatoDato(tipoRegistro, 2)                     //Tipo Registro
				fila += formatoDato(completarSecuencia(secuencia, 5), 5) //Secuencia
				fila += formatoDato("CC", 2)                             //Tip de documento del cotizante
				fila += formatoDato(value.Id, 16)                        //Número de identificación del cotizante
				fila += formatoDato(completarSecuencia(1, 2), 2)         //Tipo Cotizante
				fila += formatoDato(completarSecuencia(0, 2), 2)         //Subtipo de Cotizante
				fila += formatoDato("", 1)                               //Extranjero no obligado a cotizar pensión

				establecerNovedadesExterior(idPersona, idPreliquidacion)

				fila += formatoDato(value.PrimerApellido, 20)  //Primer apellido
				fila += formatoDato(value.SegundoApellido, 30) //Segundo apellido
				fila += formatoDato(value.PrimerNombre, 20)    //Primer nombre
				fila += formatoDato(value.SegundoNombre, 30)   //Segundo nombre

				establecerNovedades(idPersona, idPreliquidacion, cedulaPersona)

				//Código de la administradora de fondo de pensiones a la cual pertenece el afiliado
				codigoFondoPension = traerCodigoEntidadSalud(strconv.Itoa(value.IdFondoPension))
				fila += codigoFondoPension

				//Código de la admnistradora de pensiones a la cual se traslada el afiliado
				// Si hay un translado, debe aparecer el nuevo código, de lo contrario será un campo vació
				if trasladoAPensiones {
					fila += formatoDato("230301", 6)
				} else {
					fila += formatoDato("", 6)
				}

				//Código EPS o EOC a la cual pertenece el afiliado
				if codigoEps = traerCodigoEntidadSalud(strconv.Itoa(value.IdEps)); codigoEps == "      " {
					fila += formatoDato("MIN001", 6)
				} else {
					fila += codigoEps
				}

				//Código EPS o EOC a la cual se traslada el afiliado
				// Si hay un translado, debe aparecer el nuevo código, de lo contrario será un campo vació
				if trasladoAEps {
					fila += formatoDato("EPS012", 6)
				} else {
					fila += formatoDato("", 6)
				}

				// fila += traerCodigoEntidadSalud(strconv.Itoa(value.IdCajaCompensacion)) //Código CCF a la cual pertenece el afiliado
				fila += formatoDato("CCF24", 6) //Código CCF a la cual pertenece el afiliado
				diasLaborados, _ := strconv.Atoi(traerDiasCotizados(idPersona, idPreliquidacion, "salud"))
				horasLaboradas = strconv.Itoa(diasLaborados * 8)

				fila += traerDiasCotizados(idPersona, idPreliquidacion, "pension") // Número de días coitzados a pensión
				fila += traerDiasCotizados(idPersona, idPreliquidacion, "salud")   // Número de días cotizados a salud
				fila += formatoDato(completarSecuencia(diasArl, 2), 2)             // Número de días cotizados a arl
				diasCaja, pagoCaja := traerDiasCotizadosEmpleador(idPersona, idPreliquidacion, strconv.Itoa(periodoPago.Id), "caja_compensacion")
				fila += diasCaja // Número de días cotizados a caja de compensación

				if contratistas {
					err = getJson("http://"+beego.AppConfig.String("titanServicio")+
						"/detalle_preliquidacion?limit=1"+
						"&fields=ValorCalculado"+
						"&query=Preliquidacion:"+strconv.Itoa(periodoPago.Liquidacion)+",Concepto.NombreConcepto:honorarios,Persona:"+key, &preliquidacion)
				} else {
					err = getJson("http://"+beego.AppConfig.String("titanServicio")+
						"/detalle_preliquidacion?limit=1"+
						"&fields=ValorCalculado"+
						"&query=Preliquidacion:"+strconv.Itoa(periodoPago.Liquidacion)+",Concepto.NombreConcepto:salarioBase,Persona:"+key, &preliquidacion)
				}

				if err == nil {
					salarioBase = int(preliquidacion[0].ValorCalculado)
					fila += formatoDato(completarSecuencia(salarioBase, 9), 9) //Salario básico
				}

				fila += formatoDato("", 1) //Salario integral

				err = getJson("http://"+beego.AppConfig.String("titanServicio")+
					"/detalle_preliquidacion?limit=1"+
					"&fields=ValorCalculado,Id"+
					"&query=Preliquidacion:"+strconv.Itoa(periodoPago.Liquidacion)+
					",Concepto.NombreConcepto:ibc_liquidado,Persona:"+idPersona, &preliquidacion)
				if err == nil {
					ibcLiquidadoTemp := int(preliquidacion[0].ValorCalculado)
					fila += formatoDato(completarSecuencia(ibcLiquidadoTemp, 9), 9) //IBC pensión
					fila += formatoDato(completarSecuencia(ibcLiquidadoTemp, 9), 9) //IBC salud
					fila += formatoDato(completarSecuencia(ibcLiquidadoTemp, 9), 9) //IBC ARL
					fila += formatoDato(completarSecuencia(ibcLiquidadoTemp, 9), 9) //IBC CCF
					ibcLiquidado = fmt.Sprint(ibcLiquidadoTemp)
				}

				fila += formatoDato("0.16000", 7) //Tarifa de aportes pensiones

				// Aquí se traen todos los valores totales a pagar correspondientes a seguridad social
				// (revisar comentarios de la función traerDiasCotizadosEmpleador(idPersona, idPreliquidacion, idPeriodoPago, tipo_pago) )
				_, pagoPension = traerDiasCotizadosEmpleador(idPersona, idPreliquidacion, strconv.Itoa(periodoPago.Id), "pension_ud")
				_, pagoSalud = traerDiasCotizadosEmpleador(idPersona, idPreliquidacion, strconv.Itoa(periodoPago.Id), "salud_ud")
				_, pagoIcbf = traerDiasCotizadosEmpleador(idPersona, idPreliquidacion, strconv.Itoa(periodoPago.Id), "icbf")
				_, pagoArl = traerDiasCotizadosEmpleador(idPersona, idPreliquidacion, strconv.Itoa(periodoPago.Id), "arl")

				fila += formatoDato(completarSecuenciaString(pagoPension, 9), 9) // Cotización obligatoria a pensiones

				// Aporte voluntario del afiliado al fondo de pensiones
				err = getJson("http://"+beego.AppConfig.String("titanServicio")+
					"/detalle_preliquidacion"+
					"?fields=Concepto,Id"+
					"&query=Preliquidacion:"+strconv.Itoa(periodoPago.Liquidacion)+",Persona:"+key, &preliquidacion)

				if err != nil {
					fila += formatoDato(completarSecuencia(0, 9), 9)
				} else {
					valorAporteVoluntarioPension := 0.0
					filaAporteVoluntarioPension := formatoDato(completarSecuencia(0, 9), 9)
					for i := 0; i < len(preliquidacion); i++ {
						tempMap := preliquidacion[i]
						valorCalculado := strconv.FormatFloat(tempMap.ValorCalculado, 'E', -1, 64)
						switch tempMap.Concepto.NombreConcepto {
						// se busca asignar el valor calculado de alguno de estos conceptos,
						// En caso de que alguno exista en el detalle_preliquidación, se cambia el valor de la fila
						// y se cambia el valor del aporte voluntario, para luego asignar la fila y sumar el aporte voluntario con el total
						// de cotización a pensión
						case "nombreRegla2176":
							filaAporteVoluntarioPension = formatoDato(valorCalculado, 9)
							valorAporteVoluntarioPension = tempMap.ValorCalculado
						case "nombreRegla2178":
							filaAporteVoluntarioPension = formatoDato(valorCalculado, 9)
							valorAporteVoluntarioPension = tempMap.ValorCalculado
						case "nombreRegla2173":
							filaAporteVoluntarioPension = formatoDato(valorCalculado, 9)
							valorAporteVoluntarioPension = tempMap.ValorCalculado
						}
					}
					fila += filaAporteVoluntarioPension
					auxPagoPension, _ := strconv.Atoi(pagoPension)
					totalPagoPension = int(valorAporteVoluntarioPension) + auxPagoPension
					fila += formatoDato(completarSecuencia(0, 9), 9)                                    // Aporte voluntario del aportante al fondo de pensiones obligatoria
					fila += formatoDato(completarSecuencia(totalPagoPension, 9), 9)                     // Total cotización Sistema General de Pensiones
					fila += traerValorConceptoEmpleado(idPersona, idPreliquidacion, "fondoSolidaridad") // Aportes a fondo de solidaridad pensional subcuenta de solidaridad
					fila += traerValorConceptoEmpleado(idPersona, idPreliquidacion, "fondoSolidaridad") // Aportes a fondo de solidaridad pensional subcuenta de subsistencia
					fila += formatoDato(completarSecuencia(0, 9), 9)                                    // Valor no retenido por aportes voluntarios
					fila += formatoDato("0.12500", 7)                                                   // Tarifa de aportes salud
					fila += formatoDato(completarSecuenciaString(pagoSalud, 9), 9)                      // Cotización obligatoria a salud

					valorUpc = buscarUpcAsociada(idPersona)
					fila += valorUpc                                 //Valor UPC Adicional
					fila += formatoDato("", 15)                      //Nº de autorización de la incapacidad por enfermedad general
					fila += formatoDato(completarSecuencia(0, 9), 9) //Valor de la incapacidad por enfermedad general
					fila += formatoDato("", 15)                      //Nº de autorización de la licencia de maternidad o paternidad
					fila += formatoDato(completarSecuencia(0, 9), 9) //Valor de la licencia de maternidad

					fila += formatoDato("0.0052200", 9) //Tarifa de aportes a Riegos Laborales

					fila += formatoDato("        1", 9)                          // Centro de trabajo CT
					fila += formatoDato(completarSecuenciaString(pagoArl, 9), 9) // Cotización obligatoria a sistema de riesgos laborales

					fila += formatoDato("0.04000", 7)                             // Tarifa de aportes CCF
					fila += formatoDato(completarSecuenciaString(pagoCaja, 9), 9) // Cotización obligatoria a salud

					fila += formatoDato(completarSecuencia(0, 7), 7) // Tarifa de aportes SENA
					fila += formatoDato(completarSecuencia(0, 9), 9) // Valor Aportes SENA

					fila += formatoDato("0.03000", 7)                             //Tarifa de aportes ICBF
					fila += formatoDato(completarSecuenciaString(pagoIcbf, 9), 9) // Cotización obligatoria a salud

					fila += formatoDato(completarSecuencia(0, 7), 7) //Tarifa de aportes ESAP
					fila += formatoDato(completarSecuencia(0, 9), 9) //Valor de aporte ESAP
					fila += formatoDato(completarSecuencia(0, 7), 7) //Tarifa de aportes MEN
					fila += formatoDato(completarSecuencia(0, 9), 9) //Valor de aporte MEN

					// Estos campos están vacios porque solo aplican a los registros que son upc
					fila += formatoDato("", 2)  // Tipo de documento del cotizante principal
					fila += formatoDato("", 16) // Número de identificación del cotizante principal

					fila += formatoDato("N", 1)     // Cotizante exonerado de pago de aporte salud, SENA e ICBF - Ley 1607 de 2012
					fila += formatoDato("14-23", 6) // Código de la administradora de Riesgos Laborales a la cual pertenece el afiliado
					fila += formatoDato("1", 1)     // Clase de Riesgo en la que se encuentra el afiliado
					fila += formatoDato("", 1)      // Indicador tarifa especial pensiones (Actividades de alto riesgo, Senadores, CTI y Aviadores aplican)

					fila += formatoDato("", 150) // Columnas que corresponden a las fechas de las novedades

					fila += formatoDato(completarSecuenciaString(ibcLiquidado, 9), 9) //IBC otros parafiscales difenrentes a CCF
					fila += formatoDato(horasLaboradas, 3)
					fila += formatoDato("", 26)

					fila += "\n" // siguiente persona...
					filas += fila
					secuencia++
					if suspencionTemporalContrato {
						crearFilaNovedad(idPersona, idPreliquidacion, value)
					} else if licenciaNoRemunerada {
						crearFilaNovedad(idPersona, idPreliquidacion, value)
					} else if comisionServicios {
						crearFilaNovedad(idPersona, idPreliquidacion, value)
					} else if incapacidadGeneral {
						crearFilaNovedad(idPersona, idPreliquidacion, value)
					} else if licenciaMaternidad {
						crearFilaNovedad(idPersona, idPreliquidacion, value)
					} else if vacaciones {
						crearFilaNovedad(idPersona, idPreliquidacion, value)
					} else if licenciaRemunerada {
						crearFilaNovedad(idPersona, idPreliquidacion, value)
					}
				}
			}
			respuestaJSON := make(map[string]interface{})
			respuestaJSON["informacion"] = filas
			c.Data["json"] = respuestaJSON
		} else {
			c.Data["json"] = err.Error()
		}

	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// crearFilaNovedad crea las filas de acuerdo a las novedades de una persona
func crearFilaNovedad(idPersona, idPreliquidacion string, persona models.InformacionPersonaNatural) {
	var horasLaboradasNovedad = "0"
	filaAux = ""
	filaAux += formatoDato(tipoRegistro, 2)                     //Tipo Registro
	filaAux += formatoDato(completarSecuencia(secuencia, 5), 5) //Secuencia

	filaAux += formatoDato("CC", 2)                     //Tip de documento del cotizante
	filaAux += formatoDato(persona.Id, 16)              //Número de identificación del cotizante
	filaAux += formatoDato(completarSecuencia(1, 2), 2) //Tipo Cotizante
	filaAux += formatoDato(completarSecuencia(0, 2), 2) //Subtipo de Cotizante
	filaAux += formatoDato("", 1)                       //Extranjero no obligado a cotizar pensión

	marcarNovedadConOpciones(exterior, "X") //Colombiano en el exterior
	filaAux += formatoDato("11", 2)         //Código del departamento de la ubicación laboral
	filaAux += formatoDato("001", 3)        //Código del municipio de ubicación laboral

	filaAux += formatoDato(persona.PrimerApellido, 20)  //Primer apellido
	filaAux += formatoDato(persona.SegundoApellido, 30) //Segundo apellido
	filaAux += formatoDato(persona.PrimerNombre, 20)    //Primer nombre
	filaAux += formatoDato(persona.SegundoNombre, 30)   //Segundo nombre

	marcarNovedadConOpciones(ingreso, "X")     //ING: Ingreso
	marcarNovedadConOpciones(retiro, "X")      //RET: Retiro
	marcarNovedad(trasladoDesdeEps)            //TDE: Traslado desde otra EPS o EOC
	marcarNovedad(trasladoAEps)                //TAE: Traslado a otra EPS o EOC
	marcarNovedad(trasladoDesdePensiones)      //TDP: Traslado desde otra administradora de pensiones
	marcarNovedad(trasladoAPensiones)          //TAP: Traslado a otra administradora de pensiones
	marcarNovedad(variacionPermanteSalario)    //VSP: Variación permanente de salario
	marcarNovedadConOpciones(corecciones, "A") //Corecciones
	marcarNovedad(variacionTransitoriaSalario) //VST: Variación transitoria del salario
	//SLN: Suspención temporal del contrato o licencia no remunerada o comisión de servicios
	if suspencionTemporalContrato || (licenciaNoRemunerada) {
		filaAux += formatoDato("X", 1)
	} else if comisionServicios {
		filaAux += formatoDato("C", 1)
	} else {
		filaAux += formatoDato("", 1)
	}

	marcarNovedad(incapacidadGeneral) //IGE: Incapacidad temporal por enfermadad general
	marcarNovedad(licenciaMaternidad) //LMA: Licencia de maternidad o de paternidad

	//VAC - LR: Vacaciones, licencia remunerada
	if vacaciones {
		filaAux += formatoDato("X", 1)
	} else if licenciaRemunerada {
		filaAux += formatoDato("L", 1)
	} else {
		filaAux += formatoDato("", 1)
	}
	marcarNovedad(aporteVoluntario)                           //APV: Aporte voluntario
	marcarNovedad(variacionCentroTrabajo)                     //VCT: Variación centros de trabajo
	filaAux += formatoDato(completarSecuencia(diasIrl, 2), 2) // IRL

	filaAux += codigoFondoPension

	//Código de la admnistradora de pensiones a la cual se traslada el afiliado
	// Si hay un translado, debe aparecer el nuevo código, de lo contrario será un campo vació
	if trasladoAPensiones {
		filaAux += formatoDato("230301", 6)
	} else {
		filaAux += formatoDato("", 6)
	}

	//Código EPS o EOC a la cual pertenece el afiliado
	if codigoEps == "      " {
		filaAux += formatoDato("MIN001", 6)
	} else {
		filaAux += codigoEps
	}

	//Código EPS o EOC a la cual se traslada el afiliado
	// Si hay un translado, debe aparecer el nuevo código, de lo contrario será un campo vació
	if trasladoAEps {
		filaAux += formatoDato("EPS012", 6)
	} else {
		filaAux += formatoDato("", 6)
	}

	filaAux += formatoDato("CCF24", 6) //Código CCF a la cual pertenece el afiliado

	if suspencionTemporalContrato {
		escribirDiasCotizadosAux(diasSuspencionContrato)
	} else if licenciaNoRemunerada {
		escribirDiasCotizadosAux(diasLicenciaNoRem)
	} else if comisionServicios {
		escribirDiasCotizadosAux(diasComisionServicios)
	} else if licenciaMaternidad {
		escribirDiasCotizadosAux(diasLicenciaMaternidad)
	} else if vacaciones {
		escribirDiasCotizadosAux(diasVacaciones)
	} else if licenciaRemunerada {
		escribirDiasCotizadosAux(diasLicenciaRem)
	}

	filaAux += formatoDato(completarSecuencia(salarioBase, 9), 9) //Salario básico
	filaAux += formatoDato("", 1)                                 //Salario integral

	filaAux += formatoDato(completarSecuencia(0, 36), 36) // Espacios para llegar a la tarifa de pensión
	filaAux += formatoDato(tarifaPension, 7)              //Tarifa de aportes pensiones

	filaAux += formatoDato(completarSecuencia(0, 63), 63) // Espacios para llegar a la tarifa de pensión
	filaAux += formatoDato(tarifaSalud, 7)                // Tarifa de aportes salud
	filaAux += formatoDato(completarSecuencia(0, 9), 9)   // Cotización obligatoria a salud

	filaAux += valorUpc                                 //Valor UPC Adicional
	filaAux += formatoDato("", 15)                      //Nº de autorización de la incapacidad por enfermedad general
	filaAux += formatoDato(completarSecuencia(0, 9), 9) //Valor de la incapacidad por enfermedad general
	filaAux += formatoDato("", 15)                      //Nº de autorización de la licencia de maternidad o paternidad
	filaAux += formatoDato(completarSecuencia(0, 9), 9) //Valor de la licencia de maternidad

	filaAux += formatoDato(tarifaArl, 9) //Tarifa de aportes a Riegos Laborales

	filaAux += formatoDato("        1", 9)              // Centro de trabajo CT
	filaAux += formatoDato(completarSecuencia(0, 9), 9) // Cotización obligatoria a sistema de riesgos laborales

	filaAux += formatoDato(tarifaCaja, 7)               // Tarifa de aportes CCF
	filaAux += formatoDato(completarSecuencia(0, 9), 9) // Cotización obligatoria a salud

	filaAux += formatoDato(completarSecuencia(0, 7), 7) // Tarifa de aportes SENA
	filaAux += formatoDato(completarSecuencia(0, 9), 9) // Valor Aportes SENA

	filaAux += formatoDato(tarifaIcbf, 7)               //Tarifa de aportes ICBF
	filaAux += formatoDato(completarSecuencia(0, 9), 9) // Cotización obligatoria a salud

	filaAux += formatoDato(completarSecuencia(0, 7), 7) //Tarifa de aportes ESAP
	filaAux += formatoDato(completarSecuencia(0, 9), 9) //Valor de aporte ESAP
	filaAux += formatoDato(completarSecuencia(0, 7), 7) //Tarifa de aportes MEN
	filaAux += formatoDato(completarSecuencia(0, 9), 9) //Valor de aporte MEN

	// Estos campos están vacios porque solo aplican a los registros que son upc
	filaAux += formatoDato("", 2)  // Tipo de documento del cotizante principal
	filaAux += formatoDato("", 16) // Número de identificación del cotizante principal

	filaAux += formatoDato("N", 1)     // Cotizante exonerado de pago de aporte salud, SENA e ICBF - Ley 1607 de 2012
	filaAux += formatoDato("14-23", 6) // Código de la administradora de Riesgos Laborales a la cual pertenece el afiliado
	filaAux += formatoDato("1", 1)     // Clase de Riesgo en la que se encuentra el afiliado
	filaAux += formatoDato("", 1)      // Indicador tarifa especial pensiones (Actividades de alto riesgo, Senadores, CTI y Aviadores aplican)

	//Fechas de novedades (AAAA-MM-DD)
	filaAux += formatoDato(fechaIngreso, 10)          //Fecha ingreso
	filaAux += formatoDato(fechaRetiro, 10)           //Fecha retiro
	filaAux += formatoDato(fechaInicioVsp, 10)        //Fecha inicio VSP
	filaAux += formatoDato(fechaInicioSuspencion, 10) //Fecha inicio SLN
	filaAux += formatoDato(fechaFinSuspencion, 10)    //Fecha fin SLN
	filaAux += formatoDato(fechaInicioIge, 10)        //Fecha inicio IGE
	filaAux += formatoDato(fechaFinIge, 10)           //Fecha fin IGE
	filaAux += formatoDato(fechaInicioLma, 10)        //Fecha inicio LMA
	filaAux += formatoDato(fechaFinLma, 10)           //Fecha fin LMA
	filaAux += formatoDato(fechaInicioVac, 10)        //Fecha inicio VAC-LR
	filaAux += formatoDato(fechaFinVac, 10)           //Fecha fin VAC-LR
	filaAux += formatoDato(fechaInicioVct, 10)        //Fecha inicio VCT
	filaAux += formatoDato(fechaFinVct, 10)           //Fecha fin VCT
	filaAux += formatoDato(fechaInicioIrl, 10)        //Fecha inicio IRL
	filaAux += formatoDato(fechaFinIrl, 10)           //Fecha fin IRL

	filaAux += formatoDato(completarSecuencia(0, 9), 9) //IBC otros parafiscales difenrentes a CCF
	filaAux += formatoDato(completarSecuenciaString(horasLaboradasNovedad, 3), 3)
	filaAux += formatoDato("", 26)

	filaAux += "\n"
	filas += filaAux
	secuencia++
}

// escribirDiasCotizadosAux escribe los días a cotizar para la fila auxiliar de novedades
func escribirDiasCotizadosAux(diasCotizados int) {
	filaAux += formatoDato(completarSecuencia(diasCotizados, 2), 2) //Número de días cotizados a pensión
	filaAux += formatoDato(completarSecuencia(diasCotizados, 2), 2) //Número de días cotizados a salud
	filaAux += formatoDato("00", 2)                                 //Número de días cotizados a ARL
	filaAux += formatoDato(completarSecuencia(diasCotizados, 2), 2) //Número de días cotizados a Caja de Compensación Familiar
}

// buscarUpcAsociada busca todas las upc asociadas a esa persona
// devuelve el valor total de las upc asociadas a esa persona
func buscarUpcAsociada(idPersona string) string {
	var beneficiariosAdicionales []models.UpcAdicional

	err := getJson("http://"+beego.AppConfig.String("segSocialService")+
		"/upc_adicional"+
		"?limit=0"+
		"&query=PersonaAsociada:"+idPersona, &beneficiariosAdicionales)
	if err != nil {
		ImprimirError("error en buscarUpcAsociada()", err)
	}
	valorTotalBeneficiariosAdc := 0
	if beneficiariosAdicionales[0].Id != 0 { // Esto es para verificar que el arreglo que devuelve el servicio tiene al menos un elemento con valor
		for _, beneficiario := range beneficiariosAdicionales {
			valorTotalBeneficiariosAdc += int(beneficiario.TipoUpc.Valor)
		}
	}
	return formatoDato(completarSecuencia(valorTotalBeneficiariosAdc, 9), 9)
}

// traerDiasCotizados función para trer los días cotizados, tipoPago debe tener el mismo nombre que un concepto_nomina en titán
// duelve:
// diasLiquidados: días coritzados del tipo de pago
func traerDiasCotizados(idPersona, idPreliquidacion, tipoPago string) string {

	var detallePreliquidacion []models.DetallePreliquidacion

	err := getJson("http://"+beego.AppConfig.String("titanServicio")+
		"/detalle_preliquidacion"+
		"?limit=1"+
		"&query=Persona:"+idPersona+
		",Preliquidacion.Id:"+idPreliquidacion+
		",Concepto.NombreConcepto:"+tipoPago, &detallePreliquidacion)

	if err != nil {
		ImprimirError("error en traerDiasCotizados()", err)
	}

	dias := 0
	if detallePreliquidacion != nil {
		dias = int(detallePreliquidacion[0].DiasLiquidados)
	}
	diasLiquidados := formatoDato(completarSecuencia(dias, 2), 2)
	return diasLiquidados
}

// traerValorConceptoEmpleado trae el valor de un concepto que tenga que pagar el empleado
func traerValorConceptoEmpleado(idPersona, idPreliquidacion, tipoPago string) string {
	var (
		detallePreliquidacion []models.DetallePreliquidacion
	)

	err := getJson("http://"+beego.AppConfig.String("titanServicio")+
		"/detalle_preliquidacion"+
		"?limit=1"+
		"&query=Persona:"+idPersona+
		",Preliquidacion.Id:"+idPreliquidacion+
		",Concepto.NombreConcepto:"+tipoPago, &detallePreliquidacion)
	if err != nil {
		ImprimirError("error en traerDiasCotizadosEmpleador()", err)
	}
	valorTipoPago := 0.0
	if detallePreliquidacion != nil {
		valorTipoPago = detallePreliquidacion[0].ValorCalculado
	}
	valorPagoAproximado := AproximarPesoSuperior(valorTipoPago, 100)
	return formatoDato(completarSecuencia(valorPagoAproximado, 9), 9)
}

// traerDiasCotizadosEmpleador trae los días cotizados correspondientes al pago del emplador
// devuelve
// diasCotizados: días cotizados del tipo de pago
// valorTipoPago: valor del tipo de pago por parte de la universidad
// valorTotalPago: valor total del tipo de pago (sumando lo del empleado y la universidad)
func traerDiasCotizadosEmpleador(idPersona, idPreliquidacion, periodoPago, tipoPago string) (string, string) {
	var (
		pago                  []models.Pago
		detallePreliquidacion []models.DetallePreliquidacion
		conceptoNomina        []models.ConceptoNomina
	)

	valorTipoPagoTemp := GetPagoEmpleado(idPersona, idPreliquidacion, tipoPago)

	err := getJson("http://"+beego.AppConfig.String("titanServicio")+
		"/detalle_preliquidacion"+
		"?limit=1"+
		"&query=Persona:"+idPersona+
		",Preliquidacion.Id:"+idPreliquidacion+
		",Concepto.NombreConcepto:ibc_liquidado", &detallePreliquidacion)
	if err != nil {
		ImprimirError("error en traerDiasCotizadosEmpleador()", err)
	}

	err = getJson("http://"+beego.AppConfig.String("titanServicio")+
		"/concepto_nomina"+
		"?limit=1"+
		"&query=NaturalezaConcepto.Nombre:seguridad_social"+
		",NombreConcepto:"+tipoPago, &conceptoNomina)
	if err != nil {
		ImprimirError("error en traerDiasCotizadosEmpleador()", err)
	}

	err = getJson("http://"+beego.AppConfig.String("segSocialService")+
		"/pago"+
		"?limit=1"+
		"&query=DetalleLiquidacion:"+strconv.Itoa(detallePreliquidacion[0].Id)+
		",PeriodoPago.Id:"+periodoPago+
		",TipoPago:"+strconv.Itoa(conceptoNomina[0].Id), &pago)
	if err != nil {
		ImprimirError("error en traerDiasCotizadosEmpleador()", err)
	}

	diasCotizados := "00"
	valorTotalPago := formatoDato(completarSecuencia(valorTipoPagoTemp+int(pago[0].Valor), 9), 9)

	if pago[0].EntidadPago != 0 {
		diasCotizados = formatoDato(completarSecuenciaString(traerDiasCotizados(idPersona, idPreliquidacion, "salud"), 2), 2)
	}
	return diasCotizados, valorTotalPago
}

func establecerNovedadesExterior(idPersona, idPreliquidacion string) {
	var conceptoNominaPorPersona []models.ConceptoNominaPorPersona

	err := getJson("http://"+beego.AppConfig.String("titanServicio")+
		"/detalle_preliquidacion"+
		"?limit=0"+
		"&query=Persona:"+idPersona+
		",Preliquidacion.Id:"+idPreliquidacion+
		",Concepto.NombreConcepto__startswith:exterior", &conceptoNominaPorPersona)

	if err != nil {
		ImprimirError("error en establecerNovedadesExterior()", err)
	}

	if len(conceptoNominaPorPersona) > 0 {
		fila += formatoDato("X", 1) //Colombiano en el exterior
		novedadPersona = true
		exterior = true
	} else {
		fila += formatoDato("", 1) //Colombiano en el exterior
	}

	fila += formatoDato("11", 2)  //Código del departamento de la ubicación laboral
	fila += formatoDato("001", 3) //Código del municipio de ubicación laboral
}

// establecerNovedades se encarga de buscar todas las novedades de una persona en una preliquidación especifica
// y recursivamente va a llenando las novedades, para luego marcarlas en el archivo y también asignarles las fechas correspondientes
func establecerNovedades(idPersona, idPreliquidacion, cedulaPersona string) {
	reinicializarVariablesNovedades()

	var (
		detallePreliquidaicon []models.DetallePreliquidacion
		conceptoNominaPersona []models.ConceptoNominaPorPersona
		fechaInicioTemp,
		fechaFinTemp string
	)

	// diasIbc := 0
	// diasIbcNovedad := 0

	fechaIngresoTemp := revisarIngreso(idPreliquidacion, cedulaPersona)

	err := getJson("http://"+beego.AppConfig.String("titanServicio")+
		"/detalle_preliquidacion"+
		"?limit=0"+
		"&query=Persona:"+idPersona+
		",Preliquidacion.Id:"+idPreliquidacion, &detallePreliquidaicon)

	if err != nil {
		ImprimirError("error en establecerNovedades()", err)
	}

	if int(fechaIngresoTemp.Month()) == detallePreliquidaicon[0].Preliquidacion.Mes {
		ingreso = true
		fechaIngreso = fechaIngresoTemp.String()
	}

	for _, value := range detallePreliquidaicon {

		err := getJson("http://"+beego.AppConfig.String("titanServicio")+
			"/concepto_nomina_por_persona"+
			"?limit=1"+
			"&query=Persona:"+idPersona+
			",Concepto.NombreConcepto:"+value.Concepto.NombreConcepto+
			",Activo:true", &conceptoNominaPersona)
		if err != nil {
			ImprimirError("error en establecerNovedades()", err)
		}
		if conceptoNominaPersona != nil {

			auxFechaInicio := conceptoNominaPersona[0].FechaDesde
			if int(auxFechaInicio.Month()) < mesPeriodo {
				fechaInicioTemp = time.Date(auxFechaInicio.Year(), time.Month(mesPeriodo), 1, 0, 0, 0, 0, time.UTC).Format(formatoFecha)
			} else {
				fechaInicioTemp = conceptoNominaPersona[0].FechaDesde.Format(formatoFecha)
			}

			auxFechaFin := conceptoNominaPersona[0].FechaHasta
			if int(auxFechaFin.Month()) > mesPeriodo {
				fechaFinTemp = time.Date(auxFechaFin.Year(), time.Month(mesPeriodo), 30, 0, 0, 0, 0, time.UTC).Format(formatoFecha)
			} else {
				fechaFinTemp = conceptoNominaPersona[0].FechaHasta.Format(formatoFecha)
			}
		}

		switch value.Concepto.NombreConcepto {
		case "licencia_rem":
			licenciaRemunerada = true
			diasLicenciaRem = int(value.DiasLiquidados)
			fechaInicioVac = fechaInicioTemp
			fechaFinVac = fechaFinTemp
			novedadPersona = true
		case "vacaciones":
			vacaciones = true
			diasArl = int(value.DiasLiquidados)
			fechaInicioVac = fechaInicioTemp
			fechaFinVac = fechaFinTemp
			novedadPersona = true
		case "incapacidad_general":
			diasIncapacidad = int(value.DiasLiquidados)
		case "incapacidad_laboral":
			incapacidadGeneral = true
			diasIrl = int(value.DiasLiquidados)
			diasArl = diasIrl
			fechaInicioIge = fechaInicioTemp
			fechaFinIge = fechaFinTemp
			novedadPersona = true
		case "licencia_maternidad":
		case "licencia_paternidad":
			licenciaMaternidad = true
			diasLicenciaMaternidad = int(value.DiasLiquidados)
			diasArl = int(value.DiasLiquidados)
			fechaInicioLma = fechaInicioTemp
			fechaFinLma = fechaFinTemp
			novedadPersona = true
		case "comision_estudio":
		case "comision_norem":
		case "licencia_norem":
			switch value.Concepto.NombreConcepto {
			case "comision_estudio":
			case "comision_norem":
				diasComisionServicios = int(value.DiasLiquidados)
				comisionServicios = true
			case "licencia_norem":
				diasLicenciaNoRem = int(value.DiasLiquidados)
				licenciaNoRemunerada = true
			}
			fechaInicioSuspencion = fechaInicioTemp
			fechaFinSuspencion = fechaFinTemp
			tarifaArl = completarSecuencia(0, 9)
			tarifaIcbf = completarSecuencia(0, 7)
			tarifaCaja = completarSecuencia(0, 7)
			tarifaSalud = "0.08500"
			tarifaPension = "0.12000"
			novedadPersona = true
		}
	}

	fila += formatoDato("", 8)                             // son espacios del archivo plano hasta la novedad de VST
	fila += formatoDato("X", 1)                            // VST: Variación transitoria de salario
	fila += formatoDato("", 6)                             //son espacios del archivo plano hasta el código de la adminitradora de fondo de pensiones
	fila += formatoDato(completarSecuencia(diasIrl, 2), 2) // IRL
}

func revisarIngreso(idPreliquidacion, cedulaPersona string) (fechaMenor time.Time) {
	var preliquidacion models.Preliquidacion
	var contratosPersona map[string]map[string]interface{}

	err := getJson("http://"+beego.AppConfig.String("titanServicio")+
		"/preliquidacion/"+idPreliquidacion, &preliquidacion)
	if err != nil {
		ImprimirError("error en revisarIngreso()", err)
	}

	anio := fmt.Sprint(preliquidacion.Ano)
	mes := fmt.Sprint(preliquidacion.Mes)
	if preliquidacion.Mes < 9 {
		mes = "0" + fmt.Sprint(preliquidacion.Mes)
	}

	if contratistas {
		// Contrato de prestación de servicios profesionales o apoyo a la gestión
		err = getJsonWSO2("http://"+beego.AppConfig.String("argoWso2Service")+
			"/contratos_elaborado_tipo_persona/6/"+anio+"-"+mes+"/"+anio+"-"+mes+"/"+cedulaPersona, &contratosPersona)
		fmt.Println("Contratista: http://" + beego.AppConfig.String("argoWso2Service") +
			"/contratos_elaborado_tipo_persona/6/" + anio + "-" + mes + "/" + anio + "-" + mes + "/" + cedulaPersona)
		if err != nil {
			ImprimirError("error en revisarIngreso()", err)
		}

	} else {
		// Contrato de docente de vinculación especial (salarios)
		err = getJsonWSO2("http://"+beego.AppConfig.String("argoWso2Service")+
			"/contratos_elaborado_tipo_persona/2/"+anio+"-"+mes+"/"+anio+"-"+mes+"/"+cedulaPersona, &contratosPersona)
		fmt.Println("vinculación especial salarios: http://" + beego.AppConfig.String("argoWso2Service") +
			"/contratos_elaborado_tipo_persona/2/" + anio + "-" + mes + "/" + anio + "-" + mes + "/" + cedulaPersona)
		if err != nil {
			ImprimirError("error en revisarIngreso()", err)
		}
		// Contrato docente de vinculación especial (Tiempo completo ocasional TCO - MTO)
		if len(contratosPersona["contratos_tipo"]) == 0 {
			err = getJsonWSO2("http://"+beego.AppConfig.String("argoWso2Service")+
				"/contratos_elaborado_tipo_persona/18/"+anio+"-"+mes+"/"+anio+"-"+mes+"/"+cedulaPersona, &contratosPersona)
			fmt.Println("Vinculación especial TCO: http://" + beego.AppConfig.String("argoWso2Service") +
				"/contratos_elaborado_tipo_persona/18/" + anio + "-" + mes + "/" + anio + "-" + mes + "/" + cedulaPersona)
			if err != nil {
				ImprimirError("error en revisarIngreso()", err)
			}
		}
	}

	if contratosPersona != nil {
		var actaInicio map[string]map[string]interface{}
		for _, value := range contratosPersona {
			contratos := value["contrato_tipo"].([]interface{})

			for _, contrato := range contratos {
				numeroContrato := contrato.(map[string]interface{})["numero_contrato"].(string)
				vigenciaContrato := contrato.(map[string]interface{})["vigencia"].(string)

				err := getJsonWSO2("http://"+beego.AppConfig.String("argoWso2Service")+
					"/acta_inicio_elaborado/"+numeroContrato+"/"+vigenciaContrato, &actaInicio)
				fmt.Println("Acta inicio: http://" + beego.AppConfig.String("argoWso2Service") +
					"/acta_inicio_elaborado/" + numeroContrato + "/" + vigenciaContrato)
				if err != nil {
					ImprimirError("error en revisarIngreso()", err)
				}

				if len(actaInicio["actaInicio"]) != 0 {

					t, err := time.Parse(formatoFecha, actaInicio["actaInicio"]["fechaInicio"].(string))
					if err != nil {
						ImprimirError("error en revisarIngreso()", err)
					}

					if fechaMenor.IsZero() {
						fechaMenor = t
					} else {
						if t.Year() <= fechaMenor.Year() {
							if t.Month() <= fechaMenor.Month() {
								if t.Day() <= fechaMenor.Day() {
									fechaMenor = t
								} else if (t.Year() <= fechaMenor.Year()) && (t.Month() <= fechaMenor.Month()) && (t.Day() >= fechaMenor.Day()) {
									fechaMenor = t
								}
							} else if (t.Year() <= fechaMenor.Year()) && (t.Month() >= fechaMenor.Month()) {
								fechaMenor = t
							}
						}
					}
				}

			}

		}
	}
	return
}

// marcarNovedad marca el valor para la novedad o la deja vacia
func marcarNovedadConOpciones(novedad bool, valor string) {
	if novedad {
		filaAux += formatoDato(valor, 1)
	} else {
		filaAux += formatoDato("", 1)
	}
}

// marcarNovedad marca el X para la novedad o la deja vacia
func marcarNovedad(novedad bool) {
	if novedad {
		filaAux += formatoDato("X", 1)
	} else {
		filaAux += formatoDato("", 1)
	}
}

func traerCodigoEntidadSalud(idEntidad string) string {
	var infoAPI map[string]interface{}

	// Este es el servicio que se debería consumir....
	// err := getJson("http://"+beego.AppConfig.String("core_api")+
	// 	"/parametro_entidad/"+idEntidad, &infoAPI)

	err := getJson("http://"+beego.AppConfig.String("segSocialService")+
		"/pago/GetParametroEntidad/"+idEntidad, &infoAPI)

	if err != nil {
		ImprimirError("error en traerCodigoEntidadSalud()", err)
	}
	codigoEntiad := ""
	if infoAPI["Codigo"] != nil {
		codigoEntiad = formatoDato(infoAPI["Codigo"].(string), 6)
	}
	codigoEntiadFormato := formatoDato(codigoEntiad, 6)
	return codigoEntiadFormato
}

// GetPagoEmpleado obtiene los pagos correspondientes de salud y pensión del empleado
func GetPagoEmpleado(idPersona, idPreliquidacion, tipoPago string) (valorPago int) {
	var detallePreliquidacion []models.DetallePreliquidacion
	switch tipoPago {
	case "salud_ud":
		err := getJson("http://"+beego.AppConfig.String("titanServicio")+
			"/detalle_preliquidacion"+
			"?limit=1"+
			"&query=Persona:"+idPersona+
			",Preliquidacion.Id:"+idPreliquidacion+
			",Concepto.NombreConcepto:salud", &detallePreliquidacion)
		if err != nil {
			ImprimirError("error en traerDiasCotizadosEmpleador()", err)
		}
		valorPago = AproximarPesoSuperior(detallePreliquidacion[0].ValorCalculado, 100)
	case "pension_ud":
		err := getJson("http://"+beego.AppConfig.String("titanServicio")+
			"/detalle_preliquidacion"+
			"?limit=1"+
			"&query=Persona:"+idPersona+
			",Preliquidacion.Id:"+idPreliquidacion+
			",Concepto.NombreConcepto:salud", &detallePreliquidacion)
		if err != nil {
			ImprimirError("error en traerDiasCotizadosEmpleador()", err)
		}
		valorPago = AproximarPesoSuperior(detallePreliquidacion[0].ValorCalculado, 100)
	}
	return
}

// AproximarPesoSuperior apróxima un número al peso superior
func AproximarPesoSuperior(valor float64, valorAaproximar int) int {
	x := valor / float64(valorAaproximar)
	y := math.Trunc(x)

	var numero int
	if (x - y) > 0 {
		numero = AproximarPesoSuperior(float64(math.Trunc(valor)+1), valorAaproximar)

	} else {
		numero = int(valor)
	}
	return numero
}

func establecerNovedadesTranslado() {
}

func completarSecuencia(num, cantSecuencia int) (secuencia string) {
	tamanioNum := len(strconv.Itoa(num))
	for i := 0; i < cantSecuencia-tamanioNum; i++ {
		secuencia += "0"
	}
	secuencia += strconv.Itoa(num)
	return
}

// reinicializarVariablesNovedades simplemente devuelva al valor original las variables declaras al comienzo, para que no se repitan en
// otras personas
func reinicializarVariablesNovedades() {
	// Reinicializando valor de las variables para cada una de las novedades y sus días validos
	ingreso = false
	fechaIngreso = ""

	retiro = false
	fechaRetiro = ""

	trasladoDesdeEps = false
	trasladoDesdePensiones = false
	trasladoAPensiones = false
	trasladoAEps = false
	variacionPermanteSalario = false
	corecciones = false
	if tipoPreliquidacion == "HCS" {
		variacionTransitoriaSalario = true
	} else {
		variacionTransitoriaSalario = false
	}

	suspencionTemporalContrato = false

	exterior = false

	fechaInicioSuspencion = ""
	fechaFinSuspencion = ""

	licenciaNoRemunerada = false
	comisionServicios = false
	incapacidadGeneral = false
	fechaInicioIge = ""
	fechaFinIge = ""

	licenciaMaternidad = false
	fechaInicioLma = ""
	fechaFinLma = ""

	vacaciones = false
	licenciaRemunerada = false
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
	diasArl = 30
	tarifaIcbf = "0.03000"
	tarifaCaja = "0.04000"
	tarifaArl = "0.0052200"
	tarifaPension = "0.16000"
	tarifaSalud = "0.12500"

	novedadPersona = false
	diasIrl = 0
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
