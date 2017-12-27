package golog

import (
	"fmt"
	"strconv"

	. "github.com/udistrital/golog"
)

func GetInt64(reglas string, regla_inyectada string, variable_a_obtener string) []int64 {
	m := NewMachine().Consult(reglas)
	resultados := m.ProveAll(regla_inyectada)
	var res []int64
	for _, solution := range resultados {
		aux, _ := strconv.ParseInt(fmt.Sprintf("%s", solution.ByName_(variable_a_obtener)), 10, 64)
		res = append(res, aux)
	}
	return res
}

func GetString(reglas string, regla_inyectada string, variable_a_obtener string) []string {
	m := NewMachine().Consult(reglas)
	resultados := m.ProveAll(regla_inyectada)
	var res []string
	for _, solution := range resultados {
		aux := fmt.Sprintf("%s", solution.ByName_(variable_a_obtener))
		res = append(res, aux)
	}

	return res
}

func GetFloat(reglas string, regla_inyectada string, variable_a_obtener string) []float64 {
	m := NewMachine().Consult(reglas)
	resultados := m.ProveAll(regla_inyectada)
	var res []float64
	for _, solution := range resultados {
		aux, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_(variable_a_obtener)), 64)
		res = append(res, aux)
	}

	return res
}

func GetOneString(reglas string, regla_inyectada string, variable_a_obtener string) string {
	m := NewMachine().Consult(reglas)
	resultados := m.ProveAll(regla_inyectada)
	var returnedString string
	for _, solution := range resultados {
		aux := fmt.Sprintf("%s", solution.ByName_(variable_a_obtener))
		returnedString = aux
	}
	return returnedString
}

func GetOneInt64(reglas string, regla_inyectada string, variable_a_obtener string) int64 {
	m := NewMachine().Consult(reglas)
	resultados := m.ProveAll(regla_inyectada)
	var returnedInt int64
	for _, solution := range resultados {
		aux, _ := strconv.ParseInt(fmt.Sprintf("%s", solution.ByName_(variable_a_obtener)), 10, 64)
		returnedInt = aux
	}
	return returnedInt
}
