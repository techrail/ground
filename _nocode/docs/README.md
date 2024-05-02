# Documentation for Ground
This directory contains documentation related to Ground project. You can start with the [history](./history.md) to
understand where Ground stands and why it is not the standard library.

Ground does not try to provide you with a new way to deal with URL routing or creating a new driver for
valkey (or Redis) or your SQL databases. Ground does not try to give you a new way to work with anything that you are 
probably habituated to. So you might ask yourself what is Ground, after all! Here are some principals for Ground:

1. Ground will provide framework style functionality. 
2. All sub-modules of Ground will be as independent as possible. 
3. However, they must have the capability to work across different modules (e.g. The database code generator would be able to use the caching system). 
4. Ground will not rediscover wheels. If there is a good library that already does a work, Ground will provide a thin wrapper around it to get the work done. There is already some great code out there waiting to be reused.
5. Ground will provide some functionality on its own where the said functionality does not exist.





