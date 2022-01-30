# Response Converter
A tool that creates requests with the given urls and converts its response to md5 hash.

# Prerequisites
Before you begin you must have Go installed and configured properly for your computer. Please see https://golang.org/doc/install
To run the project, you need to clone this project to your computer. Link: https://github.com/ahoirg/responseConverter.git

# Run Response Converter
To run Response Converter, first go to the path where the project is installed in the terminal.
You can run it using the “go run” command, then project name. After that you can write urls.
```go
$ go run .\main.go  www.google.com http://www.google.com https://github.com/ahoirg  
http://www.google.com e2ffa54c2b8dea2d4af00aabc883038d
http://www.google.com ea57a2dd255cd6130afaf636a7d5f57c
https://github.com/ahoirg e2b5aab2c4652718c7bb3ce963e5b7e3
```

By default, up to 10 goroutines can run simultaneously. It can be changed with the "-parallel" command.
```go
go run .\main.go  -parallel 1 www.google.com http://www.google.com https://github.com/ahoirg
http://www.google.com 696eeff791c1c265e70355909e607995
http://www.google.com 510c01000a10b892be35c237c9f251cf
https://github.com/ahoirg 8e0f63239f97ac718999a17e494e57db
```

The tool will not print any response of invalid url patterns and unused domain names. It will ignore them.
```go
go run .\main.go  www.inValidDomainName.com http:///google.com

```

# Running Tests
To run the tests, "go test" command should be used in the project directory.
```go
$ go test
ok  responseConverter    7.797s
```

# How To Contribute
A url whose request will be created must be valid. I wrote this validation using the links below. 
However, urls that are still not valid can pass this validation. A better solution can be made.

# References
https://www.geeksforgeeks.org/html-dom-url-property/
https://www.google.com/support/enterprise/static/gsa/docs/admin/current/gsa_doc_set/admin_crawl/url_patterns.html

# License
This project is licensed under the GNU General Public License v3.0 - see the LICENSE.md file for details
