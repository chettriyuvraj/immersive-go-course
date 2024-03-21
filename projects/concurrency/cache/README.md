# README

## Plan

Goals

- Statistics
- Concurrency-safe
- LRU

Plan

- Let's get something working first.
- Part 1: Concurrency-safe: only one of Get/Put can be executed at a time
    - Add a global mutex i.e. same mutex used for both
    - How would you test this?
        - https://forum.golangbridge.org/t/how-do-you-unit-test-a-concurrent-data-structure/26912/2 This explains it a bit
        - Instead of testing for something, use different methods that might catch a deadlock
        - Use t.Parallel() to run normal tests, higher chance of catching race/deadlock
        - Use Go's race detector: https://go.dev/doc/articles/race_detector
        - Fuzzy testing: https://go.dev/doc/security/fuzz/, can be used to seed inputs and then generates random inputs on it's own

- Part 2: Stats, again concurrency safe
    - Cache-level (cache struct can have additional fields to accomodate this)
        - hits
        - misses
        - writes
    - Cache item level (cache can use a KV pair where V is a type of CacheItem, which stores value + statss)
        - timesRead
    - Testing mechanism similar to above
- Part 3: LRU using hashmap + double linked list:
    - Testing mechanism similar to above
