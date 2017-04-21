package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"ss_api_mid/models"
	"time"

	"github.com/astaxie/beego"
)

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
	concepto(ud, descuento, porcentaje, salud, 5, 0.085,	2017). 	%%descuento salud ud
	concepto(ud, descuento, porcentaje, pension, 5, 0.12, 2017).	%%descuento pension ud
	concepto(ud, descuento, porcentaje, arl, 5, 0.00522, 2017). %%descuento de ARL

	%%	Hechos para pensionados
	concepto(X, devengo, porcentaje, salud, pensionado, 0.12, 2017).	%%descuento de salud pensionado

	%% Hechos para aportes parafiscales
	concepto(ud, descuento, porcentaje, caja,	5, 0.04, 2017).	%%caja de compensa familiar
	concepto(ud, descuento, porcentaje,	icbf, 5, 0.03, 2017).	%%ICBF


	%% Pagos de salud funcionario
	v_salud_func(1, 105080).
	v_salud_func(2, 145200).
	v_salud_func(3, 90240).

	%% Pagos de pension funcionario
	v_pen_func(1, 105080).
	v_pen_func(2, 145200).
	v_pen_func(3, 90240).

	%%salario minimo legal mensual vigente
	smlmv(737717, 2017).

	%%		SALUD
	v_salud_ud(I,Y) :- concepto(ud,Z,T,salud,5,V,2017), ibc(I,W,salud), Y is V * W.
	v_total_salud(X,T) :- v_salud_func(X,Y), v_salud_ud(X,U), T is (Y + U).

	%%		PENSION
	v_pen_ud(I,Y) :- concepto(ud,Z,T,pension,5,V,2017), ibc(I,W,salud), Y is V * W.
	v_total_pen(X,T) :- v_pen_func(X,Y), v_pen_ud(X,U), T is Y + U.

	%%		ARL
	v_arl(I,Y) :- concepto(ud,Z,T,arl,5,V,2017), ibc(I,W,riesgos), (novedad(I,2324) -> Y is 0 * W; novedad(I,2326) -> Y is 0 * W;Y is V * W).

	%%		FONDO DE SOLIDARIDAD
	v_fondo1(X,S,D,Y) :- ibc(X,W,apf), smlmv(M,2017),
	(S is 4*M, W@>= S, D is 16*M, W@< D -> Y is W * 0.01;
	S is 16*M, W@>= S, D is 17*M, W@< D -> Y is W * 0.012;
	S is 17*M, W@>= S, D is 18*M, W@< D -> Y is W * 0.014;
	S is 18*M, W@>= S, D is 19*M, W@< D -> Y is W * 0.016;
	S is 19*M, W@>= S, D is 20*M, W@=< D -> Y is W * 0.018;
	S is 20*M, W@> S -> Y is W * 0.02).	%calculo de fondo de solidaridad 1

	%% PAGO UPC
	v_upc(I,Y,Z) :- ibc(I,W,salud), upc(Z,V,I), Y is W - V.

	%%		CAJA DE COMPENSACION FAMILIAR
	v_caja(I,Y) :- concepto(ud,Z,T,caja,5,V,2017), ibc(I,W,apf), Y is V * W.

	%%		ICBF
	v_icbf(I,Y) :- concepto(ud,Z,T,icbf,5,V,2017), ibc(I,W,apf), Y is V * W.
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
