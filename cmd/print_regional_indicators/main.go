package main

func main() {
	r := int('\U0001F1E6')

	for i := 0; i < 26; i++ {
		println("\"" + string(rune(r+i)) + "\",")
	}
}
