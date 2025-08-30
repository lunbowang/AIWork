package chatinternal

const (
	BASE_PROMPAT_TEMPLATE = `Current conversation:
{{.history}}\n\n`

	_defaultMrklPrefix = `Today is {{.today}}.
Answer the following questions as best you can. You have access to the following tools:

{{.tool_descriptions}}

IMPORTANT: You MUST complete the entire format including Action Input, Observation, and Final Answer.` + BASE_PROMPAT_TEMPLATE

	OUT_PROMPT_TEMPLATE = `<< instructions >>
- Your response should follow the JSON format.
- Your response should have the following structure: {"chatType": {{.chatType}}, "data": {{.data}} }
- "chatType" this is a fixed output`

	_defaultChatLogPrompts = `Please summarize based on the following chat conversations
- Role: You are an office assistant who helps employees summarize communication matters at work. Matters mainly include things such as "tasks to be done, leave for leave"
- work
1. You must first read the chat content to understand the relationship between employees and distinguish between superiors and subordinates
2. You need to distinguish what matters are based on the chat content, such as [tasks to be done, leave for leave]
4. You need to first obtain the overall context of the matter, and then summarize it based on the context, such as: the person initiated the leave, whether there are any follow-up, etc.
5. Personnel information needs to be fully output
6. keep Chinese output

- chatlog
{{.input}}

- require
The output content needs to be output in the following format and can be parsed by json
[
    {
       "type": int,      // task type; enum : 1. task to be done, 2 approval
       "title": string,  // title
       "content": string // content
    }, {
       "type": int,      // task typ; enum : 1. task to be done, 2 approval
       "title": string,  // title
       "content": string // content
    }
]
`
)
