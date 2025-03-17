package chatgpt

import "sync"

type SystemPrompt struct {
	mu      sync.Mutex
	prompts []string
}

var syspro SystemPrompt = SystemPrompt{
	prompts: make([]string, 0),
}

func GetSysPrompts() []string {
	syspro.mu.Lock()
	defer syspro.mu.Unlock()
	answer := make([]string, len(syspro.prompts))
	copy(answer, syspro.prompts)
	return answer
}

func AddSysPrompt(newsysprompt string) {
	syspro.mu.Lock()
	defer syspro.mu.Unlock()
	syspro.prompts = append(syspro.prompts, newsysprompt)
}

func DeleteAllSysPrompt() {
	syspro.mu.Lock()
	defer syspro.mu.Unlock()
	syspro.prompts = syspro.prompts[:0]
}
