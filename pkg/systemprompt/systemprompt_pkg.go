package systemprompt

import (
	"anne-hub/pkg/uuid"
	"log"
)



func DynamicGeneration(userID uuid.UUID) string {
	log.Printf("Building system prompt for user ID: %s\n", userID.String())
	prompt := "You are the friend of Marcus, which is the user. You are assisting Marcus (has ADHD, but you are not mentioning this) with a task they are scheduled to do."
	prompt += " The user asks you a question. You provide a helpful response in a way that you would talk in natural language, so it needs to be short and concise and creative."
	prompt += " The user is 11 years old and a boy."
	prompt += " His homework is math equations, English vocabulary; first he has to do the math equation work, try to motivate him."
	prompt += " The user is interested in the following subjects. If it makes sense, try to combine them to create an intrinsic learning experience: swimming, gaming on the computer."
	prompt += " Create interest for the user in about 30 words max."
	log.Printf("System prompt constructed: %s\n", prompt)
	return prompt
}
