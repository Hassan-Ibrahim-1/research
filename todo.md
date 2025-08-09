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
    - being able to attach urls / files
        parse the user message for certain patterns like #(url) / @(file)
        or tell the llm to tell the system to get a file and send the file back to it

        search for the @ symbol and then read until the start of a bracket
        then read until a right paren is found
        be able to escape this by adding a backslash like \@attack-file()
        if the function is not valid or parenthesis are missing
        then just ignore it and send it to the llm normally

        otherwise, remove the function and swap it out with the file contents
        if the file or url are not found display an error in the chat thing
        and don't send the prompt to the llm
        
        parse it like a function call, can have multiple arguments
        separated by a comma

        @attach-file()
        @attack-link()

        (maybe also add the ability to add files that are not part of a .gitignore)

### Tomorrow
    - add commands to llm/
    - embed stuff into prompts (these should be stored in messages)
        the ui will only get the messages with the function calls
        not the actual embedded stuff
    add some styling to differentiate user prompts and responses
        background stuff, etc
    send embedded messages to the llm

    command table in llm/
