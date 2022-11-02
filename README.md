<h1 align="center">Hi there, I'm Bodya
<img src="https://github.com/blackcater/blackcater/raw/main/images/Hi.gif" width="30" height="30"/></h1>
<h3 align="center">Computer science student, beginner programmer from Abkhazia <img src="https://cdn-icons-png.flaticon.com/512/164/164890.png"
width="30" height="30"></h3>
<h2>Information</h2>

Welcome to my GitHub page. This is my project, which was written in the Go programming language, not without the help of third-party articles, 
the essence of the project is the development of a Rest API for working with a table of books, the database was used by sqlite. Four methods have been written:
<table class="iksweb">
  <tbody>
    <tr>
      <td>Method</td>
      <td>
What is he doing</td>
    </tr>
    <tr>
      <td>GET</td>
      <td>Get all books</td>
    </tr>
    <tr>
      <td>GET</td>
      <td>Get book by id</td>
    </tr>
    <tr>
      <td>PUT</td>
      <td>Create a book</td>
    </tr>
    <tr>
      <td>DELETE</td>
      <td>Delete a specific book by the ID</td>
    </tr>
  </tbody>
</table>
<h2>Installation</h2>
<h3>To Go Programming language</h3>
Go is a statically typed, compiled programming language designed at Google by Robert Griesemer, Rob Pike, and Ken Thompson. 
It is syntactically similar to C, but with memory safety, garbage collection, structural typing, and CSP-style concurrency. 
It is often referred to as Golang because of its former domain name, golang.org, but its proper name is Go.
<br/>
<pre><a href="https://go.dev/">install Go</a></pre>
<hr/>
<h3>VScode</h3>
Visual Studio Code, also commonly referred to as VS Code, is a source-code editor made by Microsoft with the Electron Framework, for Windows, 
Linux and macOS. Features include support for debugging, syntax highlighting, intelligent code completion, snippets, code refactoring, and embedded Git. 
Users can change the theme, keyboard shortcuts, preferences, and install extensions that add additional functionality.
<br/><br/>
<pre><a href="https://code.visualstudio.com/download">install VSCode</a></pre>
<hr/>
<h3>pkg go-sqlite3</h3>
This package can be installed with the go get command:
<br/><br/>
<pre><code>go get github.com/mattn/go-sqlite3</code></pre>
go-sqlite3 is cgo package. If you want to build your app using go-sqlite3, you need gcc. However, 
after you have built and installed go-sqlite3 with <code>go install github.com/mattn/go-sqlite3</code> (which requires gcc), 
you can build your app without relying on gcc in future.
<br/><br/>
<b>Important: because this is a CGO enabled package, you are required to set the environment variable CGO_ENABLED=1 and have a gcc compile present within your path.</b>
<br/><br/>
Full information on <a href="https://github.com/mattn/go-sqlite3">GitHub</a>
<hr/>
<h3>pkg gorilla/mux</h3>
Package gorilla/mux implements a request router and dispatcher for matching incoming requests to their respective handler.
The name mux stands for "HTTP request multiplexer". Like the standard http.ServeMux, 
mux.Router matches incoming requests against a list of registered routes and calls a handler for the route that matches the URL or other conditions
<br/><br/>
With a correctly configured Go toolchain:
<br/><br/>
<pre><code>go get -u github.com/gorilla/mux</code></pre>
Full information on <a href="https://github.com/gorilla/mux">GitHub</a>
