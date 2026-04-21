// Package compute implements arithmetic operations over log entry fields.
//
// Each rule specifies a destination field and an expression of the form:
//
//	"dest=left<op>right"
//
// where <op> is one of +, -, *, /. Operands must be numeric fields present
// in the entry. Entries missing an operand or causing a divide-by-zero are
// passed through unmodified.
package compute
