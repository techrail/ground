package netserver

import "context"

// GetValueFromCtx retrieves a value from context.
// K is restricted to 'comparable' types (int, string, structs, pointers, etc.).
// V is the type of the value you want to retrieve.
func GetValueFromCtx[K comparable, V any](ctx context.Context, key K, defaultValue V) V {
	val := ctx.Value(key)

	// Type assertion: check if val is not nil AND is of type V
	if castedVal, ok := val.(V); ok {
		return castedVal
	}

	return defaultValue
}

// GetValueFromCtxByStringKey function forces the 'key' argument to be a string.
// It still uses 'V any' so you can retrieve an Int, String, or Struct using that string key.
func GetValueFromCtxByStringKey[V any](ctx context.Context, key string, defaultValue V) V {
	// We call the base function, letting the compiler know K is 'string'
	return GetValueFromCtx[string, V](ctx, key, defaultValue)
}

// File ends here
