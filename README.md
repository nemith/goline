netconf
=======

[![Build Status](https://travis-ci.org/nemith/go-goline.png)](https://travis-ci.org/nemith/go-goline)

Simple implemenation of a readline like facility that is heavly based on the C library [linenoise](https://github.com/antirez/linenoise).  Uses syscalls directly from Golang to implment the low level terminal functions and doesn't wrap any existing C library.  BSD and Linux support currently.  No Windows support.

Currently only does simple line editing.  Upcoming versions will include:
 * Custom key bindings
 * History
 * Completion

Very alpha quality right now.  API _will_ change.

See examples folder for examples API implementations

Features
--------
* Support for SSH transport using go.crypto/ssh. (Other transports are planned)
* Built in RPC support (in progress)
* Support for custom RPCs
* Independant of XML library.  Free to choose encoding/xml or another third party library to parse the results.

Install
-------

    go get github.com/nemith/go-goline/goline
    
License
-------
(BSD 2)

Copyright © 2013, Brandon Bennett

Linenoise copyrights:
Copyright (c) 2010-2013, Salvatore Sanfilippo <antirez at gmail dot com>
Copyright (c) 2010-2013, Pieter Noordhuis <pcnoordhuis at gmail dot com>

All rights reserved.

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:

(1) Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.

(2) Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS “AS IS” AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

The views and conclusions contained in the software and documentation are those of the authors and should not be interpreted as representing official policies, either expressed or implied.

Authors and Contributors
------------------------
* [Brandon Bennett](http://www.linkedin.com/in/brandonrbennett)