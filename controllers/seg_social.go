package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/udistrital/ss_mid_api/golog"
	"github.com/udistrital/ss_mid_api/models"

	"github.com/astaxie/beego"
)

// SegSocialController oprations for SegSocial
type SegSocialController struct {
	beego.Controller
}

// URLMapping ...
func (c *SegSocialController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
	c.Mapping("CalcularSegSocial", c.CalcularSegSocial)
	c.Mapping("GetConceptosIbc", c.GetConceptosIbc)
}

// Post ...
// @Title Post
// @Description create SegSocial
// @Param	body		body 	models.SegSocial	true		"body for SegSocial content"
// @Success 201 {int} models.SegSocial
// @Failure 403 body is empty
// @router / [post]
func (c *SegSocialController) Post() {
	var v models.SegSocial
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if _, err := models.AddSegSocial(&v); err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = v
		} else {
			c.Data["json"] = err.Error()
		}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get SegSocial by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.SegSocial
// @Failure 403 :id is empty
// @router /:id [get]
func (c *SegSocialController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetSegSocialById(id)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

func (c *SegSocialController) GetConceptosIbc() {
	fmt.Println("GetConceptosIbc")
	var predicados []models.Predicado
	var conceptos []models.Concepto
	var conceptosIbc []models.ConceptosIbc
	err := getJson("http://"+beego.AppConfig.String("rulerServicio")+
		"/predicado?limit=-1&query=Nombre__startswith:conceptos_ibc,Dominio.Id:5", &predicados)
	errConceptoTitan := getJson("http://"+beego.AppConfig.String("titanServicio")+
		"/concepto?limit=-1", &conceptos)
	if err != nil && errConceptoTitan != nil {
		c.Data["json"] = err.Error() + errConceptoTitan.Error()
	} else {
		nombres := golog.GetString(FormatoReglas(predicados), "conceptos_ibc(X).", "X")
		for i := 0; i < len(predicados); i++ {
			for j := 0; j < len(conceptos); j++ {
				if nombres[i] == conceptos[j].NombreConcepto {
					aux := models.ConceptosIbc{
						Id:          predicados[i].Id,
						Nombre:      nombres[i],
						Descripcion: conceptos[j].AliasConcepto}
					conceptosIbc = append(conceptosIbc, aux)
					break
				} else {
				}
			}

		}
		c.Data["json"] = conceptosIbc
	}
	c.ServeJSON()
}

func (c *SegSocialController) CalcularSegSocial() {
	idStr := c.Ctx.Input.Param(":id")
	_, err := strconv.Atoi(idStr)
	var alertas []string
	var predicado []models.Predicado
	var errores string

	var detalleLiquidacion []models.DetalleLiquidacion
	var detalleLiquSalud []models.DetalleLiquidacion
	var detalleLiquPension []models.DetalleLiquidacion

	var pagosSeguridadSocial []models.PagosSeguridadSocial

	if err != nil {
		c.Data["json"] = err.Error()
	} else {

		err := getJson("http://"+beego.AppConfig.String("titanServicio")+"/detalle_liquidacion"+
			"?limit=0&query=Liquidacion.Id:"+idStr+",Concepto.NombreConcepto:ibc_liquidado", &detalleLiquidacion)

		errSalud := getJson("http://"+beego.AppConfig.String("titanServicio")+"/detalle_liquidacion"+
			"?limit=0&query=Liquidacion.Id:"+idStr+",Concepto.NombreConcepto:salud", &detalleLiquSalud)

		errPension := getJson("http://"+beego.AppConfig.String("titanServicio")+"/detalle_liquidacion"+
			"?limit=0&query=Liquidacion.Id:"+idStr+",Concepto.NombreConcepto:pension", &detalleLiquPension)

		fmt.Println(err, errSalud, errPension)
		errores = ""

		if errores != "" {
			fmt.Println("err en la peticion", errores)
			alertas = append(alertas, "error al traer detalle liquidacion")
			fmt.Println(alertas)
			c.Data["json"] = alertas

		} else {

			for index := 0; index < len(detalleLiquidacion); index++ {
				predicado = append(predicado, models.Predicado{Nombre: "ibc(" + strconv.Itoa(detalleLiquidacion[index].Persona) + "," + strconv.Itoa(int(detalleLiquidacion[index].ValorCalculado)) + ", salud)."})
				predicado = append(predicado, models.Predicado{Nombre: "ibc(" + strconv.Itoa(detalleLiquidacion[index].Persona) + "," + strconv.Itoa(int(detalleLiquidacion[index].ValorCalculado)) + ", riesgos)."})
				predicado = append(predicado, models.Predicado{Nombre: "ibc(" + strconv.Itoa(detalleLiquidacion[index].Persona) + "," + strconv.Itoa(int(detalleLiquidacion[index].ValorCalculado)) + ", apf)."})
				predicado = append(predicado, models.Predicado{Nombre: "v_salud_func(" + strconv.Itoa(detalleLiquSalud[index].Persona) + ", " + strconv.Itoa(int(detalleLiquSalud[index].ValorCalculado)) + ")."})
				predicado = append(predicado, models.Predicado{Nombre: "v_pen_func(" + strconv.Itoa(detalleLiquPension[index].Persona) + ", " + strconv.Itoa(int(detalleLiquPension[index].ValorCalculado)) + ")."})
				reglas := CargarReglasBase() + FormatoReglas(predicado) + CargarNovedades(idStr)
				//fmt.Println(reglas)
				ids := golog.GetInt64(reglas, "v_salud_ud(I,Y).", "I")
				saludUd := golog.GetFloat(reglas, "v_salud_ud(I,Y).", "Y")
				saludTotal := golog.GetFloat(reglas, "v_total_salud(X,T).", "T")
				pensionUd := golog.GetFloat(reglas, "v_pen_ud(I,Y).", "Y")
				pensionTotal := golog.GetFloat(reglas, "v_total_pen(X,T).", "T")
				arl := golog.GetFloat(reglas, "v_arl(I,Y).", "Y")
				//concepto := golog.GetOneInt64(reglas, "novedad(I,X).", "I")
				aux := models.PagosSeguridadSocial{
					Persona:      ids[index],
					SaludUd:      saludUd[index],
					SaludTotal:   saludTotal[index],
					PensionUd:    pensionUd[index],
					PensionTotal: pensionTotal[index],
					Arl:          arl[index]}

				//aux.Conceptos = GetConceptosPersona(int(ids[index]))

				pagosSeguridadSocial = append(pagosSeguridadSocial, aux)

			}
			c.Data["json"] = pagosSeguridadSocial
		}
		c.ServeJSON()

	}
}

func CargarNovedades(id string) (novedades string) {
	var conceptos []models.DetalleLiquidacion
	var predicado []models.Predicado

	errLincNo := getJson("http://"+beego.AppConfig.String("titanServicio")+"/detalle_liquidacion?limit=0&query=Liquidacion.Id:"+id+",Concepto.naturaleza:seguridad_social&fields=Concepto,Persona", &conceptos)

	if errLincNo != nil {
		fmt.Println("error en CargarNovedades()", errLincNo)
	} else {
		for index := 0; index < len(conceptos); index++ {
			predicado = append(predicado, models.Predicado{Nombre: "novedad(" + strconv.Itoa(int(conceptos[index].Persona)) + ", " + strconv.Itoa(conceptos[index].Concepto.Id) + ")."})
			novedades += predicado[index].Nombre + "\n"
		}
	}
	return
}

// GetAll ...
// @Title Get All
// @Description get SegSocial
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.SegSocial
// @Failure 403
// @router / [get]
func (c *SegSocialController) GetAll() {
	CargarReglasBase()
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 10
	var offset int64

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	// sortby: col1,col2
	if v := c.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}
	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	// query: k:v,k:v
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				c.Data["json"] = errors.New("Error: invalid query key/value pair")
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	l, err := models.GetAllSegSocial(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = l
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the SegSocial
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.SegSocial	true		"body for SegSocial content"
// @Success 200 {object} models.SegSocial
// @Failure 403 :id is not int
// @router /:id [put]
func (c *SegSocialController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.SegSocial{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := models.UpdateSegSocialById(&v); err == nil {
			c.Data["json"] = "OK"
		} else {
			c.Data["json"] = err.Error()
		}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the SegSocial
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *SegSocialController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeleteSegSocial(id); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
