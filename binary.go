package main

// isBinary checks if data looks like a binary file by scanning for null bytes
// in the first 8000 bytes (same heuristic git uses).
func isBinary(data []byte) bool {
	if len(data) == 0 {
		return false
	}
	checkLen := 8000
	if len(data) < checkLen {
		checkLen = len(data)
	}
	for i := 0; i < checkLen; i++ {
		if data[i] == 0 {
			return true
		}
	}
	return false
}
