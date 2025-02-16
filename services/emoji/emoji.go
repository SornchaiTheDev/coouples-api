package emoji

import "math/rand"

// Known emoji ranges in Unicode
var emojiRanges = []struct {
	start rune
	end   rune
}{
	{0x1F300, 0x1F5FF}, // Misc Symbols & Pictographs
	{0x1F600, 0x1F64F}, // Emoticons
	{0x1F680, 0x1F6FF}, // Transport & Map Symbols
	{0x1F900, 0x1F9FF}, // Supplemental Symbols & Pictographs
	{0x2600, 0x26FF},   // Misc Symbols
	{0x2700, 0x27BF},   // Dingbats
}

// Function to generate a random seed emoji
func generateSeedEmoji() rune {
	rangeIndex := rand.Intn(len(emojiRanges))
	selectedRange := emojiRanges[rangeIndex]
	return rune(rand.Intn(int(selectedRange.end-selectedRange.start+1)) + int(selectedRange.start))
}

// Function to check if a rune is within valid emoji ranges
func isEmoji(r rune) bool {
	for _, rRange := range emojiRanges {
		if r >= rRange.start && r <= rRange.end {
			return true
		}
	}
	return false
}

// Generate emoji set from a seed emoji by iterating forward
func generateEmojiSet(seed rune, count int) []rune {
	emojis := []rune{}
	current := seed

	for len(emojis) < count {
		if isEmoji(current) {
			emojis = append(emojis, current)
		}
		current++ // Move to next Unicode point
	}
	return emojis
}

func GetSet(amount int) []rune {
	seed := generateSeedEmoji()
	return generateEmojiSet(seed, amount)
}
