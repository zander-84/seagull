package tool

func SliceUnique[T comparable](in []T) (out []T) {
	out = make([]T, 0)
	for i := 0; i < len(in); i++ {
		if !InSlice(in[i], out) {
			out = append(out, in[i])
		}
	}
	return
}

func InSlice[T comparable](elem T, slice []T) bool {
	for _, v := range slice {
		if elem == v {
			return true
		}
	}
	return false
}

func SliceChunk[T comparable](in []T, size int) [][]T {
	var ret [][]T
	for size < len(in) {
		ret = append(ret, in[:size])
		in = in[size:]
	}
	ret = append(ret, in)
	return ret
}

func SliceChunkAny[T comparable](in []T, size int) [][]any {
	if size <= 0 {
		size = 1
	}
	var chunks [][]any
	for i := 0; i < len(in); i += size {
		end := i + size
		if end > len(in) {
			end = len(in)
		}
		chunk := make([]any, end-i)
		for j, k := i, 0; j < end; j, k = j+1, k+1 {
			chunk[k] = in[j]
		}
		chunks = append(chunks, chunk)
	}
	return chunks
}
