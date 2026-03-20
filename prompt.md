

---
# EXECUTION PROMPT — dg-http-validation Sovereign Refactor

## ROLE

You are a **Sovereign Plugin Refactor Agent** operating under **dg-core governance rules**.

Your job is to **normalize, consolidate, and certify** the plugin  
`dg-http-validation` without introducing new behavior.

You MUST behave deterministically.  
You MUST NOT invent features, APIs, or semantics.

---

## OBJECTIVE

Eliminate **contract drift** and **redundant interface definitions** by consolidating all validation contracts into a **single canonical surface**, while preserving:

- Kernel → Plugin authority
- Sovereign contract-only semantics
- Profile **A2 (Hard)** failure behavior
- Adapter optionality
- Governance and certification integrity

---

## HARD CONSTRAINTS (NON-NEGOTIABLE)

❌ Do NOT add new features  
❌ Do NOT weaken failure semantics  
❌ Do NOT introduce logging, HTTP logic, or execution logic  
❌ Do NOT couple adapters to kernel or HTTP  
❌ Do NOT duplicate contracts  
❌ Do NOT change public behavior  

✅ One canonical contract location  
✅ Deterministic behavior only  
✅ Adapter-scoped external dependencies  
✅ Certification must pass  

---

## PHASE 1 — CONTRACT CANONICALIZATION

### Goal
There must be **exactly one authoritative contract definition**.

### Actions

1. **KEEP**
```
dg-http-validation/contracts/validator.go
```

2. **DELETE**
```
dg-http-validation/contracts.go
```

3. **MOVE / CONSOLIDATE** the following interfaces into  
`contracts/validator.go`:

```go
type Subject interface{}

type Validator interface {
    Validate(ctx context.Context, subject Subject, scene ...string) (Result, error)
}

type Result interface {
    Valid() bool
    Violations() []Violation
    ByField() map[string][]Violation
    First() Violation
}

type Violation interface {
    Field() string
    Rule() string
    Message() string
    Value() any
}
```

4. Remove all legacy error-only validation interfaces.

5. Ensure **all adapters, providers, and implementations import ONLY from `contracts/`**.

---

## PHASE 2 — NO-OP FAIL-FAST IMPLEMENTATION

### Goal
Preserve **Profile A2 (Hard)** semantics.

### Actions

Modify:
```
contracts/noop.go
```

- `NoopValidator.Validate(...)` MUST:
  - Return `(Result, ErrCapabilityNotProvided)`
- Implement `NoopResult`:

```go
type NoopResult struct{}

func (NoopResult) Valid() bool { return false }
func (NoopResult) Violations() []Violation { return nil }
func (NoopResult) ByField() map[string][]Violation { return nil }
func (NoopResult) First() Violation { return nil }
```

### Forbidden

❌ Silent success  
❌ Panics  
❌ Execution logic  

---

## PHASE 3 — GOOKIT ADAPTER ALIGNMENT

### Goal
Ensure production adapter implements **rich sovereign contracts**.

### Actions

Modify:
```
adapters/gookit/adapter.go
```

1. Update `Validate` signature:
```go
Validate(ctx context.Context, subject contracts.Subject, scene ...string) (contracts.Result, error)
```

2. Ensure:
- `ValidationResult` implements `contracts.Result`
- `ValidationViolation` implements `contracts.Violation`

3. Adapter rules:
- External libraries remain adapter-scoped
- No logging
- No HTTP codes
- No kernel imports
- Deterministic mapping only

---

## PHASE 4 — PROVIDER ALIGNMENT

### Goal
Ensure DI remains canonical and adapter-agnostic.

Modify:
```
provider.go
```

- Import ONLY from `contracts/`
- Default binding remains `NoopValidator`
- Adapter override remains optional
- No new logic allowed

---

## PHASE 5 — VERIFICATION

### Automated Tests

Run:
```bash
go test ./...
```

Add:
```
certification_test.go
```

Tests MUST verify:
- No adapter → `ErrCapabilityNotProvided`
- `NoopResult.Valid()` is `false`
- Gookit adapter populates `Result` and `Violations`
- No forbidden imports in core plugin

---

## PHASE 6 — GOVERNANCE SYNC

Update ONLY if wording drift exists:
- `SPECIFICATION.md`
- `COMPLIANCE.md`
- `CERTIFICATION.md`

Generate:
```
CERTIFICATION_REPORT.md
```

Must state:
> Canonical Result-based validation contracts enforced.

---

## FINAL ACCEPTANCE CRITERIA

You may stop ONLY if ALL are true:

- Root `contracts.go` removed
- Single canonical contract exists
- All code compiles
- All tests pass
- Fail-fast behavior preserved
- Governance docs aligned
- Plugin is certifiable

---

## OUTPUT REQUIREMENTS

At completion, produce:

1. **Change summary**
2. **Files modified / deleted**
3. **Governance compliance statement**
4. **Certification readiness confirmation**

No commentary. No suggestions. No creativity.

---

### END OF EXECUTION PROMPT
---