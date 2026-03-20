# dg-validation Compliance Declaration

**Target Kernel:** dg-core >= v1.2  
**Plugin:** dg-validation  
**Review Type:** Compliance Gauntlet

---

## 1. Kernel Authority Acknowledgement
- dg-core is the sole lifecycle authority  
- dg-validation is a subordinate contract plugin  
- All call direction flows **Kernel → Plugin**  
- Default adapter (e.g., gookit) is optional and does not change compliance

✅ **Compliant**

---

## 2. Scope Compliance
| Rule | Status |
|-----|-------|
| Defines validation semantics only | ✅ |
| Execution-agnostic | ✅ |
| Stateless contract surface | ✅ |
| Adapters implement execution outside plugin | ✅ |
| Default adapter optional and replaceable | ✅ |

---

## 3. Forbidden Coupling Verification
The following dependencies are explicitly absent:  
- go-playground/validator (in plugin)  
- Gookit/validate (in plugin)  
- Struct tags (plugin does not inspect)  
- Logging frameworks  
- Global state mutation  

> Adapters may use external libs, but must not violate plugin contract

✅ **Compliant**

---

## 4. Implementation / Adapter Compliance
- All adapters must implement `Validator` interface  
- Must respect fail-fast and strict mode options  
- DB-dependent rules or other external dependencies must remain adapter-scoped  
- Any violation invalidates plugin compliance

---

## 5. Future Auditability
- Compliance reviews should be performed periodically  
- Any new adapter or plugin extension must be audited before release

---

## 6. Final Verdict
**Compliance Status:** 🟩 **COMPLIANT**

This plugin is **admitted** into the compliant kernel ecosystem.

---

**Reviewed By:** Kernel Authority  
**Last Reviewed:** 2026-01-21
**Status:** Certified Sovereign
