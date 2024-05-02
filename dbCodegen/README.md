# Database Code Generator

Code generators are complex. Hence, in this package, we shall be adding more comments than usual to make sure that process is clear about what we are doing.

## Current capabilities and features

The database code generator reads the schema of a **PostgreSQL** database (that means it does not work with any other database) and generates some code for the database. It will generate: 

1. The type structs for the tables that it finds. 
2. DAO (Data Access Objects) for each table. This DAO type is blank and is used for scoping (via receiver functions) data access for the table. 
3. Functions against the type structs to do the basic CRUD (Create, Read, Update and Delete) operations on the table.
4. Base validations against constraints that it can find (e.g. max length checks for a VARCHAR type). 
5. Functions under DAO that corresponds to the indexes it finds for the table. 
6. Functions under type struct and DAO relating to the Foreign Keys (both incoming and outgoing connections) that the table contains (even if they are across schemas).
7. Network layer structs and functions that allow you to fill a network layer struct from a DB layer struct. This segregation allows you to both consume and respond in a structure different from your DB and allows you to hide certain fields on the network layer (e.g. a User's Password Hash can be redacted from the response automatically on the network layer).

The queries that are generated do not contain any ORM code. Instead, they use `sqlx` and _raw SQL queries_ that map to the type structs. This heavily simplifies the readability and reduces dependency on an ORM which can complicate things at wrong places!