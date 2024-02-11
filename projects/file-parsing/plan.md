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

