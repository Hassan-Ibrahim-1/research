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
    -- streaming llm responses
    -- being able to attach urls / files

### Up next
    -- add commands to llm/
    -- embed stuff into prompts (these should be stored in messages)
        -- the ui will only get the messages with the function calls
        -- not the actual embedded stuff
    add some styling to differentiate user prompts and responses
        background stuff, etc
    -- send embedded messages to the llm

    tests

    -- command table in llm/
    when at the bottom, scroll automatically

    figure out a way to let my laptop communicate with my pc
    so that i can use it to run better models

    status indicators for reading a file / fetching a link
