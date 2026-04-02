# Ejercicio 2 — Evaluación técnica inicial

## Objetivo

El objetivo es construir una aplicación mobile para gestión de tareas, integrando una API existente y considerando autenticación, roles y estados de negocio.

El foco inicial es asegurar que el equipo arranque sin bloqueos, con reglas claras y decisiones mínimas tomadas.

---

## Supuestos iniciales

Parto de estos supuestos:

- La app será mobile (iOS/Android)
- Existe una API Omnicanal externa
- El equipo es mayormente junior/mid
- Se prioriza velocidad de entrega sobre perfección
- El sistema debe poder evolucionar sin rehacerse

---

## Plan de actividades

### Fase 1 — Entendimiento (1–2 días)

Primero alineo con negocio:

- Qué puede hacer cada rol  
- Cómo se integra la API externa  
- Qué flujos son críticos  
- Qué significa “listo” para la primera versión  

En un CRM que lideré en MercadoLibre, este paso evitó rehacer flujos completos más adelante.

---

### Fase 2 — Diseño técnico (2–3 días)

No intento resolver todo, solo lo necesario para arrancar:

- modelo de dominio  
- contratos de API  
- estrategia de auth  
- puntos de integración externa  

---

### Fase 3 — Setup (1–2 días)

- repo  
- estructura base  
- entorno local con Docker  
- convenciones  

Esto es clave con equipos junior: si el setup falla, se pierde tiempo desde el día 1.

---

### Fase 4 — Desarrollo iterativo

Trabajo por flujos:

1. autenticación  
2. usuarios  
3. tareas  
4. flujos del ejecutor  
5. flujos del auditor  

Iteraciones cortas, con review y tests.

---

### Fase 5 — Validación

- testing funcional  
- ajustes  
- revisión de seguridad  
- documentación  

---

## Decisión mobile: híbrida vs nativa

Elegiría React Native o Flutter.

El foco está en llegar rápido al mercado con un equipo que probablemente no tenga especialistas en iOS y Android nativo. Para este tipo de aplicación, la performance no es un factor crítico, por lo que el trade-off es aceptable.

Si en el futuro se necesita performance nativa, se puede migrar por pantalla sin rehacer todo.

---

## Comunicación mobile ↔ backend

Para este caso optaría por un BFF (Backend for Frontend).

Esto permite simplificar la lógica en mobile, adaptar las respuestas según el flujo de UI y desacoplar la app de cambios internos del backend. Consumir REST directo sería viable, pero con roles diferenciados y estados de negocio, termina trasladando complejidad innecesaria al frontend.

---

## Arquitectura backend

Uso una arquitectura en capas simple:

- handlers  
- use cases  
- dominio  
- infraestructura  

Más que la estructura, lo importante es dónde viven las reglas. Las reglas están en el dominio, los use cases orquestan y los handlers solo traducen requests. Eso hace que el sistema sea testeable y fácil de entender para el equipo.

---

## Integración con API Omnicanal

Este es el principal riesgo técnico del sistema.

No controlamos disponibilidad, latencia ni cambios de contrato, por lo que la integración debe estar encapsulada en una capa específica. También es importante manejar timeouts y errores de forma explícita, para evitar que fallas externas se propaguen sin control.

---

## Observabilidad

Si la API externa falla o se vuelve lenta, necesito detectarlo rápido.

Implementaría logging estructurado, métricas básicas de latencia y error, y alertas sobre la integración. En un sistema de ventas, no detectar esto a tiempo impacta directamente en la operación.

---

## Riesgos

El mayor riesgo es la integración con la API externa. Si falla, impacta directamente en el negocio, por lo que se prioriza aislarla, monitorearla y manejar sus errores de forma controlada.

También existen riesgos en la definición de permisos y en la consistencia de estados de las tareas. Estos se mitigan centralizando reglas en el dominio y validando transiciones de forma estricta. La sincronización con mobile puede generar inconsistencias menores, pero se resuelve con endpoints idempotentes y manejo claro de errores.

---

## Conclusión

Mi enfoque en preventa es empezar simple, proteger las reglas de negocio y no vender arquitectura que el equipo no va a poder mantener.