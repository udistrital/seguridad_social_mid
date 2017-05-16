package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:DescSeguridadSocialController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:DescSeguridadSocialController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:DescSeguridadSocialController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:DescSeguridadSocialController"],
		beego.ControllerComments{
			Method: "GetOne",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:DescSeguridadSocialController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:DescSeguridadSocialController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:DescSeguridadSocialController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:DescSeguridadSocialController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:DescSeguridadSocialController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:DescSeguridadSocialController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:DescSeguridadSocialDetalleController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:DescSeguridadSocialDetalleController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:DescSeguridadSocialDetalleController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:DescSeguridadSocialDetalleController"],
		beego.ControllerComments{
			Method: "GetOne",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:DescSeguridadSocialDetalleController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:DescSeguridadSocialDetalleController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:DescSeguridadSocialDetalleController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:DescSeguridadSocialDetalleController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:DescSeguridadSocialDetalleController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:DescSeguridadSocialDetalleController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:EdadUpcController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:EdadUpcController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:EdadUpcController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:EdadUpcController"],
		beego.ControllerComments{
			Method: "GetOne",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:EdadUpcController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:EdadUpcController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:EdadUpcController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:EdadUpcController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:EdadUpcController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:EdadUpcController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:TipoPagoSeguridadSocialController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:TipoPagoSeguridadSocialController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:TipoPagoSeguridadSocialController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:TipoPagoSeguridadSocialController"],
		beego.ControllerComments{
			Method: "GetOne",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:TipoPagoSeguridadSocialController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:TipoPagoSeguridadSocialController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:TipoPagoSeguridadSocialController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:TipoPagoSeguridadSocialController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:TipoPagoSeguridadSocialController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:TipoPagoSeguridadSocialController"],
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

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:TipoZonaUpcController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:TipoZonaUpcController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:TipoZonaUpcController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:TipoZonaUpcController"],
		beego.ControllerComments{
			Method: "GetOne",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:TipoZonaUpcController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:TipoZonaUpcController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:TipoZonaUpcController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:TipoZonaUpcController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:TipoZonaUpcController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:TipoZonaUpcController"],
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

}
