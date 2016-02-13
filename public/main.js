
//need to call this function when a file is clicked
//also close the connection if its open already
var conn;

// $(function() {
function startConnection(identity, filename) {
    if (conn) {
        conn.close(); //close the websocket    
    }
    
    if (filename == undefined) {
        return;
    }
    if (!window["WebSocket"]) {
        return;
    }

    var content = $("#content");
    console.log("Connection: " + 'ws://' + window.location.host + '/ws?id='+identity+"&filename="+filename)
    conn = new WebSocket('ws://' + window.location.host + '/ws?id='+identity+"&filename="+filename);

    // Textarea is editable only when socket is opened.
    conn.onopen = function(e) {
        content.attr("disabled", false);
    };

    conn.onclose = function(e) {
        content.attr("disabled", true);
    };

    // Whenever we receive a message, update textarea
    conn.onmessage = function(e) {
        if (e.data != content.val()) {
            content.val(e.data);
        }
    };

    var timeoutId = null;
    var typingTimeoutId = null;
    var isTyping = false;

    content.on("keydown", function() {
        isTyping = true;
        window.clearTimeout(typingTimeoutId);
    });

    content.on("keyup", function() {
        typingTimeoutId = window.setTimeout(function() {
            isTyping = false;
        }, 1000);

        window.clearTimeout(timeoutId);
        timeoutId = window.setTimeout(function() {
            if (isTyping) return;
            conn.send(content.val());
        }, 1100);
    });
}
// });