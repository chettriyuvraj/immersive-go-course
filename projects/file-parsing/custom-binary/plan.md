# Plan

- How would you like to read the binary file?
    - Sort of 'getNextRecord' again, similar to previous files
    - There is an endian-ness mark at the very start
    - Can we create a struct which stores: endianess + io.reader + any other info
- What would you like to test for
    - Read next record
    - Decoding this byte record into player data
    - A function that does this for a file

### Observations post implementation
- Didn't implement completely, quite similar to prev strats
- Check reference implementation
- Tests don't exactly read well + the idea of creating a PlayerBinaryData struct also seems a bit unnecessary (?)