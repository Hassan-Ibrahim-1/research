# TODO

create a prompt area using a viewport. 
    - there should be a character limit with an indicater
    -- alt+enter to enter text, produces a message that the main app handles
        something like PromptEnteredMsg(string)
    -- on enter remove the text
    - if paging through the chat press a key (enter) to gain focus on the prompt area
      press another key (escape) to remove focus from the prompt area
    - cycle through previous prompts with the arrow keys
    -- support for whitespace, newlines (enter) and backspace
    - move only words to next line not individual characters
    - test prompt.go
    - scrolling
    - height changes based on text length and newlines
