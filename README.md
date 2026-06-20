# AutoConfine

[![Go Report Card](https://goreportcard.com/badge/github.com/ValentinTorassa/autoconfine)](https://goreportcard.com/report/github.com/ValentinTorassa/autoconfine)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

**AutoConfine** genera perfiles de **mínimos privilegios** para contenedores [OCI](https://opencontainers.org/) mediante observación dinámica del kernel con **eBPF**.

> Los contenedores se ejecutan por defecto con perfiles seccomp que permiten ~300 syscalls. Una aplicación real usa típicamente entre 40 y 80. AutoConfine cierra la brecha automáticamente.

## Modos de operación

| Modo | Comando | Función |
|---|---|---|
| Aprender | `autoconfine learn` | Observa un contenedor en ejecución y registra cada syscall invocada. |
| Generar | `autoconfine generate` | Convierte la traza en un perfil seccomp JSON compatible con OCI. |
| Aplicar | `autoconfine enforce` | Ejecuta el contenedor con Podman aplicando el perfil derivado. |
| Auditar | `autoconfine enforce --audit` | Alerta cuando aparece una syscall fuera del perfil aprendido (drift detection). |

## Instalación

```bash
go install github.com/ValentinTorassa/autoconfine/cmd/autoconfine@latest
```

Requisitos de runtime:

- Linux con kernel >= 5.8 (para eBPF CO-RE).
- Cabeceras del kernel o BTF disponible.
- Podman >= 4.0.
- Permisos suficientes para cargar programas eBPF (`CAP_BPF`, `CAP_PERFMON` o root).

## Uso rápido

```bash
# 1. Fase de aprendizaje
autoconfine learn --image nginx --duration 5m --out nginx.trace.jsonl

# 2. Generación del perfil
autoconfine generate nginx.trace.jsonl --out nginx-seccomp.json

# 3. Ejecución protegida
autoconfine enforce --profile nginx-seccomp.json -- podman run --rm nginx
```

## Arquitectura

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   learn     │────▶│  generate   │────▶│   enforce   │
│   (eBPF)    │     │ (seccomp)   │     │  (Podman)   │
└─────────────┘     └─────────────┘     └─────────────┘
                                              │
                                              ▼
                                        ┌─────────────┐
                                        │ drift audit │
                                        └─────────────┘
```

## Licencia

Apache 2.0 — ver [LICENSE](LICENSE).

---

Trabajo presentado al **Premio CAI Pre-Ingeniería 2026** por **Valentín Torassa Colombero**.
