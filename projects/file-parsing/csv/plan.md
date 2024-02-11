# Plan

- I think we can probably use encoding/csv for parsing data
- Get an encoding.Reader from an io.Reader
- Use Read() function to read one by one - gives a string, then split using commas and assign each struct individually (we can also get all the data together)
- Testing strat:
    - Test that Read() gives the correct string based data + also eof condition (TDD probably)
    - Test that we can get the correct struct from a given string data
    - Test all of this combined with a csv file
    - Pretty similar to our prev strategies