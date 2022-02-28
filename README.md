# turbo-twitter
Fastest Turbo twitter of all time. (hijacks twitter handle)

`AuthToken` is drived from the accuont cookie

This is actually a port from a python version with the same name, I found that on discord and ported it. [The twitter-turbo python project by Aero-blue](https://github.com/Aero-Blue/twitter-turbo)
I can't figure how it does work despite static CSRF and Bearer but it does!

The way it works is it constantly checks for usernames to become open and tries acquiring them in milliseconds (depends on your ping)

It is super fast and can have thousands of threads, for the performance part depending on your ping + CPU you can have 200k+ connections 
at the same time!

##Running
To run project simply install Golang and clone the project then
```bash
  $ cd turbo-twitter
  $ go run main.go
```
