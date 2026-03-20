# dg-http-validation

**The Authoritative Validation Authority for the dg Ecosystem**

`dg-http-validation` is a **Sovereign Semantic Contract** under the **dg-core Kernel Authority Model**. It dictates the formal semantics, interface protocols, and behavioral invariants for all validation activities within the `dg` ecosystem. 

This package is **Absolute Zero**: it contains no execution logic, no external dependencies, and is entirely transport-blind.

---

## 🏛️ The Doctrine of Absolute Zero

Under the **Kernel Authority Model (dg-core >= v1.2)**, this package represents the "Clinical Contract" layer:

1.  **Zero Execution**: This package never executes a validation rule or inspects a structure.
2.  **Zero Engines**: It is forbidden from importing or knowing about `gookit/validate`, `validator/v10`, or any other engine.
3.  **Zero Transport**: It has no knowledge of HTTP, gRPC, JSON, or any framing protocol.
4.  **Zero State**: It is strictly stateless and owns no background processes or routines.

---

## 🛡️ The Subject Doctrine

To prevent "introspection invitation" (the tendency to design contracts that assume struct reflection), `dg-http-validation` enforces the **Subject Doctrine**:

-   **Opaque Input**: All data to be validated is passed as an opaque **`Subject`**.
-   **Output-Oriented**: The contract layer is authoritative over the **Result** and **Violation** semantics, but agnostic to the input format.
-   **Adapter Sovereignty**: The interpretation of a `Subject` (e.g., via struct tags, reflection, or schema) is the **exclusive domain of the Adapter** (Application Layer).

---

## 🛠️ Authoritative Usage

### 1. Declare Intent (Domain Layer)
Domain services must only depend on the `Validator` contract:

```go
import "github.com/dgframe/dg-http-validation"

type UserService struct {
    // Depend on the AUTHORITY, not the implementation
    validator dgvalidation.Validator
}

func (s *UserService) CreateUser(ctx context.Context, req CreateUserRequest) error {
    // Pass the request as an opaque Subject
    res, err := s.validator.Validate(ctx, req)
    if err != nil {
        return err
    }
    
    if !res.Valid() {
        // Handle results mapped to dg-http-validation semantics
        return &MyDomainError{Violations: res.Violations()}
    }
}
```

### 2. Implement Execution (Adapter Layer)
Adapters translate the opaque `Subject` into engine-specific logic:

```go
// In the Skeleton / Application Layer
func (a *GokitAdapter) Validate(ctx context.Context, data interface{}, scene ...string) (dgvalidation.Result, error) {
    // data is treated as an opaque Subject here.
    // 1. Run Engine (Gokit, etc.)
    // 2. Map engine errors to dgvalidation.Violation
    // 3. Return dgvalidation.Result
}
```

---

## 📜 Governance Manifest

The authoritative behavior of this contract is governed by:
-   **[SPECIFICATION.md](governance/SPECIFICATION.md)**: Semantic rules and classifications.
-   **[INJECTION_CONTRACT.md](governance/INJECTION_CONTRACT.md)**: Formal kernel-level injection protocols.
-   **[COMPLIANCE.md](governance/COMPLIANCE.md)**: Governance-verified doctrine alignment.

---

## 🏁 Compliance Verdict

This component is a **Hardened Sovereign Contract**.

- [x] Absolute Zero Implementation
- [x] Subject Doctrine Enforcement
- [x] Kernel Authority v1.2.0 Alignment
- [x] Isomorphic Contract-to-Code Mapping
