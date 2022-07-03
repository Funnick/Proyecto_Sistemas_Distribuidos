# Plataforma de Agentes

## Integrantes

- Manuel Antonio Vilas Valiente
- Andrés León Almaguer
- Miguel Alejandro Rodríguez

## Tabla de contenidos

[TOC]

## Objetivos

Implementar una aplicación capaz de gestionar la suscripción de servicios como agentes en una plataforma, y posteriormente ser capaces de acceder a estos para revisiones, modificaciones o eliminaciones. Además, se busca que la aplicación sea distribuida, de forma transparente para el cliente y que garantice que la caída de servidores provoque la menor cantidad de problemas posibles.

## Tecnologías usadas

El proyecto está desarrollado 100% en **Golang** en su versión 1.18, lenguaje escogido por sus beneficios a la hora de realizar rutinas en paralelo y por las facilidades que brinda su biblioteca *net* a la hora de gestionar los recursos del ambiente de la red.

## La Plataforma

El proyecto está dividido en tres paquetes, el *package server*, donde se desarrolla toda la parte del servidor, el *package chord*, en el cual se estructura todo el sistema de la tabla de hash distribuida sobre la cual se sostendrá el sistema, y por último el *package client*, separado completamente de los otros dos y pensado para la interacción de los usuarios con la plataforma.

## Server

La API del servidor brinda las siguientes funcionalidades:

| Función                  | Descripción                           |
| ------------------------ | ------------------------------------- |
| `CreateNewAgent(...)`    | Registrar un agente en la plataforma. |
| `DeleteAgent(...)`       | Eliminar un agente.                   |
| `UpdateAgent(...)`       | Editar los campos de un agente.       |
| `SearchAgentByName(...)` | Buscar un agente por su nombre.       |
| `SearchAgentByDesc(...)` | Buscar un agente por su descipción.   |

## DHT-Chord

Explicación de las particularides de lo que hacemos con Chord

## CLI

La aplicación consiste en una interfaz de usuario basada en linea de comandos buscando la simplicidad y la portabilidad. Las interacciones con la aplicación tiene un conjunto de comandos y opciones previamente seleccionados que permitirán las comunicaciones con el servidor mediante el cliente.

### Comandos

#### Crear un nuevo agente

Para crear un nuevo agente invocaremos el comando **create-agent** o su forma abreviada **C**. Cuando se hace uso de este comando es necesario introducir los parametros que seran los que definiran el nuevo agente usando las opciones o flags siguientes, sin importar su orden:

* **-name**, **-n** el nombre que tendrá el agente.
* **-ip**, **-i** es el IP donde se encontrara el agente a crear.
* **-port**, **-pr** representa el puerto donde se encontrará el nuevo agente.
* **-password**, -**pass**, **-p** será la contraseña que guardará el acceso al nuevo agente.
* **-description**, **-des** una descripcion del agente.
* **-doc** la documentación de uso del agente.

Ejemplos de uso del programa:

```shell
go run cli.go C -name Pepe -ip 10.8.100.2 -port 444 -password contraseña -description lorem ipsum -doc documentation
```

```shell
go run cli.go create-agent -doc documentation -password contraseña -port 444 -description lorem ipsum -ip 10.8.100.2
```



#### Eliminar agentes

Para eliminar un agente invocaremos el comando **delete-agent** o su forma abreviada **D**. Cuando se hace uso de este comando es necesario introducir los parámetros que serán los que definirán el nuevo agente usando las opciones o flags siguientes, sin importar su orden:

* **-name**, **-n** representa el id mediante el cual identificaremos el agente.

* **-password**, -**pass**, **-p** será la contraseña que guardará el acceso al agente.

  

Ejemplos de uso del programa:

``` shell
go run cli.go D -n name -password contraseña 
```

``` shell
go run cli.go delete-agent -p contraseña -name agentname 
```



#### Buscar agentes por nombre 

Para buscar dentro de la lista de agentes disponibles usando un nombre invocaremos el comando **search-name-agent** o su forma abreviada **S**. Cuando se hace uso de este comando es necesario introducir como parámetro la cadena a buscar:

Ejemplos de uso del programa:

``` shell
go run cli.go S lorem impsum 
```

``` shell
go run cli.go search-name-agent lorem impsum 
```



#### Buscar agentes por descripción 

Para buscar dentro de la lista de agentes disponibles por descripción invocaremos el comando **search-desc-agent** o su forma abreviada **Sd**. Cuando se hace uso de este comando es necesario introducir como parámetro la cadena a buscar:

Ejemplos de uso del programa:

``` shell
go run cli.go Sd lorem impsum 
```

``` shell
go run cli.go search-desc-agent lorem impsum 
```



#### Actualizar agentes

Para actualizar los valores de un agente invocaremos el comando **update-agent** o su forma abreviada **U**. Cuando se hace uso de este comando es necesario introducir los parámetros que serán los que definirán los nuevos valores del agente usando las opciones o *flags* siguientes, sin importar su orden:

* **-name**, **-n** el nombre por cual accederemos al agente
* **-password**, -**pass**, **-p** la contraseña para obtener acceso al agente
* **-ip**, **-p** es el nuevo IP donde se encontrara el agente a crear.
* **-port**, **pr** representa el nuevo puerto donde se encontrará el agente.
* **-new-password**, **-np** será la nueva contraseña que guardará el acceso al nuevo agente.
* **-description**, **-des** la nueva descripcion del agente.
* **-doc** la nueva documentación de uso del agente.

Ejemplos de uso del programa:

```shell
go run cli.go U -name agent-name -password contraseña -ip 10.8.100.2 -port 444 -new-password contraseña -description lorem ipsum -doc documentation
```

```shell
go run cli.go update-agent -doc documentation -password contraseña -port 444 -description lorem ipsum -ip 10.8.100.2 -name agent-name -new-password contraseña
```



#### Listar Agentes 

Para obtener una lista de agente disponibles invocaremos el comando **get-agent** o su forma abreviada **A**.

Ejemplos de uso del programa:

```shell
go run cli.go get-agents
```

```sh
go run cli.go A 
```



### Opciones 

Para usar una configuración diferente a la establecida por *default* usaremos la opción o *flag*: **--config**, **-c** usando como parámetro la dirección de la configuración deseada.

Ejemplos de uso del programa:

```shell
go run cli.go --config /path/to/config get-agents 
```

## Replicación

La replicación cae en manos de la tabla de hash distribuida. Cada nodo replicará los datos almacenados en la tabla de su predecesor en su propia tabla, permitiendo así que si su predecesor cae, su información no se perderá.

## Consistencia

Realizamos una única operación atómica cuando los datos son modificados, es decir, en el momento que se está realizando la modificación, esta se actualiza simulatáneamente en todos los lugares donde la información está replicada.

## Persistencia

En *Golang*, el encoding/binary de paquetes ofrece una manera conveniente de codificar estructuras de datos con valores numéricos en una representación binaria de tamaño fijo. Puede usarse directamente o como base de protocolos binarios personalizados. Debido a que el paquete es compatible con las interfaces de flujo de IO, se puede integrar fácilmente en programas de comunicación o almacenamiento utilizando primitivas de IO de transmisión.

Además, los **gobinary** (.gob) son el tipo de archivo que más rápido *Golang* codifica y decodifica, agregándole velocidad al proceso de cargado y salvado de la información del sistema, lo cual siempre es beneficioso porque disminuye el espacio de tiempo en el que, mientras se realizan cambios, una caída del sistema pueda perjudicar la correcta persistencia de los datos. Este detalle de la velocidad también nos ayuda con la consistencia.