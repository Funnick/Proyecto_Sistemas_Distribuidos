# Plataforma de Agentes

## Integrantes

- Manuel Antonio Vilas Valiente
- Andrés León Almaguer
- Miguel Alejandro Rodríguez

## Tabla de contenidos

[TOC]

## Objetivos

Implementar una aplicación capaz de gestionar la suscripción de servicios como agentes en una plataforma, y posteriormente ser capaces de acceder a estos para revisiones, modificaciones o eliminaciones. Además, se busca que la aplicación sea distribuida, de forma transparente para el cliente y que garantice que la caída de servidores provoque la menor cantidad de problemas posibles.

## Server

La API del servidor brinda las siguientes funcionalidades:

| Función                      | Descripción                            |
| ---------------------------- | -------------------------------------- |
| `CreateNewAgent(...)`        | Registrar un agente en la plataforma.  |
| `DeleteAgent(...)`           | Eliminar un agente.                    |
| `UpdateAgent(...)`           | Editar los campos de un agente.        |
| `SearchAgentByName(...)`     | Buscar un agente por su nombre.        |
| `SearchAgentByFunction(...)` | Buscar un agente por su funcionalidad. |

## DHT-Chord

Explicación de las particularides de lo que hacemos con Chord

## CLI

Explicación de lo que hace el cli

### Manual de Usuario

Cómo se usa?

## Replicación

Replicamos en tu nodo sucesor en chord

## Consistencia

Realizamos una única operación atómica cuando los datos son modificados, es decir, en el momento que se está realizando la modificación, esta se actualiza simulatáneamente en todos los lugares donde la información está replicada.32

## Persistencia

En *Golang*, el encoding/binary de paquetes ofrece una manera conveniente de codificar estructuras de datos con valores numéricos en una representación binaria de tamaño fijo. Puede usarse directamente o como base de protocolos binarios personalizados. Debido a que el paquete es compatible con las interfaces de flujo de IO, se puede integrar fácilmente en programas de comunicación o almacenamiento utilizando primitivas de IO de transmisión.

Además, los **gobinary** (.gob) son el tipo de archivo que más rápido *Golang* codifica y decodifica, agregándole velocidad al proceso de cargado y salvado de la información del sistema, lo cual siempre es beneficioso porque disminuye el espacio de tiempo en el que, mientras se realizan cambios, una caída del sistema pueda perjudicar la correcta persistencia de los datos. Este detalle de la velocidad también nos ayuda con la consistencia.