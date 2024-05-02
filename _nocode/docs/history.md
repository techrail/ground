# The history of Ground
At the time when I started thinking about Ground, I had already been using Golang professionally for 4 years. I had worked on a large and complicated monolith in Go and had created a few services in Go as well. This included projects that were toys and experiments, large and small and were written for personal needs as well as part of what I did at work.

As time went, I realized that there are a few things which almost every service needed. For example, a web server was almost always on the cards and most of the time it always interacted with a Database (I just use PostgreSQL everywhere). But in certain cases there was more that was needed. For example, one the services that I wrote ran multiple jobs which were basically periodic in nature.

At work, I consolidated the understanding I had about what a Go service needs in a service which actually never got launched. But it gave me the base for creating the _template_ for a go service. As time went by, more and more services got created using that template. But with time, all services that were created diverged a little and accumulated some behavioral changes that were local to their needs.

While one service needed a database, another operated fully over cache. One was a mid-size app that needed awareness about its own state and another which could be allowed to crash. One needed quite a bit of JSON manipulation while another required to operate over sets. Amongst all these, there were some standard libraries that I had started using. I started preferring fasthttp over other routers and http servers. I had written my own structure for handling JSON data, my own way of enumerations. I had developed a system of writing DB layer code. I already had a set of response rendering system.

As part of one of my personal projects, I had written my own logging system and which could send the logs to a PostgreSQL database

However, the difficult part was to bring together the changes in multiple services and make them available to other services as well. They all could benefit from the progression but to update the base template would have helped only a new service that could have been created using that template.

Hence the idea for Ground was born. 

Ground aims to be a library which provides all these services (a web server, a routine manager, a cache manager, some data types, a database layers code generator that outputs simple SQL queries and so on) but does not force the user to use them by keeping them decoupled as much possible.  