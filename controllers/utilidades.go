package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/udistrital/ss_mid_api/models"

	"github.com/astaxie/beego"
)

func convertirMapa(arr []interface{}) map[string][]interface{} {
	var (
		arry     []interface{}
		contrato string
	)
	returnedMap := make(map[string][]interface{})

	for i := range arr {
		tempMap := arr[i].(map[string]interface{})
		for key, value := range tempMap {
			if key == "NumeroContrato" {
				arry = append(arry, value)
				contrato = value.(string)
			}
			returnedMap[contrato] = arr
		}
	}
	return returnedMap
}

func sendJson(url string, trequest string, target interface{}, datajson interface{}) error {
	b := new(bytes.Buffer)
	if datajson != nil {
		json.NewEncoder(b).Encode(datajson)
	}
	client := &http.Client{}
	req, err := http.NewRequest(trequest, url, b)
	r, err := client.Do(req)
	//r, err := http.Post(url, "application/json; charset=utf-8", b)
	if err != nil {
		beego.Error("error", err)
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func getJsonWSO2(urlp string, target interface{}) error {
	b := new(bytes.Buffer)
	//proxyUrl, err := url.Parse("http://10.20.4.15:3128")
	//http.DefaultTransport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	client := &http.Client{}
	req, err := http.NewRequest("GET", urlp, b)
	req.Header.Set("Accept", "application/json")
	r, err := client.Do(req)
	//r, err := http.Post(url, "application/json; charset=utf-8", b)
	if err != nil {
		beego.Error("error", err)
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func getJson(url string, target interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func diff(a, b time.Time) (year, month, day int) {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	oneDay := time.Hour * 5
	a = a.Add(oneDay)
	b = b.Add(oneDay)
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	year = int(y2 - y1)
	month = int(M2 - M1)
	day = int(d2 - d1)

	// Normalize negative values
	/*if day < 0{
				day = 0
			}
			if month < 0 {
	        month = 0
	    }*/
	if day < 0 {
		// days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}

	return
}

func describe(i interface{}) {
	fmt.Printf("(%v, %T)\n", i, i)
}

/*func CargarReglasBase(dominio string) (reglas string) {
	var reglasbase string = ``
	var v []models.Predicado

	if err:= getJson("http://localhost:8080/v1/predicado?limit=0&query=Dominio.Nombre:"+dominio, &v); err == nil {
		reglasbase = reglasbase + FormatoReglas(v)
	} else {
		fmt.Println("err: ", err)
	}
	fmt.Println(reglasbase)
	return reglasbase
}*/

func CargarReglasBase() (reglas string) {
	reglas = `
			%% 		HECHOS PARA ACTIVOS
			%%(ud, conceptoDeDescuento, porcentaje, concepto, nominaCorrespondiente, valorPorcentaje, vigencia).
			concepto(descuento, porcentaje, salud, X, 0.085,	2017). 	%%descuento salud ud
			concepto(descuento, porcentaje, pension, X, 0.12, 2017).	%%descuento pension ud
			concepto(descuento, porcentaje, arl, X, 0.00522, 2017). %%descuento de ARL

			%% HECHOS PARA CONTRATISTAS
			%%(ud, conceptoDeDescuento, porcentaje, concepto, nominaCorrespondiente, valorPorcentaje, vigencia).
			concepto(descuento, porcentaje, salud, contratistas, 0.125,	2017). 	%%descuento salud contratista
			concepto(descuento, porcentaje, pension, contratistas, 0.16, 2017).	%%descuento pension contratista
			concepto(descuento, porcentaje, arl, contratistas, 0.00522, 2017). %%descuento de ARL contratista

			%%		HECHOS PARA PENSIONADOS
			concepto(X, devengo, porcentaje, salud, pensionado, 0.12, 2017).	%%descuento de salud pensionado

			%% Hechos para aportes parafiscales
			concepto(descuento, porcentaje, caja,	5, 0.04, 2017).	%%caja de compensa familiar
			concepto(descuento, porcentaje,	icbf, 5, 0.03, 2017).	%%ICBF

			%%		NOVEDADES
			%%(descripcion, persona)
			novedad(exterior_familia, 0).
			novedad_persona(-1,-1).

			%%salario minimo legal mensual vigente
			smlmv(737717, 2017).

			%%		SALUD
			v_salud_ud(I,Y,C) :- concepto(Z,T,salud,X,V,2017), ibc(I,W,C,salud), (novedad_persona(N,I), novedad(N,U) -> Y is ((V * W) * U) approach 100; Y is (V * W) approach 100).
			v_total_salud(X,T) :- v_salud_func(X,Y), v_salud_ud(X,U), T is (Y + U) approach 100.
			v_salud_contratista(I,Y,C) :- concepto(Z,T,salud,contratista,V,2017), ibc(I,W,C,salud), Y is (V * W) approach 100.

			%%		PENSION
			v_pen_ud(I,Y,C) :- concepto(Z,T,pension,X,V,2017), ibc(I,W,C,salud), Y is (V * W) approach 100.
			v_total_pen(X,T) :- v_pen_func(X,Y), v_pen_ud(X,U), T is (Y + U) approach 100.
			v_pen_contratista(I,Y,C) :- concepto(Z,T,pension,contratista,V,2017), ibc(I,W,C,salud), Y is (V * W) approach 100.

			%%		ARL
			v_arl(I,Y) :- concepto(Z,T,arl,X,V,2017), ibc(I,W,C,riesgos), Y is (V * W) approach 100.

			%%		FONDO DE SOLIDARIDAD
			v_fondo1(X,S,D,Y) :- ibc(X,W,apf), smlmv(M,2017),
			(S is 4*M, W@>= S, D is 16*M, W@< D -> Y is W * 0.01;
				S is 16*M, W@>= S, D is 17*M, W@< D -> Y is W * 0.012;
				S is 17*M, W@>= S, D is 18*M, W@< D -> Y is W * 0.014;
				S is 18*M, W@>= S, D is 19*M, W@< D -> Y is W * 0.016;
				S is 19*M, W@>= S, D is 20*M, W@=< D -> Y is W * 0.018;
				S is 20*M, W@> S -> Y is W * 0.02), Y is Y approach 100.	%calculo de fondo de solidaridad 1

				%% 		PAGO UPC
				v_upc(I,Y,Z) :- ibc(I,W,salud,D), upc(Z,V,I), Y is W - V.

				%%		CAJA DE COMPENSACION FAMILIAR
				v_caja(I,Y) :- concepto(Z,T,caja,X,V,2017), ibc(I,W,apf), Y is (V * W) approach 100.

				%%		ICBF
				v_icbf(I,Y) :- concepto(Z,T,icbf,X,V,2017), ibc(I,W,apf), Y is (V * W) approach 100.
	`
	//fmt.Println(reglas)
	return
}

func FormatoReglas(v []models.Predicado) (reglas string) {
	var arregloReglas = make([]string, len(v))
	reglas = ""

	for i := 0; i < len(v); i++ {
		arregloReglas[i] = v[i].Nombre
	}

	for i := 0; i < len(arregloReglas); i++ {
		reglas = reglas + arregloReglas[i] + "\n"
	}

	return
}
