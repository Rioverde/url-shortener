package codegen

import "errors"

// ErrKeyTooShort is returned by the code generator when the requested key
// length is below minKeyLength. Concrete error messages wrap this sentinel
// and include the offending value, so the inner cause is still detectable
// via errors.Is(err, ErrKeyTooShort).
var ErrKeyTooShort = errors.New("key length too short")
