# Database Code Generator

Code generators are complex. Hence in this package, we shall be adding more comments than usual to make sure that process is clear about what we are doing.

## What is an isolated table?
Within the generator code, you would find the reference to this term.
This term is used to refer to a table that no other table links to, and the table itself does not link to any other table by foreign keys. 
This is useful when you want to generate code for a table that is not part of the main application but is still needed for some other purpose (e.g. Schema migration, testing, etc.).