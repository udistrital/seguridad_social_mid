//Comienza la corrección...

package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/udistrital/seguridad_social_mid/golog"
	"github.com/udistrital/seguridad_social_mid/models"
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
	detallePreliquidacion []models.DetallePreliquidacion   // detalle de toda la preliquidación
	informacionUpcs       map[string][]models.UpcAdicional // información de todos los beneficiarios adicionales
	contratistas          = false

	formatoFecha = "2006-01-02"
	fila         string
	filaAux      string

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
		salariosBase[strconv.FormatInt(value.Persona, 10)] += int(value.ValorCalculado)
	}
	return salariosBase, nil
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

			for i := range detallePreliquidacion {
				personas = append(personas, fmt.Sprint(detallePreliquidacion[i].Persona))
			}
			mapPersonas, err := GetInfoPersonas(detallePreliquidacion)
			if err != nil {
				c.Data["json"] = map[string]string{"error": err.Error()}
				c.ServeJSON()
				return
			}

			for key, value := range mapPersonas {
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
					ColombianoExterior:              models.Columna{Valor: "", Longitud: 1},
					CodigoDepartamento:              models.Columna{Valor: "11", Longitud: 2},
					CodigoMunicipio:                 models.Columna{Valor: "001", Longitud: 3},
					PrimerApellido:                  models.Columna{Valor: value.PrimerApellido, Longitud: 20},
					SegundoApellido:                 models.Columna{Valor: value.SegundoApellido, Longitud: 30},
					PrimerNombre:                    models.Columna{Valor: value.PrimerNombre, Longitud: 20},
					SegundoNombre:                   models.Columna{Valor: value.SegundoNombre, Longitud: 30},
					NovIng:                          models.Columna{Valor: "", Longitud: 1},
					NovRet:                          models.Columna{Valor: "", Longitud: 1},
					NovTde:                          models.Columna{Valor: "", Longitud: 1},
					NovTae:                          models.Columna{Valor: "", Longitud: 1},
					NovTdp:                          models.Columna{Valor: "", Longitud: 1},
					NovTap:                          models.Columna{Valor: "", Longitud: 1},
					NovVsp:                          models.Columna{Valor: "", Longitud: 1},
					NovCorrecciones:                 models.Columna{Valor: "", Longitud: 1},
					NovVst:                          models.Columna{Valor: "", Longitud: 1},
					NovSln:                          models.Columna{Valor: "", Longitud: 1},
					NovIge:                          models.Columna{Valor: "", Longitud: 1},
					NovLma:                          models.Columna{Valor: "", Longitud: 1},
					NovVac:                          models.Columna{Valor: "", Longitud: 1},
					NovAvp:                          models.Columna{Valor: "", Longitud: 1},
					NovVct:                          models.Columna{Valor: "", Longitud: 1},
					NavIrl:                          models.Columna{Valor: 0, Longitud: 2},
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
					IbcPension:                      models.Columna{Valor: ingresosCotizacion[key], Longitud: 9},
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
					FechaIngreso:                    models.Columna{Valor: "", Longitud: 10},
					FechaRetiro:                     models.Columna{Valor: "", Longitud: 10},
					FechaInicioVsp:                  models.Columna{Valor: "", Longitud: 10},
					FechaInicioSuspencion:           models.Columna{Valor: "", Longitud: 10},
					FechaFinSuspencion:              models.Columna{Valor: "", Longitud: 10},
					FechaInicioIge:                  models.Columna{Valor: "", Longitud: 10},
					FechaFinIge:                     models.Columna{Valor: "", Longitud: 10},
					FechaInicioLma:                  models.Columna{Valor: "", Longitud: 10},
					FechaFinLma:                     models.Columna{Valor: "", Longitud: 10},
					FechaInicioVac:                  models.Columna{Valor: "", Longitud: 10},
					FechaFinVac:                     models.Columna{Valor: "", Longitud: 10},
					FechaInicioVct:                  models.Columna{Valor: "", Longitud: 10},
					FechaFinVct:                     models.Columna{Valor: "", Longitud: 10},
					FechaInicioIrl:                  models.Columna{Valor: "", Longitud: 10},
					FechaFinIrl:                     models.Columna{Valor: "", Longitud: 10},
					IbcOtrosParaFiscales:            models.Columna{Valor: ingresosCotizacion[key], Longitud: 9},
					HorasLaboradas:                  models.Columna{Valor: diasLaborados * 8, Longitud: 3},
					EspacioBlanco:                   models.Columna{Valor: "", Longitud: 26},
				}

				filasPlanilla = append(filasPlanilla, filaPlanilla)

				establecerNovedadesExterior(idPersona, idPreliquidacion, &filaPlanilla)

				establecerNovedades(idPersona, idPreliquidacion, cedulaPersona)

				if suspencionTemporalContrato {
					diasNovedad = diasSuspencionContrato
					filasPlanilla = append(filasPlanilla, crearFilaNovedad(filaPlanilla, idPersona, idPreliquidacion, "licencia_norem", value))
				} else if licenciaNoRemunerada {
					diasNovedad = diasLicenciaNoRem
					filasPlanilla = append(filasPlanilla, crearFilaNovedad(filaPlanilla, idPersona, idPreliquidacion, "licencia_norem", value))
				} else if comisionServicios {
					diasNovedad = diasComisionServicios
					filasPlanilla = append(filasPlanilla, crearFilaNovedad(filaPlanilla, idPersona, idPreliquidacion, "comision_norem", value))
				} else if incapacidadGeneral {
					diasNovedad = diasIncapacidad
					filasPlanilla = append(filasPlanilla, crearFilaNovedad(filaPlanilla, idPersona, idPreliquidacion, "incapacidad_general", value))
				} else if licenciaMaternidad {
					diasNovedad = diasLicenciaMaternidad
					filasPlanilla = append(filasPlanilla, crearFilaNovedad(filaPlanilla, idPersona, idPreliquidacion, "licencia_maternidad", value))
				} else if licenciaPaternidad {
					diasNovedad = diasLicenciaMaternidad
					filasPlanilla = append(filasPlanilla, crearFilaNovedad(filaPlanilla, idPersona, idPreliquidacion, "licencia_paternidad", value))
				} else if vacaciones {
					diasNovedad = diasVacaciones
					filasPlanilla = append(filasPlanilla, crearFilaNovedad(filaPlanilla, idPersona, idPreliquidacion, "vacaciones", value))
				} else if licenciaRemunerada {
					diasNovedad = diasLicenciaRem
					filasPlanilla = append(filasPlanilla, crearFilaNovedad(filaPlanilla, idPersona, idPreliquidacion, "licencia_rem", value))
				}

			}
			log.Println("Finalizó de generar la planilla")
			log.Println("Tiempo en generar la planilla: ", time.Since(start))
			filasUpc, err := getFilasUpc(mapPersonas)
			if err != nil {
				log.Println("Falló al traer las filas de UPC")
				c.Data["json"] = map[string]string{"error": err.Error()}
			}
			filasPlanilla = append(filasPlanilla, filasUpc...)
			c.Data["json"] = filasPlanilla
		} else {
			log.Println("Falló la generación de la planilla")
			c.Data["json"] = map[string]string{"error": err.Error()}
		}

	} else {
		log.Println("Falló la generación de la planilla")
		c.Data["json"] = map[string]string{"error": err.Error()}
	}
	c.ServeJSON()
}

//getFilasUpc trae todas las fila de las upc de las personas correspondientes a ese periodo
func getFilasUpc(mapPersonas map[string]models.InformacionPersonaNatural) ([]models.PlanillaTipoE, error) {
	var filasUpc []models.PlanillaTipoE
	parametros, err := GetParametroEstandar()
	if err != nil {
		return filasUpc, err
	}
	for key, value := range mapPersonas {
		if informacionUpcs[key] != nil {
			for _, upc := range informacionUpcs[key] {
				filaPlanilla := models.PlanillaTipoE{
					TipoRegistro:                    models.Columna{Valor: "02", Longitud: 2},
					TipoDocumento:                   models.Columna{Valor: parametros[upc.ParametroEstandar], Longitud: 2},
					NumeroIdentificacion:            models.Columna{Valor: upc.NumDocumento, Longitud: 16},
					TipoCotizante:                   models.Columna{Valor: 40, Longitud: 2},
					SubTipoCotizante:                models.Columna{Valor: 1, Longitud: 2},
					ExtranjeroNoPension:             models.Columna{Valor: "", Longitud: 1},
					ColombianoExterior:              models.Columna{Valor: "", Longitud: 1},
					CodigoDepartamento:              models.Columna{Valor: "", Longitud: 2},
					CodigoMunicipio:                 models.Columna{Valor: "", Longitud: 3},
					PrimerApellido:                  models.Columna{Valor: upc.PrimerApellido, Longitud: 20},
					SegundoApellido:                 models.Columna{Valor: upc.SegundoApellido, Longitud: 30},
					PrimerNombre:                    models.Columna{Valor: upc.PrimerNombre, Longitud: 20},
					SegundoNombre:                   models.Columna{Valor: upc.SegundoNombre, Longitud: 30},
					NovIng:                          models.Columna{Valor: "", Longitud: 1},
					NovRet:                          models.Columna{Valor: "", Longitud: 1},
					NovTde:                          models.Columna{Valor: "", Longitud: 1},
					NovTae:                          models.Columna{Valor: "", Longitud: 1},
					NovTdp:                          models.Columna{Valor: "", Longitud: 1},
					NovTap:                          models.Columna{Valor: "", Longitud: 1},
					NovVsp:                          models.Columna{Valor: "", Longitud: 1},
					NovCorrecciones:                 models.Columna{Valor: "", Longitud: 1},
					NovVst:                          models.Columna{Valor: "", Longitud: 1},
					NovSln:                          models.Columna{Valor: "", Longitud: 1},
					NovIge:                          models.Columna{Valor: "", Longitud: 1},
					NovLma:                          models.Columna{Valor: "", Longitud: 1},
					NovVac:                          models.Columna{Valor: "", Longitud: 1},
					NovAvp:                          models.Columna{Valor: "", Longitud: 1},
					NovVct:                          models.Columna{Valor: "", Longitud: 1},
					NavIrl:                          models.Columna{Valor: 0, Longitud: 2},
					CodigoFondoPension:              models.Columna{Valor: "", Longitud: 6},
					TrasladoPension:                 models.Columna{Valor: "", Longitud: 6},
					CodigoEps:                       models.Columna{Valor: "", Longitud: 6},
					TrasladoEps:                     models.Columna{Valor: "", Longitud: 6},
					CodigoCCF:                       models.Columna{Valor: "", Longitud: 6},
					DiasLaborados:                   models.Columna{Valor: 0, Longitud: 2},
					DiasPension:                     models.Columna{Valor: 0, Longitud: 2},
					DiasSalud:                       models.Columna{Valor: 0, Longitud: 2},
					DiasArl:                         models.Columna{Valor: 0, Longitud: 2},
					DiasCaja:                        models.Columna{Valor: 0, Longitud: 2},
					SalarioBase:                     models.Columna{Valor: 0, Longitud: 9},
					SalarioIntegral:                 models.Columna{Valor: "", Longitud: 1},
					IbcPension:                      models.Columna{Valor: 0, Longitud: 9},
					IbcSalud:                        models.Columna{Valor: 0, Longitud: 9},
					IbcArl:                          models.Columna{Valor: 0, Longitud: 9},
					IbcCcf:                          models.Columna{Valor: 0, Longitud: 9},
					TarifaPension:                   models.Columna{Valor: 0, Longitud: 7},
					PagoPension:                     models.Columna{Valor: "", Longitud: 9},
					AportePension:                   models.Columna{Valor: "", Longitud: 9},
					TotalPension:                    models.Columna{Valor: "", Longitud: 9},
					FondoSolidaridad:                models.Columna{Valor: "", Longitud: 9},
					FondoSubsistencia:               models.Columna{Valor: "", Longitud: 9},
					NoRetenidoAportesVolunarios:     models.Columna{Valor: 0, Longitud: 9},
					TarifaSalud:                     models.Columna{Valor: "", Longitud: 7},
					PagoSalud:                       models.Columna{Valor: "", Longitud: 9},
					ValorUpc:                        models.Columna{Valor: int(upc.TipoUpc.Valor), Longitud: 9},
					AutorizacionEnfermedadGeneral:   models.Columna{Valor: "", Longitud: 15},
					ValorIncapacidadGeneral:         models.Columna{Valor: 0, Longitud: 15},
					AutotizacionLicenciaMarternidad: models.Columna{Valor: "", Longitud: 15},
					ValorLicenciaMaternidad:         models.Columna{Valor: 0, Longitud: 15},
					TarifaArl:                       models.Columna{Valor: "", Longitud: 9},
					CentroTrabajo:                   models.Columna{Valor: 0, Longitud: 9},
					PagoArl:                         models.Columna{Valor: 0, Longitud: 9},
					TarifaCaja:                      models.Columna{Valor: 0, Longitud: 7},
					PagoCaja:                        models.Columna{Valor: 0, Longitud: 9},
					TarifaSena:                      models.Columna{Valor: 0, Longitud: 7},
					PagoSena:                        models.Columna{Valor: 0, Longitud: 9},
					TarifaIcbf:                      models.Columna{Valor: 0, Longitud: 7},
					PagoIcbf:                        models.Columna{Valor: 0, Longitud: 9},
					TarifaEsap:                      models.Columna{Valor: 0, Longitud: 9},
					PagoEsap:                        models.Columna{Valor: 0, Longitud: 9},
					TarifaMen:                       models.Columna{Valor: 0, Longitud: 9},
					PagoMen:                         models.Columna{Valor: 0, Longitud: 9},
					TipoDocumentoCotizantePrincipal: models.Columna{Valor: "CC", Longitud: 2},
					DocumentoCotizantePrincipal:     models.Columna{Valor: value.Id, Longitud: 16},
					ExoneradoPagoSalud:              models.Columna{Valor: "", Longitud: 1},
					CodigoArl:                       models.Columna{Valor: "", Longitud: 6},
					ClaseRiesgo:                     models.Columna{Valor: "", Longitud: 1},
					IndicadorTarifaEspecialPension:  models.Columna{Valor: "", Longitud: 1},
					FechaIngreso:                    models.Columna{Valor: "", Longitud: 10},
					FechaRetiro:                     models.Columna{Valor: "", Longitud: 10},
					FechaInicioVsp:                  models.Columna{Valor: "", Longitud: 10},
					FechaInicioSuspencion:           models.Columna{Valor: "", Longitud: 10},
					FechaFinSuspencion:              models.Columna{Valor: "", Longitud: 10},
					FechaInicioIge:                  models.Columna{Valor: "", Longitud: 10},
					FechaFinIge:                     models.Columna{Valor: "", Longitud: 10},
					FechaInicioLma:                  models.Columna{Valor: "", Longitud: 10},
					FechaFinLma:                     models.Columna{Valor: "", Longitud: 10},
					FechaInicioVac:                  models.Columna{Valor: "", Longitud: 10},
					FechaFinVac:                     models.Columna{Valor: "", Longitud: 10},
					FechaInicioVct:                  models.Columna{Valor: "", Longitud: 10},
					FechaFinVct:                     models.Columna{Valor: "", Longitud: 10},
					FechaInicioIrl:                  models.Columna{Valor: "", Longitud: 10},
					FechaFinIrl:                     models.Columna{Valor: "", Longitud: 10},
					IbcOtrosParaFiscales:            models.Columna{Valor: 0, Longitud: 9},
					HorasLaboradas:                  models.Columna{Valor: 0, Longitud: 3},
					EspacioBlanco:                   models.Columna{Valor: "", Longitud: 26},
				}
				filasUpc = append(filasUpc, filaPlanilla)
			}

		}
	}

	return filasUpc, err
}

// crearFilaNovedad crea las filas de acuerdo a las novedades de una persona
func crearFilaNovedad(filaPersona models.PlanillaTipoE, idPersona, idPreliquidacion, novedad string, persona models.InformacionPersonaNatural) (personaNovedad models.PlanillaTipoE) {
	var pagoSeguridadSocial models.PagoSeguridadSocial

	filaPersona.ColombianoExterior.Valor = marcarNovedad(exterior, "X")

	filaPersona.NovIng.Valor = marcarNovedad(ingreso, "X")
	filaPersona.NovRet.Valor = marcarNovedad(retiro, "X")
	filaPersona.NovTde.Valor = marcarNovedad(trasladoDesdeEps, "X")
	filaPersona.NovTae.Valor = marcarNovedad(trasladoDesdeEps, "X")
	filaPersona.NovTdp.Valor = marcarNovedad(trasladoDesdePensiones, "X")
	filaPersona.NovTap.Valor = marcarNovedad(trasladoDesdeEps, "X")
	filaPersona.NovVsp.Valor = marcarNovedad(variacionPermanteSalario, "X")
	filaPersona.NovCorrecciones.Valor = marcarNovedad(corecciones, "A")
	filaPersona.NovVst.Valor = marcarNovedad(variacionTransitoriaSalario, "X")

	//SLN: Suspención temporal del contrato o licencia no remunerada o comisión de servicios
	if suspencionTemporalContrato || (licenciaNoRemunerada) {
		filaPersona.NovSln.Valor = "X"
	}
	if comisionServicios {
		filaPersona.NovSln.Valor = "C"
	}

	filaPersona.NovIge.Valor = marcarNovedad(incapacidadGeneral, "X")
	filaPersona.NovLma.Valor = marcarNovedad(licenciaMaternidad, "X")

	//VAC - LR: Vacaciones, licencia remunerada
	if vacaciones {
		filaPersona.NovVac.Valor = "X"
	}
	if licenciaRemunerada {
		filaPersona.NovVac.Valor = "L"
	}

	filaPersona.NovAvp.Valor = marcarNovedad(aporteVoluntario, "X")
	filaPersona.NovVct.Valor = marcarNovedad(variacionCentroTrabajo, "X")

	//Código de la admnistradora de pensiones a la cual se traslada el afiliado
	// Si hay un translado, debe aparecer el nuevo código, de lo contrario será un campo vació
	if trasladoAPensiones {
		filaPersona.TrasladoPension.Valor = "230301"
	}

	//Código EPS o EOC a la cual pertenece el afiliado
	if filaPersona.CodigoEps.Valor == "      " {
		filaPersona.CodigoEps.Valor = "MIN001"
	}

	if trasladoAEps {
		filaPersona.TrasladoEps.Valor = "EPS012"
	}

	if suspencionTemporalContrato {
		escribirDiasCotizadosAux(&filaPersona, diasSuspencionContrato)
	} else if licenciaNoRemunerada {
		escribirDiasCotizadosAux(&filaPersona, diasLicenciaNoRem)
	} else if comisionServicios {
		escribirDiasCotizadosAux(&filaPersona, diasComisionServicios)
	} else if licenciaMaternidad {
		escribirDiasCotizadosAux(&filaPersona, diasLicenciaMaternidad)
	} else if vacaciones {
		escribirDiasCotizadosAux(&filaPersona, diasVacaciones)
	} else if licenciaRemunerada {
		escribirDiasCotizadosAux(&filaPersona, diasLicenciaRem)
	}

	ibcNovedad := calcularIbc(idPersona, persona.Id, novedad, idPreliquidacion)

	filaPersona.IbcPension.Valor = ibcNovedad
	filaPersona.IbcSalud.Valor = ibcNovedad
	filaPersona.IbcArl.Valor = ibcNovedad
	filaPersona.IbcCcf.Valor = ibcNovedad

	if novedad == "licencia_norem" || (novedad == "comision_norem") {
		pagoSeguridadSocial = *calcularValoresAux(false, idPersona, ibcNovedad)
	} else {
		pagoSeguridadSocial = *calcularValoresAux(true, idPersona, ibcNovedad)
	}

	filaPersona.TarifaPension.Valor = tarifaPension
	filaPersona.PagoPension.Valor = int(pagoSeguridadSocial.PensionTotal)
	filaPersona.TotalPension.Valor = int(pagoSeguridadSocial.PensionTotal)
	filaPersona.FondoSolidaridad.Valor = 0
	filaPersona.FondoSubsistencia.Valor = 0
	filaPersona.NoRetenidoAportesVolunarios.Valor = 0
	filaPersona.TarifaSalud.Valor = tarifaSalud
	filaPersona.PagoSalud.Valor = int(pagoSeguridadSocial.SaludTotal)
	filaPersona.TarifaArl.Valor = tarifaArl
	filaPersona.PagoArl.Valor = int(pagoSeguridadSocial.Arl)
	filaPersona.TarifaCaja.Valor = tarifaCaja
	filaPersona.PagoCaja.Valor = int(pagoSeguridadSocial.Caja)
	filaPersona.TarifaIcbf.Valor = tarifaIcbf
	filaPersona.PagoIcbf.Valor = int(pagoSeguridadSocial.Icbf)
	filaPersona.FechaIngreso.Valor = fechaIngreso
	filaPersona.FechaRetiro.Valor = fechaRetiro
	filaPersona.FechaInicioVsp.Valor = fechaInicioVsp
	filaPersona.FechaInicioSuspencion.Valor = fechaInicioSuspencion
	filaPersona.FechaInicioIge.Valor = fechaInicioIge
	filaPersona.FechaFinIge.Valor = fechaFinIge
	filaPersona.FechaInicioLma.Valor = fechaInicioLma
	filaPersona.FechaFinLma.Valor = fechaFinLma
	filaPersona.FechaInicioVac.Valor = fechaInicioVac
	filaPersona.FechaFinVac.Valor = fechaFinVac
	filaPersona.FechaInicioVct.Valor = fechaInicioVct
	filaPersona.FechaFinVct.Valor = fechaFinVct
	filaPersona.FechaInicioIrl.Valor = fechaInicioIrl
	filaPersona.FechaFinIrl.Valor = fechaFinIrl
	filaPersona.IbcOtrosParaFiscales.Valor = 0
	filaPersona.HorasLaboradas.Valor = 0

	personaNovedad = filaPersona
	return
}

// escribirDiasCotizadosAux escribe los días a cotizar para la fila auxiliar de novedades
func escribirDiasCotizadosAux(filaPersona *models.PlanillaTipoE, diasCotizados int) {
	filaPersona.DiasPension.Valor = diasCotizados
	filaPersona.DiasSalud.Valor = diasCotizados
	filaPersona.DiasArl.Valor = 0
	filaPersona.DiasCaja.Valor = diasCotizados
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

	novedadSeguridadSocial := "seg_social(" + novedad + ", " + strconv.Itoa(fechaInicioNovedad.Year()) + ", " +
		strconv.Itoa(int(fechaInicioNovedad.Month())) + ", " + strconv.Itoa(fechaInicioNovedad.Day()) + ", " +
		strconv.Itoa(fechaFinNovedad.Year()) + ", " + strconv.Itoa(int(fechaFinNovedad.Month())) + ", " +
		strconv.Itoa(fechaInicioNovedad.Day()) + ")."
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
	informacionUpcs = make(map[string][]models.UpcAdicional)
	err := getJson("http://"+beego.AppConfig.String("segSocialService")+
		"/upc_adicional"+
		"?limit=-1", &beneficiariosAdicionales)
	if err != nil {
		ImprimirError("error en buscarUpcAsociada()", err)
		return nil, err
	}

	if beneficiariosAdicionales[0].Id > 0 { // Esto es para verificar que el arreglo que devuelve el servicio tiene al menos un elemento con valor
		for _, beneficiario := range beneficiariosAdicionales {
			mapaBeneficiarios[beneficiario.PersonaAsociada] += int(beneficiario.TipoUpc.Valor)
			informacionUpcs[beneficiario.PersonaAsociada] = append(informacionUpcs[beneficiario.PersonaAsociada], beneficiario)
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
		idProveedor := strconv.FormatInt(value.Persona, 10)
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

func establecerNovedadesExterior(idPersona, idPreliquidacion string, filaPersona *models.PlanillaTipoE) {
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
		novedadPersona = true
		exterior = true
		filaPersona.CodigoDepartamento = models.Columna{Valor: "", Longitud: 2}
		filaPersona.CodigoMunicipio = models.Columna{Valor: "", Longitud: 3}
	}
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
			incapacidadGeneral = true
			diasIncapacidad = int(value.DiasLiquidados)
			diasIrl = int(value.DiasLiquidados)
			diasArl = diasIrl
			fechaInicioIge = fechaInicioTemp
			fechaFinIge = fechaFinTemp
			novedadPersona = true
		case "incapacidad_laboral":
			incapacidadGeneral = true
			diasIncapacidad = int(value.DiasLiquidados)
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
}

// marcarNovedad marca el valor para la novedad o la deja vacia
func marcarNovedadConOpciones(novedad bool, valor string) string {
	if novedad {
		return valor
	} else {
		return ""
	}
}

// marcarNovedad una novedad con el valor que se le da o la deja vacia
func marcarNovedad(novedad bool, valor string) string {
	if novedad {
		return valor
	}
	return ""
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

func formatoDato(texto string, longitud int) (textoEscribir string) {
	for _, r := range texto {
		textoEscribir += string(r)
	}
	for i := 0; i < longitud-len(texto); i++ {
		textoEscribir += " "
	}
	return
}
