//Comienza la corrección...

package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/udistrital/ss_mid_api/golog"
	"github.com/udistrital/ss_mid_api/models"
)

// PlanillasController operations for Planillas
type PlanillasController struct {
	beego.Controller
}

// URLMapping ...
func (c *PlanillasController) URLMapping() {
	c.Mapping("GenerarPlanillaActivos", c.GenerarPlanillaActivos)
	c.Mapping("PruebaPlanilla", c.PruebaPlanilla)
}

var (
	detallePreliquidacion []models.DetallePreliquidacion // detalle de toda la preliquidación
	contratistas          = false

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
	licenciaPaternidad = false
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
	diasNovedad    int
	novedadPersona = false
	diasIrl        = 0

	fechaInicioNovedad,
	fechaFinNovedad time.Time

	tipoRegistro = "02"
	secuencia    = 0

	contratosElaboradosPeriodo map[string]time.Time
	actasInicioPeriodo         map[string]map[string]interface{}
)

func getInfoContratosElaboradoTipoFecha(anioPeriodo, mesPeriodo int, tipoLiquidacion string) (err error) {
	var contratos, actasInicio map[string]map[string][]map[string]interface{}
	var idTipoContrato string

	anio := fmt.Sprint(anioPeriodo)
	mes := fmt.Sprint(mesPeriodo)
	fechaMenor := time.Now()
	contratosElaboradosPeriodo = make(map[string]time.Time)

	switch tipoLiquidacion {
	case "DP":
		idTipoContrato = "21"
	case "FP":
		idTipoContrato = "20"
	case "HCH":
		idTipoContrato = "3"
	case "HCS":
		idTipoContrato = "18"
	case "CT":
		idTipoContrato = "6"
	}

	err = getJsonWSO2("http://"+beego.AppConfig.String("argoWso2Service")+"/contratos_elaborado_tipo_fecha/"+
		idTipoContrato+"/"+anio+"-"+mes+"/"+anio+"-"+mes,
		&contratos)
	if err != nil {
		ImprimirError("error en getInfoContratosElaboradoTipoFecha():", err)
		return
	}

	err = getJsonWSO2("http://"+beego.AppConfig.String("argoWso2Service")+"/acta_inicio_elaborado_vigencia/"+anio, &actasInicio)
	if err != nil {
		ImprimirError("error en getInfoContratosElaboradoTipoFecha():", err)
		return
	}

	for _, valueContratos := range contratos["contratos_tipo"]["contrato_tipo"] {
		for _, valueActas := range actasInicio["contratos"]["contrato"] {
			if valueContratos["numero_contrato"].(string) == valueActas["numeroContrato"].(string) {

				t, _ := time.Parse(formatoFecha, valueActas["fechaInicio"].(string))
				if !t.IsZero() {
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

				contratosElaboradosPeriodo[valueContratos["numero_documento"].(string)] = fechaMenor
				break
			}
		}
	}

	return
}

func getValorConcepto(preliquidacion string, concepto string) (map[string]int, error) {
	var detallesPreliquidacion []models.DetallePreliquidacion

	salariosBase := make(map[string]int)

	err := getJson("http://"+beego.AppConfig.String("titanServicio")+
		"/detalle_preliquidacion?limit=-1"+
		"&fields=ValorCalculado,Persona"+
		"&query=Preliquidacion:"+preliquidacion+",Concepto.NombreConcepto:"+concepto, &detallesPreliquidacion)

	if err != nil {
		return nil, err
	}

	for _, value := range detallesPreliquidacion {
		// fmt.Println("value: ", value)
		salariosBase[strconv.Itoa(value.Persona)] += int(value.ValorCalculado)
	}
	return salariosBase, nil
}

// PruebaPlanilla ...
// @Title Generar planilla de activos
// @Description Recibe un periodo pago y devuelve un arreglo de json con la información de la planilla
// @Param	body body PeriodoPago true	"body for PeriodoPago"
// @Success 200 {string} string
// @Failure 403 body is empty
// @router /PruebaPlanilla/:limit [get]
func (c *PlanillasController) PruebaPlanilla() {
	limit := c.Ctx.Input.Param(":limit")
	num := rand.Intn(100)
	fmt.Println(num)
	fmt.Println(limit)
	c.Data["json"] = map[string]int{"test": num}
	c.ServeJSON()
}

// GenerarPlanillaActivos ...
// @Title Generar planilla de activos
// @Description Recibe un periodo pago y devuelve un arreglo de json con la información de la planilla
// @Param	body body PeriodoPago true	"body for PeriodoPago"
// @Success 200 {string} string
// @Failure 403 body is empty
// @router /GenerarPlanillaActivos/:limit/:offset [post]
func (c *PlanillasController) GenerarPlanillaActivos() {
	secuencia = 1
	start := time.Now()

	var filasPlanilla []models.PlanillaTipoE

	limit := c.Ctx.Input.Param(":limit")
	offset := c.Ctx.Input.Param(":offset")

	log.Println(limit, offset)

	log.Println("Comenzó a generar la planilla")
	var (
		periodoPago *models.PeriodoPago
		personas    []string
	)

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &periodoPago); err == nil {

		tipoPreliquidacion = periodoPago.TipoLiquidacion
		mesPeriodo = int(periodoPago.Mes)
		anioPeriodo = int(periodoPago.Anio)

		err = getInfoContratosElaboradoTipoFecha(anioPeriodo, mesPeriodo, periodoPago.TipoLiquidacion)
		if err != nil {
			c.Data["json"] = map[string]string{"error": err.Error()}
			log.Println("error: ", err.Error())
			c.ServeJSON()
			return
		}

		if err = getJson("http://"+beego.AppConfig.String("titanServicio")+"/detalle_preliquidacion?"+
			"limit="+limit+
			"&query=Preliquidacion.Id:"+strconv.Itoa(periodoPago.Liquidacion)+
			",Concepto.NombreConcepto:ibc_liquidado"+
			"&sortby=Persona&order=asc&offset="+offset, &detallePreliquidacion); err == nil {

			if periodoPago.TipoLiquidacion == "CT" || periodoPago.TipoLiquidacion == "HCH" {
				contratistas = true
			}

			concepto := "salarioBase"
			if contratistas {
				concepto = "honorarios"
			}

			idPreliquidacion := strconv.Itoa(detallePreliquidacion[0].Preliquidacion.Id)

			salariosBase, err := getValorConcepto(strconv.Itoa(periodoPago.Liquidacion), concepto)
			if err != nil {
				c.Data["json"] = map[string]string{"error": err.Error()}
				c.ServeJSON()
				return
			}

			ingresosCotizacion, err := getValorConcepto(strconv.Itoa(periodoPago.Liquidacion), "ibc_liquidado")
			if err != nil {
				c.Data["json"] = map[string]string{"error": err.Error()}
				c.ServeJSON()
				return
			}
			idPeriodoPago := strconv.Itoa(periodoPago.Id)

			informacionPagosPension, err := traerDiasCotizadosEmpleador(idPreliquidacion, idPeriodoPago, "pension_ud")
			if err != nil {
				c.Data["json"] = map[string]string{"error": err.Error()}
				c.ServeJSON()
				return
			}

			informacionPagosPensionVoluntaria, err := traerDiasCotizadosEmpleador(idPreliquidacion, idPeriodoPago, "nombreRegla2176")
			if err != nil {
				c.Data["json"] = map[string]string{"error": err.Error()}
				c.ServeJSON()
				return
			}

			informacionPagosSalud, err := traerDiasCotizadosEmpleador(idPreliquidacion, idPeriodoPago, "salud_ud")
			if err != nil {
				c.Data["json"] = map[string]string{"error": err.Error()}
				c.ServeJSON()
				return
			}
			informacionPagosIcbf, err := traerDiasCotizadosEmpleador(idPreliquidacion, idPeriodoPago, "icbf")
			if err != nil {
				c.Data["json"] = map[string]string{"error": err.Error()}
				c.ServeJSON()
				return
			}
			informacionPagosArl, err := traerDiasCotizadosEmpleador(idPreliquidacion, idPeriodoPago, "arl")
			if err != nil {
				c.Data["json"] = map[string]string{"error": err.Error()}
				c.ServeJSON()
				return
			}
			informacionPagosCaja, err := traerDiasCotizadosEmpleador(idPreliquidacion, idPeriodoPago, "caja_compensacion")
			if err != nil {
				c.Data["json"] = map[string]string{"error": err.Error()}
				c.ServeJSON()
				return
			}
			informacionPagosFondoSoli, err := traerDiasCotizadosEmpleador(idPreliquidacion, idPeriodoPago, "fondoSolidaridad")
			if err != nil {
				c.Data["json"] = map[string]string{"error": err.Error()}
				c.ServeJSON()
				return
			}

			informacionUpc, err := buscarUpcAsociada()
			if err != nil {
				c.Data["json"] = map[string]string{"error": err.Error()}
				c.ServeJSON()
				return
			}

			filas = ""

			for i := range detallePreliquidacion {
				personas = append(personas, fmt.Sprint(detallePreliquidacion[i].Persona))
			}
			mapPersonas, err := GetInfoPersonas(detallePreliquidacion)
			if err != nil {
				c.Data["json"] = map[string]string{"error": err.Error()}
				c.ServeJSON()
				return
			}

			contPersonas := 0

			for key, value := range mapPersonas {
				var (
				// preliquidacion []models.DetallePreliquidacion
				)
				idPersona := key //idProveedor

				cedulaPersona := fmt.Sprint(value.Id)
				fila = ""
				filaAux = ""

				diasLaborados := traerDiasCotizados(idPersona, idPreliquidacion, "salud")
				if err != nil {
					c.Data["json"] = map[string]string{"error": err.Error()}
					c.ServeJSON()
					return
				}

				filaPlanilla := models.PlanillaTipoE{
					TipoRegistro:                    models.Columna{Valor: "02", Longitud: 2},
					TipoDocumento:                   models.Columna{Valor: "CC", Longitud: 2},
					NumeroIdentificacion:            models.Columna{Valor: value.Id, Longitud: 16},
					TipoCotizante:                   models.Columna{Valor: 1, Longitud: 2},
					SubTipoCotizante:                models.Columna{Valor: 0, Longitud: 2},
					ExtranjeroNoPension:             models.Columna{Valor: "", Longitud: 1},
					PrimerApellido:                  models.Columna{Valor: value.PrimerApellido, Longitud: 20},
					SegundoApellido:                 models.Columna{Valor: value.SegundoApellido, Longitud: 30},
					PrimerNombre:                    models.Columna{Valor: value.PrimerNombre, Longitud: 20},
					SegundoNombre:                   models.Columna{Valor: value.SegundoNombre, Longitud: 30},
					CodigoFondoPension:              models.Columna{Valor: traerCodigoEntidadSalud(strconv.Itoa(value.IdFondoPension)), Longitud: 6},
					TrasladoPension:                 models.Columna{Valor: "", Longitud: 6},
					CodigoEps:                       models.Columna{Valor: traerCodigoEntidadSalud(strconv.Itoa(value.IdEps)), Longitud: 6},
					TrasladoEps:                     models.Columna{Valor: "", Longitud: 6},
					CodigoCCF:                       models.Columna{Valor: "CCF24", Longitud: 6},
					DiasLaborados:                   models.Columna{Valor: diasLaborados, Longitud: 2},
					DiasPension:                     models.Columna{Valor: diasLaborados, Longitud: 2},
					DiasSalud:                       models.Columna{Valor: traerDiasCotizados(idPersona, idPreliquidacion, "salud"), Longitud: 2},
					DiasArl:                         models.Columna{Valor: informacionPagosArl[key]["dias"], Longitud: 2},
					DiasCaja:                        models.Columna{Valor: informacionPagosCaja[key]["dias"], Longitud: 2},
					SalarioBase:                     models.Columna{Valor: salariosBase[key], Longitud: 9},
					SalarioIntegral:                 models.Columna{Valor: "", Longitud: 1},
					Ibcension:                       models.Columna{Valor: ingresosCotizacion[key], Longitud: 9},
					IbcSalud:                        models.Columna{Valor: ingresosCotizacion[key], Longitud: 9},
					IbcArl:                          models.Columna{Valor: ingresosCotizacion[key], Longitud: 9},
					IbcCcf:                          models.Columna{Valor: ingresosCotizacion[key], Longitud: 9},
					TarifaPension:                   models.Columna{Valor: "0.16000", Longitud: 7},
					PagoPension:                     models.Columna{Valor: informacionPagosPension[key]["valor"], Longitud: 9},
					AportePension:                   models.Columna{Valor: informacionPagosPensionVoluntaria[key]["valor"], Longitud: 9},
					TotalPension:                    models.Columna{Valor: informacionPagosPension[key]["valor"] + informacionPagosPensionVoluntaria[key]["valor"], Longitud: 9},
					FondoSolidaridad:                models.Columna{Valor: informacionPagosFondoSoli[key]["valor"], Longitud: 9},
					FondoSubsistencia:               models.Columna{Valor: informacionPagosFondoSoli[key]["valor"], Longitud: 9},
					NoRetenidoAportesVolunarios:     models.Columna{Valor: 0, Longitud: 9},
					TarifaSalud:                     models.Columna{Valor: "0.12500", Longitud: 7},
					PagoSalud:                       models.Columna{Valor: informacionPagosSalud[key]["valor"], Longitud: 9},
					ValorUpc:                        models.Columna{Valor: informacionUpc[key], Longitud: 9},
					AutorizacionEnfermedadGeneral:   models.Columna{Valor: "", Longitud: 15},
					ValorIncapacidadGeneral:         models.Columna{Valor: 0, Longitud: 15},
					AutotizacionLicenciaMarternidad: models.Columna{Valor: "", Longitud: 15},
					ValorLicenciaMaternidad:         models.Columna{Valor: 0, Longitud: 15},
					TarifaArl:                       models.Columna{Valor: "0.0052200", Longitud: 9},
					CentroTrabajo:                   models.Columna{Valor: "1", Longitud: 9},
					PagoArl:                         models.Columna{Valor: informacionPagosArl[key]["valor"], Longitud: 9},
					TarifaCaja:                      models.Columna{Valor: "0.04000", Longitud: 7},
					PagoCaja:                        models.Columna{Valor: informacionPagosCaja[key]["valor"], Longitud: 9},
					TarifaSena:                      models.Columna{Valor: 0, Longitud: 7},
					PagoSena:                        models.Columna{Valor: 0, Longitud: 9},
					TarifaIcbf:                      models.Columna{Valor: "0.03000", Longitud: 7},
					PagoIcbf:                        models.Columna{Valor: informacionPagosIcbf[key]["valor"], Longitud: 9},
					TarifaEsap:                      models.Columna{Valor: 0, Longitud: 9},
					PagoEsap:                        models.Columna{Valor: 0, Longitud: 9},
					TarifaMen:                       models.Columna{Valor: 0, Longitud: 9},
					PagoMen:                         models.Columna{Valor: 0, Longitud: 9},
					TipoDocumentoCotizantePrincipal: models.Columna{Valor: "", Longitud: 2},
					DocumentoCotizantePrincipal:     models.Columna{Valor: "", Longitud: 16},
					ExoneradoPagoSalud:              models.Columna{Valor: "N", Longitud: 1},
					CodigoArl:                       models.Columna{Valor: "14-23", Longitud: 6},
					ClaseRiesgo:                     models.Columna{Valor: "1", Longitud: 1},
					IndicadorTarifaEspecialPension:  models.Columna{Valor: "", Longitud: 1},
					FechasNovedades:                 models.Columna{Valor: "", Longitud: 150},
					IbcOtrosParaFiscales:            models.Columna{Valor: ingresosCotizacion[key], Longitud: 9},
					HorasLaboradas:                  models.Columna{Valor: diasLaborados * 8, Longitud: 3},
					EspacioBlanco:                   models.Columna{Valor: "", Longitud: 26},
				}

				log.Println("filaPlanilla: ", filaPlanilla.NumeroIdentificacion)

				filasPlanilla = append(filasPlanilla, filaPlanilla)
				contPersonas++
				log.Println("contPersonas:", contPersonas)

				establecerNovedadesExterior(idPersona, idPreliquidacion)

				establecerNovedades(idPersona, idPreliquidacion, cedulaPersona)

				// // Aporte voluntario del afiliado al fondo de pensiones
				// err = getJson("http://"+beego.AppConfig.String("titanServicio")+
				// 	"/detalle_preliquidacion"+
				// 	"?fields=Concepto,Id"+
				// 	"&query=Preliquidacion:"+strconv.Itoa(periodoPago.Liquidacion)+",Persona:"+key, &preliquidacion)

				// if err != nil {
				// 	fila += formatoDato(completarSecuencia(0, 9), 9)
				// } else {

				// 	fila += formatoDato("", 26)

				// 	fila += "\n" // siguiente persona...
				// 	filas += fila
				// 	secuencia++
				// 	if suspencionTemporalContrato {
				// 		diasNovedad = diasSuspencionContrato
				// 		crearFilaNovedad(idPersona, idPreliquidacion, "licencia_norem", value)
				// 	} else if licenciaNoRemunerada {
				// 		diasNovedad = diasLicenciaNoRem
				// 		crearFilaNovedad(idPersona, idPreliquidacion, "licencia_norem", value)
				// 	} else if comisionServicios {
				// 		diasNovedad = diasComisionServicios
				// 		crearFilaNovedad(idPersona, idPreliquidacion, "comision_norem", value)
				// 	} else if incapacidadGeneral {
				// 		diasNovedad = diasIncapacidad
				// 		crearFilaNovedad(idPersona, idPreliquidacion, "incapacidad_general", value)
				// 	} else if licenciaMaternidad {
				// 		diasNovedad = diasLicenciaMaternidad
				// 		crearFilaNovedad(idPersona, idPreliquidacion, "licencia_maternidad", value)
				// 	} else if licenciaPaternidad {
				// 		diasNovedad = diasLicenciaMaternidad
				// 		crearFilaNovedad(idPersona, idPreliquidacion, "licencia_paternidad", value)
				// 	} else if vacaciones {
				// 		diasNovedad = diasVacaciones
				// 		crearFilaNovedad(idPersona, idPreliquidacion, "vacaciones", value)
				// 	} else if licenciaRemunerada {
				// 		diasNovedad = diasLicenciaRem
				// 		crearFilaNovedad(idPersona, idPreliquidacion, "licencia_rem", value)
				// 	}
				// }
			}
			log.Println("Finalizó de generar la planilla")
			log.Println("Tiempo en generar la planilla: ", time.Since(start))
			respuestaJSON := make(map[string]interface{})
			respuestaJSON["informacion"] = filas
			c.Data["json"] = respuestaJSON
			c.Data["json"] = filasPlanilla
		} else {
			log.Println("Fallo la generación de la planilla")
			c.Data["json"] = err.Error()
		}

	} else {
		log.Println("Fallo la generación de la planilla")
		c.Data["json"] = map[string]string{
			"error": err.Error(),
		}
	}
	c.ServeJSON()
}

// crearFilaNovedad crea las filas de acuerdo a las novedades de una persona
func crearFilaNovedad(idPersona, idPreliquidacion, novedad string, persona models.InformacionPersonaNatural) {
	var horasLaboradasNovedad = "0"
	var pagoSeguridadSocial models.PagoSeguridadSocial

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

	ibcNovedad := calcularIbc(idPersona, persona.Id, novedad, idPreliquidacion)
	// beego.Info("ibcNovedad:", completarSecuencia(ibcNovedad, 9), 9)
	filaAux += formatoDato(completarSecuencia(ibcNovedad, 9), 9) //IBC pensión
	filaAux += formatoDato(completarSecuencia(ibcNovedad, 9), 9) //IBC salud
	filaAux += formatoDato(completarSecuencia(ibcNovedad, 9), 9) //IBC ARL
	filaAux += formatoDato(completarSecuencia(ibcNovedad, 9), 9) //IBC CCF

	if novedad == "licencia_norem" || (novedad == "comision_norem") {
		pagoSeguridadSocial = *calcularValoresAux(false, idPersona, ibcNovedad)
	} else {
		pagoSeguridadSocial = *calcularValoresAux(true, idPersona, ibcNovedad)
	}
	// beego.Info("pagoSeguirdadSocial:", pagoSeguridadSocial)

	filaAux += formatoDato(tarifaPension, 7) //Tarifa de aportes pensiones

	filaAux += formatoDato(completarSecuencia(int(pagoSeguridadSocial.PensionTotal), 9), 9) // Cotización obligatoria de pensión
	filaAux += formatoDato(completarSecuencia(0, 9), 9)                                     // Aporte voluntario del afiliado al fondo de pensiones obligatorias
	filaAux += formatoDato(completarSecuencia(0, 9), 9)                                     // Aporte voluntario del aportante al fondo de pensiones obligatorias
	filaAux += formatoDato(completarSecuencia(int(pagoSeguridadSocial.PensionTotal), 9), 9) // Total cotizaicón sistema general de pensiones

	filaAux += formatoDato(completarSecuencia(0, 9), 9) // Aportes a fondo de solidaridad pensional subcuenta de solidaridad
	filaAux += formatoDato(completarSecuencia(0, 9), 9) // Aportes al fondo de solidaridad pensional subcuenta de subsistencia
	filaAux += formatoDato(completarSecuencia(0, 9), 9) // Valor no retenido por aportes voluntarios

	filaAux += formatoDato(tarifaSalud, 7)                                                // Tarifa de aportes salud
	filaAux += formatoDato(completarSecuencia(int(pagoSeguridadSocial.SaludTotal), 9), 9) // Cotización obligatoria a salud

	filaAux += valorUpc                                 //Valor UPC Adicional
	filaAux += formatoDato("", 15)                      //Nº de autorización de la incapacidad por enfermedad general
	filaAux += formatoDato(completarSecuencia(0, 9), 9) //Valor de la incapacidad por enfermedad general
	filaAux += formatoDato("", 15)                      //Nº de autorización de la licencia de maternidad o paternidad
	filaAux += formatoDato(completarSecuencia(0, 9), 9) //Valor de la licencia de maternidad

	filaAux += formatoDato(tarifaArl, 9) //Tarifa de aportes a Riegos Laborales

	filaAux += formatoDato("        1", 9)                                         // Centro de trabajo CT
	filaAux += formatoDato(completarSecuencia(int(pagoSeguridadSocial.Arl), 9), 9) // Cotización obligatoria a sistema de riesgos laborales

	filaAux += formatoDato(tarifaCaja, 7)                                           // Tarifa de aportes CCF
	filaAux += formatoDato(completarSecuencia(int(pagoSeguridadSocial.Caja), 9), 9) // Valor aporte CCF

	filaAux += formatoDato(completarSecuencia(0, 7), 7) // Tarifa de aportes SENA
	filaAux += formatoDato(completarSecuencia(0, 9), 9) // Valor Aportes SENA

	filaAux += formatoDato(tarifaIcbf, 7)                                           //Tarifa de aportes ICBF
	filaAux += formatoDato(completarSecuencia(0, int(pagoSeguridadSocial.Icbf)), 9) // Valor aporte ICBF

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

// calcularValoresAux cálcula los valores de la fila auxiliar
func calcularValoresAux(pagoCompleto bool, idPersona string, ibc int) *models.PagoSeguridadSocial {

	var predicados []models.Predicado
	var pago *models.PagoSeguridadSocial
	ibcAux := strconv.Itoa(ibc)
	predicados = append(predicados, models.Predicado{Nombre: "ibc(" + idPersona + "," + ibcAux + ", salud)."})
	predicados = append(predicados, models.Predicado{Nombre: "ibc(" + idPersona + "," + ibcAux + ", riesgos)."})
	predicados = append(predicados, models.Predicado{Nombre: "ibc(" + idPersona + "," + ibcAux + ", apf)."})

	reglas := CrearReglas(pagoCompleto) + FormatoReglas(predicados)

	idProveedor := golog.GetInt64(reglas, "v_total_salud(X,T).", "X")
	saludTotal := golog.GetInt64(reglas, "v_total_salud(X,T).", "T")
	pensionTotal := golog.GetInt64(reglas, "v_total_pen(X,T).", "T")
	arl := golog.GetInt64(reglas, "v_arl(I,Y).", "Y")
	caja := golog.GetInt64(reglas, "v_caja(I,Y).", "Y")
	icbf := golog.GetInt64(reglas, "v_icbf(I,Y).", "Y")

	for index := 0; index < len(idProveedor); index++ {

		pago = &models.PagoSeguridadSocial{
			NombrePersona:           "",
			IdProveedor:             idProveedor[index],
			SaludUd:                 0,
			SaludTotal:              saludTotal[index],
			PensionUd:               0,
			PensionTotal:            pensionTotal[index],
			FondoSolidaridad:        0,
			Caja:                    caja[index],
			Icbf:                    icbf[index],
			IdPreliquidacion:        0,
			IdDetallePreliquidacion: 0,
			Arl:                     arl[index]}
	}
	return pago
}

func calcularIbc(idPersona, documentoPersona, novedad, idPreliquidacion string) int {
	idInt, err := strconv.Atoi(idPersona)
	if err != nil {
		ImprimirError("error en calcularIbc()", err)
	}

	documentoInt, err := strconv.Atoi(documentoPersona)
	if err != nil {
		ImprimirError("error en calcularIbc()", err)
	}

	var respuestaAPI interface{}

	// "seg_social(exterior_familia,2018,3,7,2018,4,7). concepto(420,seguridad_social, seguridad_social, exterior_familia, 0, 2018)."

	// seg_social(nombre de regla, año inicio, mes inicio, dia inicio, año fin, mes fin, dia fin).
	novedadSeguridadSocial := "seg_social(" + novedad + ", " + strconv.Itoa(fechaInicioNovedad.Year()) + ", " +
		strconv.Itoa(int(fechaInicioNovedad.Month())) + ", " + strconv.Itoa(fechaInicioNovedad.Day()) + ", " +
		strconv.Itoa(fechaFinNovedad.Year()) + ", " + strconv.Itoa(int(fechaFinNovedad.Month())) + ", " +
		strconv.Itoa(fechaInicioNovedad.Day()) + ")."
	// concepto(id proveedor, seguridad_social, seguridad_social, nombre de regla, 0, vigencia)
	concepto := "concepto(" + idPersona + ", seguridad_social, seguridad_social, " + novedad + ", " +
		strconv.Itoa(fechaInicioNovedad.Year()) + ")."

	infoRequest := make(map[string]interface{})
	infoRequest["IdPersona"] = idInt
	infoRequest["NumDocumento"] = documentoInt
	infoRequest["NombreNomina"] = tipoPreliquidacion
	infoRequest["Novedad"] = novedadSeguridadSocial + " " + concepto
	infoRequest["Ano"] = anioPeriodo
	infoRequest["Mes"] = mesPeriodo

	err = sendJson("http://"+beego.AppConfig.String("titanMidService")+"/preliquidacion/get_ibc_novedad", "POST", &respuestaAPI, infoRequest)
	if err != nil {
		ImprimirError("error en calcularIbc()", err)
	}

	return int(respuestaAPI.(float64))
}

// buscarUpcAsociada busca todas las upc asociadas a esa persona
// devuelve un mapa cuya llave es la persona asociada a la upc y el valor del mapa es el valor de la(s) upc asocaida a la persona
func buscarUpcAsociada() (map[string]int, error) {
	var beneficiariosAdicionales []models.UpcAdicional
	mapaBeneficiarios := make(map[string]int)

	err := getJson("http://"+beego.AppConfig.String("segSocialService")+
		"/upc_adicional"+
		"?limit=-1", &beneficiariosAdicionales)
	if err != nil {
		ImprimirError("error en buscarUpcAsociada()", err)
		return nil, err
	}

	if beneficiariosAdicionales[0].Id > 0 { // Esto es para verificar que el arreglo que devuelve el servicio tiene al menos un elemento con valor
		for _, beneficiario := range beneficiariosAdicionales {
			mapaBeneficiarios[strconv.Itoa(beneficiario.PersonaAsociada)] += int(beneficiario.TipoUpc.Valor)
		}
	}
	return mapaBeneficiarios, nil
}

// traerDiasCotizados función para trer los días cotizados, tipoPago debe tener el mismo nombre que un concepto_nomina en titán
// duelve:
// diasLiquidados: días coritzados del tipo de pago
func traerDiasCotizados(idPersona, idPreliquidacion, tipoPago string) int {

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

	return dias
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

// traerDiasCotizadosEmpleador trae los días cotizados correspondientes al pago del empleador
// devuelve
// diasCotizados: días cotizados del tipo de pago
// valorTipoPago: valor del tipo de pago por parte de la universidad
// valorTotalPago: valor total del tipo de pago (sumando lo del empleado y la universidad)
func traerDiasCotizadosEmpleador(idPreliquidacion, periodoPago, tipoPago string) (map[string]map[string]int, error) {
	var (
		pago           []models.Pago
		conceptoNomina []models.ConceptoNomina
	)
	infoPago := make(map[string]map[string]int)

	err := getJson("http://"+beego.AppConfig.String("titanServicio")+
		"/concepto_nomina"+
		"?limit=1"+
		"&query=NaturalezaConcepto.Nombre:seguridad_social"+
		",NombreConcepto:"+tipoPago, &conceptoNomina)
	if err != nil {
		ImprimirError("error en traerDiasCotizadosEmpleador()", err)
	}

	if conceptoNomina == nil {
		err := getJson("http://"+beego.AppConfig.String("titanServicio")+
			"/concepto_nomina"+
			"?limit=1"+
			"&query=NombreConcepto:"+tipoPago, &conceptoNomina)
		if err != nil {
			ImprimirError("error en traerDiasCotizadosEmpleador()", err)
		}
	}

	err = getJson("http://"+beego.AppConfig.String("segSocialService")+
		"/pago"+
		"?limit=1"+
		"&query=DetalleLiquidacion:"+idPreliquidacion+
		",PeriodoPago.Id:"+periodoPago+
		",TipoPago:"+strconv.Itoa(conceptoNomina[0].Id), &pago)
	if err != nil {
		ImprimirError("error en traerDiasCotizadosEmpleador()", err)
		return nil, err
	}

	for _, value := range detallePreliquidacion {
		auxiliarMap := make(map[string]int) // map con la información interna
		idProveedor := strconv.Itoa(value.Persona)
		valorTipoPagoTemp := GetPagoEmpleado(idProveedor, idPreliquidacion, tipoPago)

		auxiliarMap["valor"] += valorTipoPagoTemp + int(pago[0].Valor)
		auxiliarMap["dias"] = 0

		if pago[0].EntidadPago != 0 {
			auxiliarMap["dias"] = traerDiasCotizados(idProveedor, idPreliquidacion, "salud")

		}
		infoPago[idProveedor] = auxiliarMap
	}

	return infoPago, nil
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

	fechaIngresoTemp := contratosElaboradosPeriodo[cedulaPersona]

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
				fechaInicioNovedad = time.Date(auxFechaInicio.Year(), time.Month(mesPeriodo), 1, 0, 0, 0, 0, time.UTC)
				fechaInicioTemp = time.Date(auxFechaInicio.Year(), time.Month(mesPeriodo), 1, 0, 0, 0, 0, time.UTC).Format(formatoFecha)
			} else {
				fechaInicioNovedad = conceptoNominaPersona[0].FechaDesde
				fechaInicioTemp = conceptoNominaPersona[0].FechaDesde.Format(formatoFecha)
			}

			auxFechaFin := conceptoNominaPersona[0].FechaHasta
			if int(auxFechaFin.Month()) > mesPeriodo {
				fechaFinNovedad = time.Date(auxFechaFin.Year(), time.Month(mesPeriodo), 30, 0, 0, 0, 0, time.UTC)
				fechaFinTemp = time.Date(auxFechaFin.Year(), time.Month(mesPeriodo), 30, 0, 0, 0, 0, time.UTC).Format(formatoFecha)
			} else {
				fechaFinNovedad = conceptoNominaPersona[0].FechaHasta
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
			switch value.Concepto.NombreConcepto {
			case "licencia_maternidad":
				licenciaMaternidad = true
			case "licencia_paternidad":
				licenciaPaternidad = true
			}
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

// func revisarIngreso(idPreliquidacion, cedulaPersona string) (fechaMenor time.Time) {
// 	var preliquidacion models.Preliquidacion

// 	err := getJson("http://"+beego.AppConfig.String("titanServicio")+
// 		"/preliquidacion/"+idPreliquidacion, &preliquidacion)
// 	if err != nil {
// 		ImprimirError("error en revisarIngreso()", err)
// 	}

// 	// beego.Info(contratosElaboradosPeriodo)
// 	var actaInicio map[string]map[string]interface{}

// 	for key, value := range contratosElaboradosPeriodo {

// 		// numeroContrato := value["numero_contrato"].(string)
// 		vigenciaContrato := value["vigencia"].(string)

// 		err := getJsonWSO2("http://"+beego.AppConfig.String("argoWso2Service")+
// 			"/acta_inicio_elaborado/"+key+"/"+vigenciaContrato, &actaInicio)
// 		// beego.Info("http://" + beego.AppConfig.String("argoWso2Service") +
// 		// 	"/acta_inicio_elaborado/" + numeroContrato + "/" + vigenciaContrato)
// 		if err != nil {
// 			ImprimirError("error en revisarIngreso()", err)
// 		}

// 		if len(actaInicio["actaInicio"]) != 0 {

// 			t, err := time.Parse(formatoFecha, actaInicio["actaInicio"]["fechaInicio"].(string))
// 			if err != nil {
// 				ImprimirError("error en revisarIngreso()", err)
// 			}

// 			if fechaMenor.IsZero() {
// 				fechaMenor = t
// 			} else {
// 				if t.Year() <= fechaMenor.Year() {
// 					if t.Month() <= fechaMenor.Month() {
// 						if t.Day() <= fechaMenor.Day() {
// 							fechaMenor = t
// 						} else if (t.Year() <= fechaMenor.Year()) && (t.Month() <= fechaMenor.Month()) && (t.Day() >= fechaMenor.Day()) {
// 							fechaMenor = t
// 						}
// 					} else if (t.Year() <= fechaMenor.Year()) && (t.Month() >= fechaMenor.Month()) {
// 						fechaMenor = t
// 					}
// 				}
// 			}
// 		}
// 	}
// 	return
// }

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
	var acumuladorPago float64
	switch tipoPago {
	case "salud_ud":
		err := getJson("http://"+beego.AppConfig.String("titanServicio")+
			"/detalle_preliquidacion"+
			"?limit=0"+
			"&query=Persona:"+idPersona+
			",Preliquidacion.Id:"+idPreliquidacion+
			",Concepto.NombreConcepto:salud", &detallePreliquidacion)
		if err != nil {
			ImprimirError("error en traerDiasCotizadosEmpleador()", err)
		}

	case "pension_ud":
		err := getJson("http://"+beego.AppConfig.String("titanServicio")+
			"/detalle_preliquidacion"+
			"?limit=0"+
			"&query=Persona:"+idPersona+
			",Preliquidacion.Id:"+idPreliquidacion+
			",Concepto.NombreConcepto:salud", &detallePreliquidacion)
		if err != nil {
			ImprimirError("error en traerDiasCotizadosEmpleador()", err)
		}

	}
	for _, value := range detallePreliquidacion {
		acumuladorPago += value.ValorCalculado
	}
	valorPago = AproximarPesoSuperior(acumuladorPago, 100)
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
