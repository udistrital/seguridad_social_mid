package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/udistrital/b/ss_mid_api/controllers:DescSeguridadSocialController"] = append(beego.GlobalControllerRouter["github.com/udistrital/b/ss_mid_api/controllers:DescSeguridadSocialController"],
		beego.ControllerComments{
			Method:           "CalcularSegSocial",
			Router:           `CalcularSegSocial/:id`,
			AllowHTTPMethods: []string{"get"},
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/b/ss_mid_api/controllers:DescSeguridadSocialDetalleController"] = append(beego.GlobalControllerRouter["github.com/udistrital/b/ss_mid_api/controllers:DescSeguridadSocialDetalleController"],
		beego.ControllerComments{
			Method:           "GetNovedadesPorPersona",
			Router:           `GetNovedadesPorPersona/:persona`,
			AllowHTTPMethods: []string{"get"},
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/b/ss_mid_api/controllers:FuncionarioProveedorController"] = append(beego.GlobalControllerRouter["github.com/miguelramirez93/titan_api_crud/controllers:FuncionarioProveedorController"],
		beego.ControllerComments{
			Method:           "ConsultarIDProveedor",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/b/ss_mid_api/controllers:DescSeguridadSocialController"] = append(beego.GlobalControllerRouter["github.com/udistrital/b/ss_mid_api/controllers:DescSeguridadSocialController"],
		beego.ControllerComments{
			Method:           "GetConceptosIbc",
			Router:           `ConceptosIbc`,
			AllowHTTPMethods: []string{"get"},
			Params:           nil})
}
