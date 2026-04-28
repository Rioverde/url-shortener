package domain

// CodeGenerator generates random short codes used as URL keys.
// The interface lives in this (consumer) package so URLService is not
// coupled to a specific generation strategy. Concrete implementations
// live elsewhere (see internal/lib/codegen).
type CodeGenerator interface {
	GenerateRandomString(n int) (string, error)
}
