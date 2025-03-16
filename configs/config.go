package configs

var (
	/*
		Temperature: Behavior => Use Case
		0.0: Deterministic, focused, repetitive => Factual answers, when precision is key
		0.3: More accurate, fewer creative jumps => Standard Q&A, clear explanations
		0.7: Balanced, some creativity and variation => Conversational, general chat, casual writing
		1.0: Very creative, diverse, but less predictable => Brainstorming, storytelling, creative writing
		>1.0: Highly random, often incoherent => Rarely used unless you need chaotic output
	*/
	DefaultTemperature float32 = 0.7

	// Response limit. Cuts long responses.
	DefaultMaxTokens int = 1000
)
