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

## Installation

This is a bit more stable than my other Chrome extension, but installation is
by hand because I don't see there being a large installation base for this:

1. Clone the repository from GitHub.
2. Enable "Developer Mode" in the Google Chrome extensions configuration page.
3. Selected "Load unpacked extension...".
4. Browse to and select the location where you cloned the repository.

Raise an issue if I'm wrong about the demand.

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

