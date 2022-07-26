ocpp1.6-go
=================
![License](https://img.shields.io/dub/l/vibe-d.svg)

Open Charge Point Protocol (OCPP) is a standard open protocol for communication between Charge Points and Central System and is designed to accommodate any type of charging technique. 

The library is representing implementation of OCPP version 1.6 in Go.
Code in branch is based on JSON type communication (OCPP-J) as SOAP will no longer be supported in future versions.

## Status & Roadmap

Planned milestones and features:

- [x] OCPP 1.6 repo structure
- [x] Example of OCPP implementation as Central Sysytem (3 actions)
- [ ] OCPP 1.6 Core
- [ ] Add test cases for library

## OCPP 1.6 Usage

Go version 1.18+ is required.

Installation
------------

Use go get.

	go get github.com/CoderSergiy/ocpp16-go

Then import the validator package into your own code.

	import "github.com/CoderSergiy/ocpp16-go"


How to Contribute
------

Make a pull request...

License
-------
Distributed under MIT License, please see license file within the code for more details.

Maintainers
-----------
This project mantained by one person at this point.
If you are interested in project please reach out to me https://github.com/CoderSergiy
