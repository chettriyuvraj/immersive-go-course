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



### Observations on my solution vs reference solution

- Reference Solution: https://github.com/CodeYourFuture/immersive-go-course/tree/impl/interfaces/projects/interfaces
- OurBytesBuffer
    - Our implementation of ourbytesbuffer seems a bit more verbose than the reference impl, this is mainly due to our choice of having extra buffer space (2x) any time the buffer gets filled + at the start, thus the endPoint of the content is not always equal to the length of the buffer
    - I have implemented the non-table-driven version of tests (inadvertently - I thought we had to test EXACTLY as the question asked)
    - Each 'factor' to test is segregated in the reference solution, I have combined multiple ones. Should be segregated for readability.
    - Reference solution uses the require package, which can be used to improve readability in my solution BY A LOT.
    - The section explaining trade-offs between table and non-table driven tests is VERY GOOD!

- Filtering pipe
    - My interpretation: _The only thing we can test is the count of the bytes written_ -> Acc to the reference soln, _io.Writer is documented to return the number of bytes processed, not used_ so n will always be the entire length of the string including digits
    - What else can we test then? We can pass a ReadWriter() like Buffer and at the end check the contents of the buffer has no digits
    - Also, the table-driven style goes very well here! 


- Tests
    - In non-table driven tests (buffer_test.go), the tests were always as small as possible i.e. fine grained, while I have chosen to test multiple things in a single test, I think this has hampered readability in my tests.