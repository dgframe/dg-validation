

# dg-validation Middleware Policy

This document outlines **middleware policy expectations** for adapters implementing dg-validation.  
The plugin itself does not implement middleware; it defines the governance rules for middleware behavior in transport layers.

---

## 1. Purpose

- Ensure middleware behavior aligns with **sovereign contract principles**.
- Maintain **kernel-plugin call direction** integrity.
- Prevent adapters from violating contract boundaries in middleware.

---

## 2. Scope

- Applies to all HTTP or transport adapters that wrap validation as middleware.
- Covers pre-processing, validation execution, and result propagation.
- Does **not** include core plugin logic, which remains contract-only.

---

## 3. Kernel/Plugin Call Policy

- Middleware **must call the plugin via the Validator interface** only.
- Middleware **must not** access kernel internals directly.
- Middleware **must not** perform logging, database operations, or global state mutation outside adapter scope.

---

## 4. Adapter Guidelines

- Middleware should be **stateless or deterministic**.
- Must respect **fail-fast and strict mode configurations** from adapter options.
- Should propagate errors in a **standardized Result/Violation format**.
- Middleware should be pluggable and replaceable without altering plugin contracts.

---

## 5. Notes

- Middleware is optional; adapters may provide transport-specific helpers.
- Adapters must follow these policies to maintain **SOVEREIGN + CERTIFIED** status.