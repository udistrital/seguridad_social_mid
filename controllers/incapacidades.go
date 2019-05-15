package controllers

import (
	"fmt"
	"log"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/manucorporat/try"
	"github.com/udistrital/ss_mid_api/models"
)

// IncapacidadesController operations for Incapacidades
type IncapacidadesController struct {
	beego.Controller
}

// URLMapping ...
func (c *IncapacidadesController) URLMapping() {
	c.Mapping("BuscarPersonas", c.BuscarPersonas)
}

// BuscarPersonas ...
// @Title BuscarPersonas
// @Description obtiene todas las personas que pueden aplicar a cualquier nómina
// @Param	documento		query	string false		"documento de la persona"
// @Success 200 {object} interface{}
// @Failure 403
// @router /BuscarPersonas/:documento [get]
func (c *IncapacidadesController) BuscarPersonas() {
	var proveedores, contratos, personaNatural, respuesta []map[string]interface{}
	documento := c.Ctx.Input.Param(":documento")
	log.Println(documento)
	try.This(func() {
		if err := getJson("http://"+beego.AppConfig.String("administrativaService")+"informacion_proveedor?"+
			"limit=6&query=NumDocumento__icontains:"+documento+",TipoPersona:NATURAL", &proveedores); err != nil {
			panic(err)
		}

		for _, proveedor := range proveedores {
			idProveedor := strconv.Itoa(int(proveedor["Id"].(float64)))
			if err := getJson("http://"+beego.AppConfig.String("administrativaService")+"contrato_general?"+
				"query=Estado:true,Contratista:"+idProveedor+"&fields=Id,VigenciaContrato", &contratos); err != nil {
				fmt.Println("error en contrato general")
				panic(err)
			}

			if err := getJson("http://"+beego.AppConfig.String("administrativaService")+"informacion_persona_natural?"+
				"limit=1&query=Id:"+proveedor["NumDocumento"].(string), &personaNatural); err != nil {
				fmt.Println("error en contrato informacion_persona_natural")
				panic(err)
			}

			var contratosPropios []map[string]interface{} // contratos de cada proveedor

			for _, contrato := range contratos {
				numeroContrato := contrato["Id"].(string)
				vigenciaContrato := contrato["VigenciaContrato"].(float64)

				temp := map[string]interface{}{
					"NumeroContrato":   numeroContrato,
					"VigenciaContrato": vigenciaContrato,
				}
				contratosPropios = append(contratosPropios, temp)
			}

			resp := map[string]interface{}{
				"display":       proveedor["NomProveedor"],
				"value":         proveedor["NumDocumento"],
				"documento":     proveedor["NumDocumento"],
				"contratos":     contratosPropios,
				"tipoDocumento": personaNatural[0]["TipoDocumento"].(map[string]interface{})["Abreviatura"],
				"idProveedor":   idProveedor,
			}

			respuesta = append(respuesta, resp)
		}
		if respuesta == nil {
			respuesta = append(respuesta, map[string]interface{}{})
		}

		c.Data["json"] = respuesta

	}).Catch(func(e try.E) {
		fmt.Println("Error en GetPersonas() ", e)
		c.Data["json"] = models.Alert{Code: "error"}
	})
	c.ServeJSON()
}

// IncapacidadesPorPersona ...
// @Title IncapacidadesPorPersona
// @Description obtiene todas las incapacidades activdas de una persona
// @Param	documento		query	string false		"documento de la persona"
// @Success 200 {object} interface{}
// @Failure 403
// @router /incapacidadesPersona/:contrato/:vigencia [get]
func (c *IncapacidadesController) IncapacidadesPorPersona() {
	contrato := c.Ctx.Input.Param(":contrato")
	vigencia := c.Ctx.Input.Param(":vigencia")

	var incapacidades []map[string]interface{}
	try.This(func() {
		incapacidadesLaborales, err := traerIncapacidades("incapacidad_laboral", contrato, vigencia)
		if err != nil {
			panic(err.Error())
		}

		incapacidaGenerales, err := traerIncapacidades("incapacidad_general", contrato, vigencia)
		if err != nil {
			panic(err.Error())
		}

		prorrogas, err := traerIncapacidades("prorroga_incapacidad", contrato, vigencia)
		if err != nil {
			panic(err.Error())
		}

		incapacidades = append(incapacidades, incapacidadesLaborales...)
		incapacidades = append(incapacidades, incapacidaGenerales...)
		incapacidades = append(incapacidades, prorrogas...)
		c.Data["json"] = incapacidades
	}).Catch(func(e try.E) {
		log.Panicf("Error en IncapacidadesPorPersona() ", e)
		c.Data["json"] = map[string]interface{}{"error": e}
	})

	c.ServeJSON()
}

func traerIncapacidades(tipoIncapacidad, contrato, vigencia string) (incapacidades []map[string]interface{}, err error) {
	var detalleNovedad []map[string]interface{}
	err = getJson("http://"+beego.AppConfig.String("titanServicio")+"/concepto_nomina_por_persona?query=Concepto.Nombreconcepto:"+tipoIncapacidad+
		",NumeroContrato:"+contrato+",VigenciaContrato:"+vigencia+",Activo:true&limit=0", &incapacidades)

	for i, v := range incapacidades {
		conceptoNominaPorPesona := strconv.Itoa(int(v["Id"].(float64)))
		err = getJson("http://"+beego.AppConfig.String("segSocialService")+"/detalle_novedad_seguridad_social?"+
			"query=ConceptoNominaPorPersona:"+conceptoNominaPorPesona+"&limit=1", &detalleNovedad)
		incapacidades[i]["Codigo"] = detalleNovedad[0]["Descripcion"]
	}
	return
}
