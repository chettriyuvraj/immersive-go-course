# Plan


## JSON Data

- Quickly read up the official Go post "Intro to JSON". Marshal, Unmarshal, Encode and Decode make sense now.
- Decoder: Accepts an io.Reader. On triggering decode with a struct argument(empty struct ig), Reads it's contents into a []byte internally, then Marshals this []byte into the struct
- Encoder: Accepts an io.Writer. On triggering encode with a struct argument containing marshalled data, unmarshals it to []byte and writes the []byte to the writer using it's write method.
- Let's try thinking about tests before writing the code
- What would I like to test? 
    - Essentially that ALL records have been parsed + all fields
    - I would frame my test as a single test, not TDD
    - First a function to parse the
    - Essentially a decoder, accepts any io.Reader, on using decode - decodes to a struct of specified type (would be an array of the struct in this case), then checks the fields + count
    - Also a function for highest lowest scores
- Observations after writing code + tests:
    - Fine grained tests to two tests
    - One simply accepts a reader and decodes it to json
    - Other accepts a file and does the same i.e. reuses the first function
    - Seems neat enough


- Plans and observations for each written in their respective folders


### Observations post implementation

- The overall strategy for each one seemed quite similar - fine grained functions to parse, read next record, and cast the raw data into a particular struct
- Is the folder struct with separate test cases sensible ? Check reference implementation
    - Could have been same package with separate code and test files for each (?), would have led to many files in a single folder though
- The PlayerDataStruct is repeated in each folder, can isolate it to a separate struct - check reference impl
- Apart from the custom-binary, other implementations seem pretty neat (?) - compare it to reference impl


### Discussion

Encoding schemes:

- Using Terminating Character:
    - Implementation becomes simple
    - Space conservation in terms of we use exactly what we use
    - Space waste in terms of we always use an extra 'x' bytes, might add up in case of large amount of data
    - We can't use terminating character
- Using fixed length:
    - Implementation becomes simple
    - Space wastage if small data and large space
    - Inflexible, we are stuck with it come what may - more data less data doesn't matter
- Variable length encoding
    - Implementation becomes a tad more complex with many checks depedning on which encoding
    - Space conservation is high
    - Space wastage in terms of how many bytes are used to encode info about the size of data (may be higher or lower than using a terminating character)
    - We can use any character - no bar!


### Comparison with reference implementation

- The reference uses a Parser interface with a Parse() function, this is very neat because it simplifies function names and the overall structure of all the code a lot. I had to write separate func names (often different) and slightly different structures for each file
- I have considered each as a separate parser in itself, while using the parser interface also allows us to unify things such as opening the file only once instead of in all packages
- I think I have written more fine grained tests for helper functions as well (good thing imo)
- JSON:
    - Core func same
- Repeated JSON:
    - Core func same
    - I have fine grained functionalities, but is this too much abstraction which is making my code difficult to follow - reference implementation is pretty straightforward to read. Imo this is a trade off to consider, do you really want to fine grain to the level of 'IsLineAComment()' function. Removing things like this instantly remove one chunk of test + code = lot less for people to read and understand
- CSV:
    - Here the parse function is a little bit larger, the approach of fine grained tests would improve readability here imo
    - Reference impl also has tests
    - I haven't implemented it entirely
- Custom-Binary:
    - Use predefined values such as binary.LittleEndian and binary.BigEndian instead of what I have done (true/false)
    - My over-abstraction has forced me to use an additional struct i.e. more for code-reader to read and process

- Overall thoughts:
    - I have veered towards over-abstraction in this exercise
    - Try to think when making interfaces will make sense i.e. here using something like the Parser interface was a very natural fit
    - The concept of interfaces in terms of readers writers decoders encoders et al has been well understood by me
    - The actual logic in terms of parsing was almost always the same in my impl and the reference impl, structuring of the code was different at times
    

### Using tools instead of code

- Practice using jq csvq etc in the shell