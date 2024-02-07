# README

## OurBytesBuffer

- Functions being implemented
    - Behaviour same as bytes.Buffer:
        - NewBuffer
        - Read
        - Write
        - Bytes
    - Implementation details and limitations:
        - NewBuffer: Initializes a byte slice underneath of length 2 * given slice and copies contents to internal slice
        - Read:
            - Keeps an internal byte offset and reads the data + updates pointer depending on size of provided slice
            - Never clears the internal storage slice, but always moves pointer ahead
        - Bytes:
            - Reads data - internalSlice[byteOffset:]
            - doesn't move pointer ahead
        - Write:
            - If data to be written exceeds the internal byte slice capacity, provisions extra space all at once and then writes
            - Thus writing one large chunk of data would probably be faster (single system call internally to provision the space)
            - As opposed to writing the same data in smaller chunks (multiple system calls to provision the space)
        - ResetOffset:
            - sets the internal offset to 0 so Read and Bytes can essentially re-read the data

        - We let our buffer length be constrained by max size of the offset pointer = 2 ^ 31
        - Thread Safety (?):


## FilteringPipe

- Not implementing this, writing down pseudocode
- Define struct with an io.writer field
- Let io.writer be fed on initialization
- It's write method should first take 1 pass through byte array and remove all numerical bytes
- After which it should call write of the underlying io.Writer and return result
- Since this is just a writer; how should we write tests for it (?)
    - The only thing we can test is the count of the bytes written, i.e. pass a normal string without numbers, n should return count equal to string length - write one of this
    - Write a couple of different types of number tests (continuing eg. "abc467") (single digits eg. "1a2b") (only digits eg. "123")