# dg-validation Injection Contract

**Target Kernel:** dg-core >= v1.2

---

## 1. Purpose

This document defines the **formal injection contract** between the dg-core kernel and the dg-validation plugin.

The contract specifies:
- Call direction and authority
- Required abstractions
- Ownership boundaries

> Note: A default adapter (e.g., gookit) is provided optionally for convenience.
> Users may replace it with a custom adapter without violating the sovereign contract.

---

## 2. Authority and Call Direction

- The dg-core kernel is authoritative.
- dg-validation is subordinate.
- Call direction is strictly **Kernel → dg-validation**.

---

## 3. Ownership Boundaries

### Kernel Owns
- Validation lifecycle orchestration
- Failure reaction strategy
- Transport mapping (HTTP, gRPC, CLI)
- Error rendering and localization

### dg-validation Owns
- Validation semantics
- Rule identity model
- **Validator**, **Result**, and **Violation** contracts
- Determinism guarantees

### Skeleton Owns
- Concrete validation engines (e.g., gookit adapter)
- Rule execution logic
- Mapping engine output to contracts

> Adapters are optional and replaceable; the plugin itself enforces the contract only.

---

## 4. Abstractions

### 4.1 Validator
The primary injection surface. Accepts a **Subject** and returns a **Result**.

### 4.2 Result
Represent an immutable outcome of a validation pass.

### 4.3 Violation
Represents a single rule failure with logical field path.

---

## 5. Prohibited Behaviors

The dg-validation plugin and its adapters MUST NOT:
- Depend on external validator libraries directly
- Assume struct-tag-based validation
- Return HTTP status codes
- Access global state

Violations invalidate compliance.

---

**Status:** Sovereign Contract v1.2.0 (Isomorphic)
