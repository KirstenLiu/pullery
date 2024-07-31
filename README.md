# Kirsten Liu

## Execution

- Execution instructions: put the input json into file "input.json", the output should be printed directly in STDOUT. 
- Additional information that may be helpful to the reviewer:
	1. I also include some useful logs for debug usage, please feel free to change the log level.
	2. In the transformation instruction, **N** denotes the value's data type, and the sanitize of trailing and leading whitespace is only defined to "value". So by logic from instruction, whitespace in dataType is not processed and considered illegal.
	For example, "null_1": { "NULL ": "true"} should be considered illegal.
	But in the sample output it is included, so I process the trailing zero for data type as well.
	3. Since map is always sorted in fmt print for morden Go version, I skip the implemetation of the part of sorting the map lexically.
	4. My apologies, but I didn't see the instruction about "report the implementation processing time" from the beginning. If I did, I would remember to create the repo on my local env so git log can tell exactly the process. I only saw it in the middle (while I was impementing the null part), and my git log started there. I will estimate that I spent roughly 1 to 1.5 hours before I started the git repo. So in total I spent about 3 hours in this.
	5. My local go version: go1.21.5 darwin/arm64
	6. I added the original test file (input.json) and an additional test file I used to test different cases (input2.json) in the repo as well, please feel free to use it. It's not fully coverage, I just add some classic test cases to check the corectness of the code. I normally like to add tests in my go code, but since it's not in the instruction and might be out of scope, I don't add them this time.

