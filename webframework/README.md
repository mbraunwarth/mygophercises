# Web Frameworks in Go

This project is rather based on a [blog post](https://www.nicolasmerouze.com/build-web-framework-golang) from Nicolas Mérouze.
It is intendet to build an understanding for the functionality of web frameworks in general
and in particular in Go.

### Frameworks

A framework serves as a building structure for projects including not only code in form
of libraries but also components and the structural design for the project you build.  
There are two main kinds of web frameworks:
- Rails-like frameworks which have almost every feature built-in you can think of making the development process fast and
- Sinatra-like frameworks with some built-in features and a router for managing your endpoints.

A simple framework includes a **router** managing requests and for the handlers, a **middelware system** and the **handler** itself which processes the incoming request and serving a response.