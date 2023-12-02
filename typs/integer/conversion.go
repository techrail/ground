package integer

func Base10ToBase36(base10Num int64) string {
	const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// Handle the case where the input number is zero
	if base10Num == 0 {
		return "0"
	}

	// Compute the base36 string representation of the input number
	var base36Str string
	for base10Num > 0 {
		rem := base10Num % 36
		base36Str = string(charset[rem]) + base36Str
		base10Num /= 36
	}

	return base36Str
}
