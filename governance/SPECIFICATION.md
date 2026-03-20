# dg-http-validation Specification

**Target Kernel:** dg-core >= v1.2

---

## 1. Purpose

dg-http-validation defines the **authoritative validation semantics** for the dg ecosystem.

It is a **contract-only plugin** whose sole responsibility is to describe:
- What validation *means*
- How validation rules are identified
- How validation failures are expressed and propagated

The dg-core kernel is authoritative. dg-http-validation exists to serve the kernel.

---

## 2. Scope

### In Scope

The dg-http-validation plugin MAY define:
- Validation concepts and terminology
- Rule identity and naming semantics
- Field and scope addressing models
- Validation result and violation semantics

### The Subject Doctrine
Validation targets are represented as an opaque **Subject**.
- dg-http-validation handles the *output* of validation (Results/Violations).
- dg-http-validation does *not* handle the *input* structure.
- Interpreting the Subject (via reflection, tags, or schema) is the **exclusive right of the Adapter**.

### Explicitly Out of Scope

The dg-http-validation plugin MUST NOT:
- Execute validation logic
- Inspect structs, DTOs, or schemas
- Use reflection or tags
- Import or reference any validation engine
- Bind validation to transports (HTTP, gRPC, CLI)

Execution belongs exclusively to adapters.

---

## 3. Kernel Authority Model

The dg-core kernel:
- Owns lifecycle orchestration
- Initiates validation requests
- Receives validation results

The dg-http-validation plugin:
- Declares **validation contracts only**; no execution is performed.
- Cannot assume execution context.
- Cannot call the kernel.
- Default adapter (e.g., gookit) is provided optionally and can be replaced by a custom adapter.

Call direction is strictly **Kernel → Plugin**.

---

## 4. Forbidden Couplings (Hard Prohibitions)

dg-http-validation MUST NOT reference:
- go-playground/validator
- Gookit/validate
- Struct tags or annotations
- HTTP status codes
- Logging frameworks

Violations invalidate compliance.

---

**Status:** Final — Certified Sovereign v1.1.0
