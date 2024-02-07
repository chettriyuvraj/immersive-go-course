# README


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