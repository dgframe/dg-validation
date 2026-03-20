package gookit

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/google/uuid"
	"github.com/gookit/validate"

	dgvalidation "github.com/dgframe/dg-validation"
)

// ErrDatabaseRequired is returned when a database-backed rule is invoked but no database is provided.
var ErrDatabaseRequired = errors.New("dg-http-validation: database connection required for unique/exists rules")

// Adapter implements contracts.Validator using gookit/validate.
// Adapter implements dgvalidation.Validator using gookit/validate.
type Adapter struct {
	db          *sql.DB
	stopOnError bool
	skipOnEmpty bool
}

// ValidationResult wraps gookit validation errors to satisfy dgvalidation.Result.
type ValidationResult struct {
	errors validate.Errors
	valid  bool
}

func (r *ValidationResult) Valid() bool {
	return r.valid
}

func (r *ValidationResult) Violations() []dgvalidation.Violation {
	var violations []dgvalidation.Violation
	for field, errMap := range r.errors {
		for rule, msg := range errMap {
			violations = append(violations, &ValidationViolation{
				FieldPath:   field,
				RuleName:    rule,
				MessageText: msg,
			})
		}
	}
	return violations
}

func (r *ValidationResult) ByField() map[string][]dgvalidation.Violation {
	res := make(map[string][]dgvalidation.Violation)
	for field, errMap := range r.errors {
		for rule, msg := range errMap {
			res[field] = append(res[field], &ValidationViolation{
				FieldPath:   field,
				RuleName:    rule,
				MessageText: msg,
			})
		}
	}
	return res
}

func (r *ValidationResult) First() dgvalidation.Violation {
	for field, errMap := range r.errors {
		for rule, msg := range errMap {
			return &ValidationViolation{
				FieldPath:   field,
				RuleName:    rule,
				MessageText: msg,
			}
		}
	}
	return nil
}

// ValidationViolation implements contracts.Violation.
type ValidationViolation struct {
	FieldPath   string `json:"field"`
	RuleName    string `json:"rule"`
	MessageText string `json:"message"`
	Actual      any    `json:"value,omitempty"`
}

func (v *ValidationViolation) Field() string   { return v.FieldPath }
func (v *ValidationViolation) Rule() string    { return v.RuleName }
func (v *ValidationViolation) Message() string { return v.MessageText }
func (v *ValidationViolation) Value() any      { return v.Actual }

// NewAdapter creates a new gookit validator adapter.
func NewAdapter() *Adapter {
	return &Adapter{
		stopOnError: true, // Default to fail-fast behavior
		skipOnEmpty: true,
	}
}

// WithDatabase sets the database connection for database-aware validations.
func (a *Adapter) WithDatabase(db *sql.DB) *Adapter {
	a.db = db
	return a
}

// StopOnError sets whether validation should stop on the first error.
func (a *Adapter) StopOnError(stop bool) *Adapter {
	a.stopOnError = stop
	return a
}

// SkipOnEmpty sets whether validation should be skipped for empty fields.
func (a *Adapter) SkipOnEmpty(skip bool) *Adapter {
	a.skipOnEmpty = skip
	return a
}

// Validate checks a Subject against its defined rules.
func (a *Adapter) Validate(ctx context.Context, data dgvalidation.Subject, scene ...string) (dgvalidation.Result, error) {
	var v *validate.Validation

	switch d := data.(type) {
	case map[string]interface{}:
		v = validate.Map(d)
	case map[string]string:
		// Convert map[string]string to map[string]interface{} for gookit
		converted := make(map[string]interface{})
		for k, val := range d {
			converted[k] = val
		}
		v = validate.Map(converted)
	default:
		v = validate.Struct(data)
	}

	v.StopOnError = a.stopOnError
	v.SkipOnEmpty = a.skipOnEmpty

	if err := a.registerCustomValidators(v); err != nil {
		return nil, err
	}

	if len(scene) > 0 {
		v.SetScene(scene[0])
	}

	valid := v.Validate()
	return &ValidationResult{
		errors: v.Errors,
		valid:  valid,
	}, nil
}

func (a *Adapter) registerCustomValidators(v *validate.Validation) error {
	// Rule: Database-backed validators MUST fail if DB is missing.

	// unique:table,column,ignore_col,ignore_val
	v.AddValidator("unique", func(val interface{}, args ...string) bool {
		if a.db == nil {
			v.AddError("_internal", "unique", ErrDatabaseRequired.Error())
			return false
		}
		if len(args) < 2 {
			return false
		}
		table, column := args[0], args[1]
		query := fmt.Sprintf("SELECT count(*) FROM %s WHERE %s = ?", table, column)
		qArgs := []interface{}{val}

		if len(args) >= 4 {
			ignoreCol, ignoreVal := args[2], args[3]
			query += fmt.Sprintf(" AND %s != ?", ignoreCol)
			qArgs = append(qArgs, ignoreVal)
		}

		var count int64
		err := a.db.QueryRow(query, qArgs...).Scan(&count)
		return err == nil && count == 0
	})

	// exists:table,column,extra_col,extra_val...
	v.AddValidator("exists", func(val interface{}, args ...string) bool {
		if a.db == nil {
			v.AddError("_internal", "exists", ErrDatabaseRequired.Error())
			return false
		}
		if len(args) < 2 {
			return false
		}
		table, column := args[0], args[1]
		query := fmt.Sprintf("SELECT count(*) FROM %s WHERE %s = ?", table, column)
		qArgs := []interface{}{val}

		// Handle extra filters (pairs of col,val)
		if len(args) > 2 {
			for i := 2; i < len(args); i += 2 {
				if i+1 < len(args) {
					col := args[i]
					val := args[i+1]
					query += fmt.Sprintf(" AND %s = ?", col)
					qArgs = append(qArgs, val)
				}
			}
		}

		var count int64
		err := a.db.QueryRow(query, qArgs...).Scan(&count)
		return err == nil && count > 0
	})

	// unique_multi:table,col1,val1,col2,val2...
	v.AddValidator("unique_multi", func(val interface{}, args ...string) bool {
		if a.db == nil {
			v.AddError("_internal", "unique_multi", ErrDatabaseRequired.Error())
			return false
		}
		if len(args) < 2 {
			return false
		}
		table, column := args[0], args[1]
		query := fmt.Sprintf("SELECT count(*) FROM %s WHERE %s = ?", table, column)
		qArgs := []interface{}{val}

		if len(args) > 2 {
			for i := 2; i < len(args); i += 2 {
				if i+1 < len(args) {
					col := args[i]
					val := args[i+1]
					query += fmt.Sprintf(" AND %s = ?", col)
					qArgs = append(qArgs, val)
				}
			}
		}

		var count int64
		err := a.db.QueryRow(query, qArgs...).Scan(&count)
		return err == nil && count == 0
	})

	// Custom Validators Ported from Legacy

	// uuid
	v.AddValidator("uuid", func(val any) bool {
		s, ok := val.(string)
		if !ok {
			return false
		}
		_, err := uuid.Parse(s)
		return err == nil
	})

	// slug
	v.AddValidator("slug", func(val any) bool {
		slug, ok := val.(string)
		if !ok || slug == "" {
			return false
		}
		slugPattern := `^[a-z0-9]+(?:-[a-z0-9]+)*$`
		matched, _ := regexp.MatchString(slugPattern, slug)
		return matched
	})

	// phone
	v.AddValidator("phone", func(val any) bool {
		phone, ok := val.(string)
		if !ok || phone == "" {
			return false
		}
		cleaned := strings.Map(func(r rune) rune {
			if unicode.IsDigit(r) || r == '+' {
				return r
			}
			return -1
		}, phone)
		length := len(cleaned)
		if strings.HasPrefix(cleaned, "+") {
			return length >= 11 && length <= 16
		}
		return length >= 10 && length <= 15
	})

	// password
	v.AddValidator("password", func(val any) bool {
		password, ok := val.(string)
		if !ok || len(password) < 8 {
			return false
		}
		var hasUpper, hasLower, hasNumber bool
		for _, char := range password {
			switch {
			case unicode.IsUpper(char):
				hasUpper = true
			case unicode.IsLower(char):
				hasLower = true
			case unicode.IsDigit(char):
				hasNumber = true
			}
		}
		return hasUpper && hasLower && hasNumber
	})

	// username
	v.AddValidator("username", func(val any) bool {
		username, ok := val.(string)
		if !ok || len(username) < 3 || len(username) > 20 {
			return false
		}
		// Alphanumeric + underscore + hyphen
		pattern := `^[a-zA-Z0-9_-]+$`
		matched, _ := regexp.MatchString(pattern, username)
		return matched
	})

	// alpha_space
	v.AddValidator("alpha_space", func(val any) bool {
		value, ok := val.(string)
		if !ok {
			return false
		}
		for _, char := range value {
			if !unicode.IsLetter(char) && !unicode.IsSpace(char) {
				return false
			}
		}
		return true
	})

	// no_sql
	v.AddValidator("no_sql", func(val any) bool {
		value, ok := val.(string)
		if !ok {
			return false
		}
		value = strings.ToLower(value)
		sqlKeywords := []string{
			"select", "insert", "update", "delete", "drop", "create",
			"alter", "exec", "execute", "union", "declare", "--", "/*", "*/",
			"xp_", "sp_", "0x", "char(", "nchar(", "varchar(", "nvarchar(",
		}
		for _, keyword := range sqlKeywords {
			if strings.Contains(value, keyword) {
				return false
			}
		}
		return true
	})

	// no_xss
	v.AddValidator("no_xss", func(val any) bool {
		value, ok := val.(string)
		if !ok {
			return false
		}
		value = strings.ToLower(value)
		xssPatterns := []string{
			"<script", "</script", "javascript:", "onerror=", "onload=",
			"onclick=", "onmouseover=", "<iframe", "<object", "<embed",
			"eval(", "expression(", "vbscript:", "data:text/html",
		}
		for _, pattern := range xssPatterns {
			if strings.Contains(value, pattern) {
				return false
			}
		}
		return true
	})

	// color_hex
	v.AddValidator("color_hex", func(val any) bool {
		color, ok := val.(string)
		if !ok {
			return false
		}
		hexPattern := `^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$`
		matched, _ := regexp.MatchString(hexPattern, color)
		return matched
	})

	// timezone
	v.AddValidator("timezone", func(val any) bool {
		tz, ok := val.(string)
		if !ok {
			return false
		}
		if tz == "UTC" || tz == "GMT" {
			return true
		}
		// Basic format validation: Area/Location
		tzPattern := `^[A-Z][a-z]+/[A-Z][a-z_]+$`
		matched, _ := regexp.MatchString(tzPattern, tz)
		return matched
	})

	return nil
}
