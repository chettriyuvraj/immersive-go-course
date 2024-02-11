# README - solution

## Thoughts
	- Didn't end up implementing as I had done a similar exercise here: https://github.com/chettriyuvraj/One2N-Golang-Bootcamp/blob/main/Building-Linux-CLIs/
	- Implementation and tests are relatively straightforward
	- Avoid using global rootCmd as I have done in the repo above, instead create a new rootCmd in Execute() func
	- Testing made easy by cmd.OutOrStdout() function, use an io.Writer() and test the output, demonstrates the simplicity and power of interfaces.
 	- Reading closing dirs/files using Go calls very similar to system calls for files and dirs so pretty straightforward (uses them under the hood ofcourse)
