/**
 * Implementation of the browsing stream monitor.
 *
 * @author b.lempereur@outlook.com <Brett Lempereur>
 */

// Address of the web browsing stream.
var STREAM_ADDRESS = 'ws://icra.eu-gb.mybluemix.net/ws/browsing';

// Maximum number of items to keep in the history stream.
var STREAM_LIMIT = 50;

// Client-side application.
var app = angular.module('icraApp', ['angular-websocket']);

// Interface between the application and the stream of browsing information
// from the websocket.
app.factory('HistoryStream', function($websocket) {

  var socket = $websocket(STREAM_ADDRESS);
  var items = [];

  socket.onMessage(function(message) {

    console.log('Received message: ', message);

    // Decode the message and convert its fields to more useful types.
    var item = JSON.parse(message.data);
    item.timestamp = new Date(item.timestamp);
    item.uri = new URI(item.uri);

    // Add the item to the queue and maintain its size.
    items.unshift(item);
    if (items.length > STREAM_LIMIT) {
      items.pop();
    }

  });

  // Construct the service.
  return {
    items: items,
  };

});

// Controller for the stream of browsing information.
app.controller('StreamCtrl', function($scope, HistoryStream) {

  // Make the history stream available.
  $scope.stream = HistoryStream;

});
