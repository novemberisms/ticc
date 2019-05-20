# ticc

A command line tool to help writing large and complex games for the TIC-80 fantasy console.

# NOTICE

This is still a work in progress. It doesn't even work yet. Eventually this tool will have the following features:

* Support for all the languages supported in the TIC-80
* Preprocessor directives
* Compiler flags
* Optional source code minification
* Unit testing suite

# Reason

One of the limitations of the TIC-80 is that if you want to edit with a proper text editor like VSCode or Sublime, all the source code for your game must be in a single file, which must then be imported into the TIC-80. This can get very unmanageable very quickly, and so usually people write these programs that take a bunch of source files and 'stitch' them together into one file. 

This is my implementation of one of those programs.

I decided to use Go because I've always wanted to learn it, and this seems like a nice easy project to learn with while implementing it. This is always how I've learned languages: by completing some project with it. So far I am very impressed with Go and want to keep 'go'-ing. (haha)

# License

Copyright 2019 Novemberisms

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.