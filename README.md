# RESTAD Group Manager

This project provides the front end and api that talks to a RESTful backend (currently RESTAD, but moving to [Addict](http://github.com/dthree/addict) in the future) for managing Active Directory groups for the BYU Physics and Astronomy Department.

![Group manager in use](example.gif)


## Build
The fontend must be built first because it is embeded into the resulting server binary.

### Frontend
Enter directory `react-improved` and fetch dependencies with `yarn install`. Finally run `yarn webpack`

## Build backend
Make sure your `GOPATH` is defined and then run `make setup` to fetch all the dependencies.
Once you have the dependencies you can build the application with `make release`. Note that CGO is required so cross compiling for linux from windows is not well supported.

