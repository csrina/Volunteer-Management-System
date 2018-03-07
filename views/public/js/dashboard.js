function load() {    
    let req = new XMLHttpRequest();
    req.addEventListener("load", function(evt) {
	    let data = JSON.parse(req.response);
	    input(data);
    });
    req.open("GET", "http://localhost:8080/api/v1/dashboard");
    req.send();
}
function input(data) {
    let needed = 0;
    if (data.children == 1) {
	needed = 2.5;
    } else {
	needed = 5;
    }
    let done = document.getElementById("hoursDone");
    let booked = document.getElementById("hoursBooked");
    let table = document.getElementById("events");

    if (data.hoursDone/needed > 0.99) {
	done.style.color = "green"
    } else if (data.hoursDone/needed > 0.66) {
	done.style.color = "yellow"
    } else if (data.hoursDone/needed > 0.33) {
	done.style.color = "orange"
    } else {
	done.style.color = "red"
    }
    
    done.innerHTML = data.hoursDone;
    booked.innerHTML = data.hoursBooked;
   /* for (let i = 0; i < data.eventlist.length; i++) {
	let item = document.createElement("li");
	item.innerHTML = data.eventlist[i];
	table.appendChild(item)
    } */
}

// This function will send a booking removal request to the server
function requestRemoval(event) {
    var promptStr = "Are you sure you want to remove yourself from:\n";

    if (!confirm(promptStr + event.start.toString() + ", in the " + event.room + " room")) {
        return;
    }
    // Block info for booking
    var booking_json = JSON.stringify({
        id:         event.id
    });

    // Make ajax POST request with booking request or request bookign delete if already booked
    $.ajax({
        url: '/api/v1/events/book',
        type: 'POST',
        contentType:'json',
        data: booking_json,
        dataType:'json',
        success: function(data) {  // We expect the server to return json data with a msg field
            alert(data.msg);
            event.booked = !event.booked;
            event.bookingCount--;
            $('#calendar').fullCalendar('updateEvent', event);
        },
        error: function(xhr, ajaxOptions, thrownError) {
            alert("Request failed: " + thrownError);
        }
    });
}
// Configures calendar options
$(document).ready(function() {
    // page is now ready, initialize the calendar...
    $('#calendar').fullCalendar({
        weekends: false,
        header: {
            left: 'title',
            right: 'prev, today, next'
        },
        defaultView: "list",
        duration: {days: 14},        // two week intervals shown for upcoming events
        events: "/api/v1/events/dash",    // EventsFeed with dash as its target
        allDayDefault: false,        // blocks are not all-day unless specified
        themeSystem: "bootstrap3",
        editable: false,                 // Need to use templating engine to change bool based on user's rolego ,
        eventRender: function(event, element, view) {
            element.find('.fc-title').append("  " + event.bookingCount + " / 3  ");
        },
        // DOM-Event handling for Calendar Eventblocks (why do js people suck at naming)
        eventOverlap: false,
        eventClick: function(event, jsEvent, view) {
            requestRemoval(event);
        },
        businessHours: {
            // days of week. an array of zero-based day of week integers (0=Sunday)
            dow: [1, 2, 3, 4, 5], // Monday - Thursday
            start: '8:00',
            end: '18:00'
        }
    })
});

load();
