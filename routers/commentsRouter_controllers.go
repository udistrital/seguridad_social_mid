package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["ss_api_mid/controllers:SegSocialController"] = append(beego.GlobalControllerRouter["ss_api_mid/controllers:SegSocialController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["ss_api_mid/controllers:SegSocialController"] = append(beego.GlobalControllerRouter["ss_api_mid/controllers:SegSocialController"],
		beego.ControllerComments{
			Method: "GetOne",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["ss_api_mid/controllers:SegSocialController"] = append(beego.GlobalControllerRouter["ss_api_mid/controllers:SegSocialController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["ss_api_mid/controllers:SegSocialController"] = append(beego.GlobalControllerRouter["ss_api_mid/controllers:SegSocialController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			Params: nil})

	beego.GlobalControllerRouter["ss_api_mid/controllers:SegSocialController"] = append(beego.GlobalControllerRouter["ss_api_mid/controllers:SegSocialController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			Params: nil})

	beego.GlobalControllerRouter["ss_api_mid/controllers:SegSocialDetalleController"] = append(beego.GlobalControllerRouter["ss_api_mid/controllers:SegSocialDetalleController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["ss_api_mid/controllers:SegSocialDetalleController"] = append(beego.GlobalControllerRouter["ss_api_mid/controllers:SegSocialDetalleController"],
		beego.ControllerComments{
			Method: "GetOne",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["ss_api_mid/controllers:SegSocialDetalleController"] = append(beego.GlobalControllerRouter["ss_api_mid/controllers:SegSocialDetalleController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["ss_api_mid/controllers:SegSocialDetalleController"] = append(beego.GlobalControllerRouter["ss_api_mid/controllers:SegSocialDetalleController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			Params: nil})

	beego.GlobalControllerRouter["ss_api_mid/controllers:SegSocialDetalleController"] = append(beego.GlobalControllerRouter["ss_api_mid/controllers:SegSocialDetalleController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			Params: nil})

	beego.GlobalControllerRouter["ss_api_mid/controllers:TipoUpcController"] = append(beego.GlobalControllerRouter["ss_api_mid/controllers:TipoUpcController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["ss_api_mid/controllers:TipoUpcController"] = append(beego.GlobalControllerRouter["ss_api_mid/controllers:TipoUpcController"],
		beego.ControllerComments{
			Method: "GetOne",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["ss_api_mid/controllers:TipoUpcController"] = append(beego.GlobalControllerRouter["ss_api_mid/controllers:TipoUpcController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["ss_api_mid/controllers:TipoUpcController"] = append(beego.GlobalControllerRouter["ss_api_mid/controllers:TipoUpcController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			Params: nil})

	beego.GlobalControllerRouter["ss_api_mid/controllers:TipoUpcController"] = append(beego.GlobalControllerRouter["ss_api_mid/controllers:TipoUpcController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			Params: nil})

	beego.GlobalControllerRouter["ss_api_mid/controllers:UpcAdicionalController"] = append(beego.GlobalControllerRouter["ss_api_mid/controllers:UpcAdicionalController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["ss_api_mid/controllers:UpcAdicionalController"] = append(beego.GlobalControllerRouter["ss_api_mid/controllers:UpcAdicionalController"],
		beego.ControllerComments{
			Method: "GetOne",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["ss_api_mid/controllers:UpcAdicionalController"] = append(beego.GlobalControllerRouter["ss_api_mid/controllers:UpcAdicionalController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["ss_api_mid/controllers:UpcAdicionalController"] = append(beego.GlobalControllerRouter["ss_api_mid/controllers:UpcAdicionalController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			Params: nil})

	beego.GlobalControllerRouter["ss_api_mid/controllers:UpcAdicionalController"] = append(beego.GlobalControllerRouter["ss_api_mid/controllers:UpcAdicionalController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			Params: nil})

	beego.GlobalControllerRouter["ss_api_mid/controllers:ZonaController"] = append(beego.GlobalControllerRouter["ss_api_mid/controllers:ZonaController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["ss_api_mid/controllers:ZonaController"] = append(beego.GlobalControllerRouter["ss_api_mid/controllers:ZonaController"],
		beego.ControllerComments{
			Method: "GetOne",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["ss_api_mid/controllers:ZonaController"] = append(beego.GlobalControllerRouter["ss_api_mid/controllers:ZonaController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["ss_api_mid/controllers:ZonaController"] = append(beego.GlobalControllerRouter["ss_api_mid/controllers:ZonaController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			Params: nil})

	beego.GlobalControllerRouter["ss_api_mid/controllers:ZonaController"] = append(beego.GlobalControllerRouter["ss_api_mid/controllers:ZonaController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			Params: nil})

}
