/**
 * Implements the NavigationMonitor object.
 *
 * @author b.lempereur@outlook.com <Brett Lempereur>
 */

/**
 * Maintains a connection to a message broker, and broadcasts navigation events
 * as messages.
 *
 * @constructor
 */
function NavigationMonitor() {

    /**
     * The message broker client.
     *
     * @type {Paho.MQTT.Client}
     * @private
     */
    this.client_ = null;

    /**
     * The identity of the client.
     *
     * @type {!string}
     * @private
     */
    this.identity_ = "";

    /**
     * Whether to include paths in messages.
     *
     * @type {!bool}
     * @private
     */
    this.paths_ = false;

    /**
     * The status of the connection to the broker;
     *
     * @type {!string}
     * @private
     */
    this.status_ = "Disconnected";

    // Bind to the chrome navigation events.
    chrome.webNavigation.onCompleted.addListener(
        this.onNavigationCompleted_.bind(this)
    );

    // Bind to receive messages from the user interface.
    chrome.runtime.onMessage.addListener(this.onChromeMessage_.bind(this));

}

NavigationMonitor.prototype = {

    /**
     * Connect to the message broker.
     */
    connect: function() {

        var monitor = this;

        // If we have a connection, disconnect it first.
        if (this.client_ !== null) {
            this.disconnect();
        }

        // Get the current configuration settings and connect to the broker.
        chrome.storage.local.get([
            "hostname", "port", "identity", "paths", "ssl", "username",
            "password"
        ], function(values) {
            monitor.identity_ = values.identity;
            monitor.paths_ = values.paths;
            monitor.client_ = new Paho.MQTT.Client(
                values.hostname,
                values.port,
                values.identity
            );
            monitor.client_.onConnectionLost = monitor.onBrokerConnectionLost_.bind(this);
            monitor.client_.connect({
                userName: values.username,
                password: values.password,
                useSSL: values.ssl,
                onSuccess: function() {
                    monitor.status_ = "Connected";
                },
                onFailure: function() {
                    monitor.client_ = null;
                    monitor.status_ = "Failed to connect";
                }
            });
        });
    },

    /**
     * Disconnect from the message broker.
     */
    disconnect: function() {
        if (this.client_ == null) {
            return;
        }
        this.client_.disconnect();
        this.client_ = null;
    },

    /**
     * Handler for the 'onCompleted' event of the web navigation interface that
     * constructs and broadcasts a message describing the visited page.
     *
     * @param {!Object} details The event data generated for this request.
     * @private
     */
    onNavigationCompleted_: function(details) {

        // Check whether we are connected.
        if (this.client_ == null) {
            return;
        }

        // Parse and clean the url.
        var uri = URI.parse(new URI(details.url).normalize().toString());
        delete uri.username;
        delete uri.password;
        delete uri.query;
        delete uri.fragment;
        if (!this.paths_) {
            uri.path = null;
        }

        // We are only interested in HTTP and HTTPS traffic.
        if (uri.protocol != "http" && uri.protocol != "https") {
            return;
        }
        if (uri.protocol == "https") {
          uri.path = null;
        }

        // Construct and send the update message.
        var payload = JSON.stringify({
            timestamp: new Date(details.timeStamp),
            identity: this.identity_,
            uri: uri
        });
        var message = new Paho.MQTT.Message(payload);
        message.destinationName = "Browsing/" + this.identity_;
        this.client_.send(message);

    },

    /**
     * Handler for the 'connectionLost' event of the message broker client.
     */
    onBrokerConnectionLost_: function(reason) {
        monitor.client_ = null;
        monitor.status_ = "Disconnected";
        if (reason.errorCode != 0) {
          monitor.connect_();
        }
    },

    /**
     * Handler for messages from the popup.
     *
     * @param {!{type:string}} message  The external message to answer.
     * @param {!MessageSender} sender   The context of the sender.
     * @param {!function}      respond  Function to call to send a response.
     * @private
     */
    onChromeMessage_: function(message, sender, respond) {
        switch (message.type) {
        case "connect":
            this.connect();
            respond({});
            break;
        case "disconnect":
            this.disconnect();
            respond({});
            break;
        case "getStatus":
            respond({result: this.status_});
            break;
        default:
            respond({});
        }
    }

};
