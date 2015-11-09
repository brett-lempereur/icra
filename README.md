# ICRA

Giving the home secretary a helping hand.

## What

I'm not sure if people get the scope of the "snooper's charter" that is still
on the table. So, I'll be using the extension in this repository to stream my
web browsing metadata in real-time to anyone that cares to look at it for a
while.

If anyone else wants to take part, you're more than welcome; I'll setup an
account and topic on the MQTT server.

When I get the front-end written, I'll put a link up here.

## Chrome Extension

This is a bit more stable than my other Chrome extension, but installation is
by hand because I don't see there being a large installation base for this:

1. Clone the repository from GitHub.
2. Enable "Developer Mode" in the Google Chrome extensions configuration page.
3. Selected "Load unpacked extension...".
4. Browse to and select the "extension/chrome" directory where you cloned the
   repository.

Raise an issue if I'm wrong about the demand.

## Website

This is the web application that users can connect to and monitor any browsing
activity from ICRA extensions.

To run the application locally:

1. [Install Go](https://golang.org/doc/install).
2. Download and extract this code into your $GOPATH/src directory.
3. Change into the application directory.
4. Run 'go run app.go'.
5. Access the application in a web browser at http://localhost:8080.

## Todo

As ever, help with any of these things would be appreciated:

1. ~~Building a single-page site that will stream my traffic.~~
2. ~~Deploying an MQTT server without a 10-client limit.~~
3. ~~Fixing outstanding bugs in the extension and making it more robust.~~
4. Localising the extension user interface.

Help with the second todo would be great.

## Thanks

The [Eclipse Paho](https://projects.eclipse.org/projects/technology.paho) team
and Rodney Rehm for [URI.js](https://medialize.github.io/URI.js/).
