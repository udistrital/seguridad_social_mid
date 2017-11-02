package utilidades

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
)

func FillStruct(m interface{}, s interface{}) (err error) {
	j, _ := json.Marshal(m)
	err = json.Unmarshal(j, s)
	return
}

func SetField(obj interface{}, name string, value interface{}) error {

	structValue := reflect.ValueOf(obj).Elem()
	fieldVal := structValue.FieldByName(name)

	if !fieldVal.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !fieldVal.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	val := reflect.ValueOf(value)

	if fieldVal.Type() != val.Type() {

		if m, ok := value.(map[string]interface{}); ok {

			// if field value is struct
			if fieldVal.Kind() == reflect.Struct {
				return FillStruct(m, fieldVal.Addr().Interface())
			}

			// if field value is a pointer to struct
			if fieldVal.Kind() == reflect.Ptr && fieldVal.Type().Elem().Kind() == reflect.Struct {
				if fieldVal.IsNil() {
					fieldVal.Set(reflect.New(fieldVal.Type().Elem()))
				}
				// fmt.Printf("recursive: %v %v\n", m,fieldVal.Interface())
				return FillStruct(m, fieldVal.Interface())
			}

		}

		return fmt.Errorf("Provided value type didn't match obj field type")
	}

	fieldVal.Set(val)
	return nil

}

func FillDataStruct(m map[string]interface{}, s interface{}) error {
	for k, v := range m {
		err := SetField(s, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}
func FillStructDeep(m map[string]interface{}, fields string, s interface{}) (err error) {
	f := strings.Split(fields, ".")
	if len(f) == 0 {
		err = errors.New("invalid fields.")
		return
	}

	var aux map[string]interface{}
	var load interface{}
	for i, value := range f {

		if i == 0 {
			//fmt.Println(m[value])
			FillStruct(m[value], &load)
		} else {
			FillStruct(load, &aux)
			FillStruct(aux[value], &load)
			//fmt.Println(aux[value])
		}
	}
	j, _ := json.Marshal(load)
	err = json.Unmarshal(j, s)
	return
}

//funcion para generar canales de interface{}
func GenChanInterface(mp ...interface{}) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		for _, ch := range mp {
			out <- ch
		}
		close(out)
	}()
	return out
}

func digester(done <-chan interface{}, f func(interface{}, ...interface{}) interface{}, params []interface{}, in <-chan interface{}, out chan<- interface{}) {
	for intfc := range in {
		res := f(intfc, params...)
		select {
		case out <- res:
		case <-done:
			return
		}
	}
}

//funcion para administrar las go rutines armadas para la consulta de solicitudes de rp.
func Digest(done <-chan interface{}, f func(interface{}, ...interface{}) interface{}, in <-chan interface{}, params []interface{}) (outchan <-chan interface{}) {
	out := make(chan interface{})
	var wg sync.WaitGroup
	const numDigesters = 20
	wg.Add(numDigesters)
	for i := 0; i < numDigesters; i++ {
		go func() {
			digester(done, f, params, in, out)
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
