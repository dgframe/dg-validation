# dg-http-validation Configuration

This document outlines configuration options for the **dg-http-validation** plugin.  
As a sovereign contract plugin, it does not require execution or external dependencies, but implementation adapters (e.g., gookit) may provide configurable behaviors.

---

## 1. Default Adapter Configuration

- **Adapter:** gookit (default, optional)
- **Purpose:** Provides a ready-to-use validation engine without breaking the contract.
- **Notes:** Users may replace the default adapter with a custom implementation via dependency injection.

### Example (conceptual)

```yaml
validation:
  adapter: gookit
  options:
    strict_mode: true
    message_format: json
```

> The above configuration is optional. All adapters must respect the dg-http-validation contract.

---

## 2. Adapter Options

- `strict_mode` – boolean, determines whether all fields must be validated.
- `message_format` – string, defines output format of validation results (e.g., `json`, `text`).

> These options are adapter-specific; the plugin itself enforces no behavior.

---

## 3. Environment Variables

- `DG_HTTP_VALIDATION_ADAPTER` – override adapter at runtime.
- `DG_HTTP_VALIDATION_STRICT` – enable/disable strict mode.

> Environment variables provide runtime flexibility without breaking the contract.

---

## 4. Future Extensions

- Custom adapters can introduce additional options.  
- All new options must **never violate the sovereign contract**.

---

## Notes

- dg-http-validation **does not execute validation logic**.  
- Configurations only influence **adapter behavior**, not the contract itself.
