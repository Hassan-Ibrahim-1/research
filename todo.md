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

    -- better distinction between user messages and llm responses
    - streaming llm responses
        don't format the viewport until the response is fully streamed
        don't run glamour.Render() on every redrawViewport, its just wasteful
        only call it on new user messages and, when on llm responses on when they are fully sent
    - being able to attach urls / files
        parse the user message for certain patterns like #(url) / @(file)
        or tell the llm to tell the system to get a file and send the file back to it
        (maybe also add the ability to add files that are not part of a .gitignore)
