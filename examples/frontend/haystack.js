/*
	Haystack JS tracker to push data to the default stream
*/

(function () {
	// protocol selector
	function buildUrl(url) {
        return ("https:" == a.location.protocol ? "https" : "http") + "://" + url
    }

    "use strict";
    var endPoint = buildUrl("localhost/hstk");

    // wrap around events | start with pageviews only
    // add listeners
    // construct message
    // push message to endpoint

})()