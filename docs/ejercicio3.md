# Ejercicio 3 — Buenas prácticas para liderar el equipo

## Contexto

Con equipos junior/mid el mayor desafío no es la tecnología sino la claridad. Mi rol es asegurar que puedan avanzar sin fricción.

---

## Diseño

Evito sobre-ingeniería.

Prefiero capas claras, responsabilidades simples y reglas de negocio en el dominio, no en los handlers. En equipos menos maduros, la complejidad técnica se traduce directamente en errores.

---

## Desarrollo

### Pull Requests obligatorios

Nadie mergea directo.

Pero lo importante no es el control, es el aprendizaje.

---

### Code review como formación

Uso el review para enseñar.

En un equipo que lideré, empecé a pedir que cada PR incluya una breve explicación de por qué eligió ese enfoque. Solo ese cambio mejoró la calidad de los PRs en pocas semanas.

---

### Cambios pequeños

Prefiero PRs chicos. Son más fáciles de revisar, reducen errores y permiten iterar rápido.

---

## Definition of Done

Para mí es clave, especialmente con equipos junior.

Una tarea no está lista si no tiene tests, no pasó code review o no respeta las reglas de negocio.

---

## Testing

No busco cobertura alta.

Prefiero validar reglas importantes y evitar bugs críticos antes que tener tests artificiales.

---

## Metodología

Trabajo con herramientas simples como GitHub Projects o Jira, con tareas claras y visibles. Evito procesos complejos que el equipo termina ignorando.

---

## Onboarding

Cuando entra un junior:

Día 1:
- explico arquitectura  
- muestro un flujo completo  

Primera semana:
- tareas pequeñas  
- revisión cercana  
- pairing si es necesario  

En mi experiencia, un junior que entiende el flujo completo desde el día 1 empieza a contribuir antes y con menos errores.

---

## Manejo de conflictos técnicos

Siempre aparece.

En un equipo anterior tuvimos un debate sobre dónde debían vivir las validaciones. Un desarrollador proponía resolverlo en los handlers por simplicidad, otro en el dominio.

Decidí llevarlas al dominio y expliqué el motivo: si cambia el transporte, las reglas no deberían romperse. Ese tipo de decisiones son las que marcan la mantenibilidad del sistema.

---

## Cómo mido progreso

No me baso solo en velocidad.

Miro la calidad de los PRs, cuántas correcciones necesitan y cuánto mejora la autonomía del desarrollador. Cuando un junior deja de depender constantemente, es una señal clara de avance.

---

## Deuda técnica

Con equipos junior, la deuda aparece rápido.

La dejo visible (issues o backlog), reservo capacidad en cada iteración y no la ignoro hasta que explote. Con perfiles junior, si no la nombrás explícitamente, no existe para ellos.

---

## Comunicación

Si el equipo no puede explicar por qué hace algo, es señal de que no está claro. Mi trabajo es asegurar que entiendan qué construimos y por qué importa, no solo cómo implementarlo.

---

## Rol del TL

El rol del TL no es ser el mejor programador. Es asegurarse de que el equipo avance sin fricción y tome buenas decisiones.

---

## Conclusión

Con equipos junior/mid, el éxito no depende de la arquitectura elegida. Depende de que cada desarrollador entienda las reglas, tenga feedback rápido y pueda mejorar. Eso es lo que intento construir.