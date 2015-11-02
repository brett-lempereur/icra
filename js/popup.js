/**
 * Controller for the in-browser popup window.
 *
 * @author b.lempereur@outlook.com <Brett Lempereur>
 */

/**
 * Update configuration settings and connect to the message broker.
 */
function connect() {

    var identity = document.getElementById('identity').value;
    var paths = document.getElementById('paths').checked;
    var hostname = document.getElementById('hostname').value;
    var port = Number(document.getElementById('port').value);
    var ssl = document.getElementById('ssl').checked;
    var username = document.getElementById('username').value;
    var password = document.getElementById('password').value;

    // Store configuration settings so they are available when the browser
    // starts and trigger a connect event.
    chrome.storage.local.set({
        identity: identity,
        paths: paths,
        hostname: hostname,
        port: port,
        ssl: ssl,
        username: username,
        password: password
    }, function() {
        chrome.runtime.sendMessage({"type": "connect"});
    });

}

/**
 * Disconnect from the message broker.
 */
function disconnect() {
    chrome.runtime.sendMessage({"type": "disconnect"});
}

/**
 * Restore configuration settings from local storage.
 */
function restore() {

    var identity = document.getElementById('identity');
    var paths = document.getElementById('paths');
    var hostname = document.getElementById('hostname');
    var port = document.getElementById('port');
    var ssl = document.getElementById('ssl');
    var username = document.getElementById('username');
    var password = document.getElementById('password');

    // Get configuration settings and modify the interface.
    chrome.storage.local.get({
        identity: "",
        paths: false,
        hostname: "localhost",
        port: 8080,
        ssl: true,
        username: "",
        password: ""
    }, function(values) {
        identity.value = values.identity;
        paths.checked = values.paths;
        hostname.value = values.hostname;
        port.value = values.port;
        ssl.checked = values.ssl;
        username.value = values.username;
        password.value = values.password;
    });

}

/**
 * Update the status display.
 */
function updateStatus() {

    chrome.runtime.sendMessage(
        {"type": "getStatus"},
        function displayStatus(response) {
            document.getElementById("status").innerHTML = response.result;
        }
    );

    setTimeout(updateStatus, 500);

}

// Start the update status process, restore configuration values, and bind
// actions to the buttons when the document is loaded.
document.addEventListener('DOMContentLoaded', function() {
    updateStatus();
    restore();
    document.getElementById('connect').addEventListener('click', connect);
    document.getElementById('disconnect').addEventListener('click', disconnect);
});
