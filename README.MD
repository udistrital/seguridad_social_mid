# ss_mid_api

# seguridad_social_crud

seguridad_social_crud, API de entidades que utiliza el sistema de Seguridad Social

Este proyecto está escrito en el lenguaje Go, utilizando el framework beego (https://beego.me/).

Se utiliza para administrar los modelos, transacciones y efectos de cambios directamente en la base de datos, nada de lógica. 


# Ejecución 

1. Clonar el repositorio

```go
  go get github.com/udistrital/ss_mid_api`
```

2. Moverse a la carpeta del proyecto

```sh
  cd $GOPATH/src/github.com/udistrital/ss_mid_api
```

3. Actualizar la rama develop

```sh
  git pull origin develop
```


4. Cambiar a la rama develop

```git
  git checkout develop
```

5. Asignar variables de entorno (ver fichero **/conf/app.conf**) 

6. espués de asignar las variables de entorno ejecutar

```
bee run
```

7. Para crear/actualizar la documentación y archivo swagger de forma automática ejecutar dentro del proyecto

```
bee run -downdoc=true -gendoc=true
```


## Derechos de Autor

This program is free software: you can redistribute it 
and/or modify it under the terms of the GNU General Public 
License as published by the Free Software Foundation, either
version 3 of the License, or (at your option) any later
version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

### UNIVERSIDAD DISTRITAL Francisco José de Caldas

### Oficina asesora de Sistemas 2019

### 2019
