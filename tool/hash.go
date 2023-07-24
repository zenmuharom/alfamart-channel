package tool

// Check if a string is an MD5 hash
func IsMD5Hash(hash string) bool {
	return len(hash) == 32
}

// Check if a string is a SHA256 hash
func IsSHA256Hash(hash string) bool {
	return len(hash) == 64
}
