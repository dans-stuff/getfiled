# getfiled file generator

Getfiled will appear to be a normal file, but when accessed, will generate the contents of the file on the fly.

As an example, you could ask getfiled to provide your Nginx configuration, but have it fetch the latest configuration from a central server, like so: `getfiled /etc/nginx.conf "ssh root@master.example.com cat /etc/nginx.conf"`

For more complex tasks, you can just have it run a binary, which could have caching/diffing from a server, generate custom configuration files, or anything else.

# how it works

Getfiled is atomic and can be read by hundreds of programs simultaneously. Internally, getfiled created a FIFO with the expected name, and as soon as something reads that FIFO, it renames it and created a new FIFO for the next reader to find and spawns the generator binary. It then points the spawned process's stdout to that FIFO, so the spawned process can take as long as it needs to write, or even write in chunks. Once that process closes, getfiled closes the FIFO.
