# Standard responses
Some response bodies are required way too commonly, or frequently or both in a typical API server. Some of those are: 

1. A blank response with no content (this can be true for 200, 300 as well as 400 range of HTTP response codes)
2. A response with a single message (maybe for displaying some message upon success or error)
3. A response representing an error with details (such as an error code, a display message, some other details etc.)
4. A response for reporting some basic information about the service, indicating that it is alive.

This directory/module shall contain such responses and shall be available for anyone importing the ground project. 
This list is in no way exhaustive and any project using Ground should have its own set of network responses for its
own use-case.
