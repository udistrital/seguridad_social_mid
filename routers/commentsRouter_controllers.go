package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:SegSocialController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:SegSocialController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:SegSocialController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:SegSocialController"],
		beego.ControllerComments{
			Method: "GetOne",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:SegSocialController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:SegSocialController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:SegSocialController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:SegSocialController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:SegSocialController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:SegSocialController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:SegSocialDetalleController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:SegSocialDetalleController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:SegSocialDetalleController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:SegSocialDetalleController"],
		beego.ControllerComments{
			Method: "GetOne",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:SegSocialDetalleController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:SegSocialDetalleController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:SegSocialDetalleController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:SegSocialDetalleController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:SegSocialDetalleController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:SegSocialDetalleController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:TipoUpcController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:TipoUpcController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:TipoUpcController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:TipoUpcController"],
		beego.ControllerComments{
			Method: "GetOne",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:TipoUpcController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:TipoUpcController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:TipoUpcController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:TipoUpcController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:TipoUpcController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:TipoUpcController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:UpcAdicionalController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:UpcAdicionalController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:UpcAdicionalController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:UpcAdicionalController"],
		beego.ControllerComments{
			Method: "GetOne",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:UpcAdicionalController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:UpcAdicionalController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:UpcAdicionalController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:UpcAdicionalController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:UpcAdicionalController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:UpcAdicionalController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:ZonaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:ZonaController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:ZonaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:ZonaController"],
		beego.ControllerComments{
			Method: "GetOne",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:ZonaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:ZonaController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:ZonaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:ZonaController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:ZonaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:ZonaController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			Params: nil})

}
