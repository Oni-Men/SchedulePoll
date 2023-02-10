package emoji

const (
	ProgressFG = "◻️"
	ProgressBG = "◼️"
	Calendar   = "🗓️"
)

var ABCs = [26]rune{}

func init() {
	r := int('\U0001F1E6')

	for i := 0; i < 26; i++ {
		ABCs[i] = rune(r + i)
	}
}

func ABCDEmoji(i int) string {
	return string([]rune{ABCs[i%26]})
}
