# getfiled file generator

Getfiled will appear to be a normal file, but when accessed, will generate the contents of the file on the fly. Because it's a normal file, programs do not need to be aware that the file is dynamically generated- The only difference would be extra latency when reading the file.

As an example, you could ask getfiled to provide your Nginx configuration, but have it fetch the latest configuration from a central server, like so:  
`getfiled /etc/nginx.conf "ssh root@master.example.com cat /etc/nginx.conf"`  
The same method could be used to centralize all of your other configurations, like cron, supervisor, etc.

For more complex tasks, you can just have it run a binary, which could have caching/diffing from a server, generate custom configuration files, reloading programs, or anything else. Getfiled is byte-friendly, so you can even have other resources such as jpegs or dll's be dynamically generated. If you want your binary to be run often, such as for updating a cached copy, use cron.

# how it works

Getfiled is atomic and can be read by hundreds of programs simultaneously. Internally, getfiled created a FIFO with the expected name, and as soon as something reads that FIFO, it renames it and created a new FIFO for the next reader to find and spawns the generator binary. It then points the spawned process's stdout to that FIFO, so the spawned process can take as long as it needs to write, or even write in chunks. Once that process closes, getfiled closes the FIFO.

# future features

It might be useful to implement polling+caching of the generated files as an option. Also, it might be good to run the binary once, and just write a filename to it's STDIN when it needs to generate a file to simplify getfiled, splitting the live-generation to a different binary.
