# Ground Core

The core contains the structures and mechanisms to view and control the state of different components that have been initialised and being used. It is important because depending on which components have been initialized, we might have to alter the behavior of other components. 

For example, if the logging system has been initialized then it means that the project using ground is ready to utilize the bark logging system. In such a case, we can use the bark client to log. If not, we might have to use `fmt.Println` statements (or something similar). 

Another example - if the database systems have been configured then maybe other parts of the system can use that. But if not, we can stop such parts from actually being called.