#Plan


- In a normal JSON file, ideally parsing would involve passing the io.Reader straight to decoder
- Here, seems like we would need an additional step
- We can go about it two ways:
    - parse the file in one go inside function, ignoring all the commented out lines and return the clean data
    - separate the reading and parsing, so we read the file line-by-line, parsing each line as it comes
- Second approach seems better as it involves separation of functions + not holding a lot of data in memory at a time
- Test wise:
    - Test the reading i.e. we pass an io.reader to function and it returns to us the next newline separated line (as a byte slice) - so some sort of buffered io, (use bufio package?)
    - Test the parsing i.e. we pass a byte slice and it returns true false depending on whether the line starts with '#' or not
    - Test these two combined i.e parsingRepeatedJsonFile