// Callback function for drag/drops and resizes of existing events
// Note: We dont want this to be populated if we aren't admin.
// post-demo will refactor this out into templates populated differently based on the role of the user
function storeChangesToEvent(event, delta, revertFunc, jsEvent, ui, view) {
    // Extract block data required for updating on server
    let event_json = JSON.stringify({
        id: event.id,
        start: event.start,
        end:   event.end,
        modifier: event.value
    });
    // Make ajax post request with updated event data
    $.ajax({
        url: '/api/v1/events/update',
        type: 'POST',
        contentType:'json',
        data: event_json,
        dataType:'json',
        success: function(data) { 
            alert(data.msg);
        },
        error: function(xhr, ajaxOptions, thrownError) {
            revertFunc();
            alert("Request failed: " + xhr + "\n" + thrownError);
        }
    });
}

// This function will send a booking request to the server
// Like the updateEvent function, post-demo I will refactor this out
// and have the templates populate based on role. Additional auth checks
// server side to ensure correct user/role and such should still take place
function requestBooking(event, jsEvent, view) {
    let promptStr = "Confirm booking ";
    console.log(event.booked);
    if (event.booked === true) {
        promptStr += "Cancellation (";
    } else {
        promptStr += "Booking (";
    }
    if (!confirm(promptStr + event.start.toString() + ", in the " + event.color + " room)")) {
        return;
    }
    // Block info for booking
    let booking_json = JSON.stringify({
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
            if (event.booked === true) {
                event.bookingCount++;
            } else {
                event.bookingCount--;
            }
            $('#calendar').fullCalendar('updateEvent', event);
        },
        error: function(xhr, ajaxOptions, thrownError) {
            alert("Request failed: " + thrownError);
        }
    });
}

$(document).ready(function() {
    loadAddEvent();
    // page is now ready, initialize the calendar...
    $('#calendar').fullCalendar({
        // Education use (both now and if deployed!)
        weekends: false,
        header: {
            left: 'today',
            center: 'prev, title, next',
            right: 'agendaWeek, month'
        },
        defaultView: "agendaWeek",
        events: "/api/v1/events",    // link to events (bookings + blocks feed)
        allDayDefault: false,        // blocks are not all-day unless specified
        themeSystem: "bootstrap3",
        editable: true,                 // Need to use templating engine to change bool based on user's rolego ,
        eventRender: function(event, element, view) {
            element.find('.fc-time').css("font-size", "1.5em");
            element.find('.fc-title').css("font-size", "1.5em");
            if (event.booked) {
                element.find('.fc-title').prepend("<span class='glyphicon glyphicon-pushpin' aria-valuetext='You are booked in this block!'></span><br/>");
            } else {
                element.find('.fc-title').prepend("<br/>");
            }
            element.find('.fc-title').append("<br/>" + event.bookingCount + " / 3<br/>");
        },
        // DOM-Event handling for Calendar Eventblocks (why do js people suck at naming)
        eventOverlap: function(stillEvent, movingEvent) {
            return stillEvent.color === movingEvent.color;
        },
        // When and event is drag/dropped to new day/time --> updates db & stuff
        // revertFunc is called should our update request fail
        eventDrop: function(ev, delta, revertFunc, jsEvent, ui, view) {
            storeChangesToEvent(ev, delta, revertFunc, jsEvent, ui, view);
        },
        // When an event is resized (post duration change) it will callback the function
        // revertFunc is fullCalendar function which reverts the display should the request fail
        eventResize: function(ev, delta, revertFunc, jsEvent, ui, view) {
            storeChangesToEvent(ev, delta, revertFunc, jsEvent, ui, view);
        },
        eventClick: function(event, jsEvent, view) {
                requestBooking(event, jsEvent, view);
        },
        businessHours: {
            // days of week. an array of zero-based day of week integers (0=Sunday)
            dow: [1, 2, 3, 4, 5], // Monday - Thursday
            start: '8:00', // a start time (10am in this example)
            end: '18:00', // an end time (6pm in this example)
        },
        // Controls view of agendaWeek
        minTime: '07:00:00',
        maxTime: '18:00:00',
        allDaySlot: false,       // shows slot @ top for allday events
        slotDuration: '00:30:00' // hourly divisions
    })
});


function loadAddEvent() {
    let xhttp = new XMLHttpRequest();
    xhttp.addEventListener("loadend", () => {
        if (xhttp.response > 300) {
            alert("ERROR: Could not load class list");
        }
        let classes = JSON.parse(xhttp.response);
        let tmpl = document.querySelector("#tmpl_EventForm").innerHTML;
        let func = doT.template(tmpl);
        document.querySelector("#eventForm").innerHTML = func(classes);
        document.querySelector("#submit").addEventListener("click", submitEvent);
    });
    xhttp.open("GET", "http://localhost:8080/api/v1/admin/classes")
    xhttp.send();
}

function submitEvent() {
    if (document.querySelector("#start").value == "" 
        || document.querySelector("#end").value == ""
        || document.querySelector("#room").value == ""
        || document.querySelector("#modifier").value == "" ) {
        alert('Please fill out all options');
        return;
    }

    let xhttp = new XMLHttpRequest();
    xhttp.addEventListener("loadend", () => {
        if (xhttp.status > 300) {
            alert('ERROR: Could not create event.')
            return;
        }
        //loadAddEvent();
    });

    let newStart = document.querySelector("#start").value;
    let newEnd = document.querySelector("#end").value;
    let newRoom = parseInt(document.querySelector("#room").value);
    let newMod = parseInt(document.querySelector("#modifier").value);
    let newNote = document.querySelector("#note").value;
    xhttp.open("POST", "http://localhost:8080/api/v1/events/add");
    xhttp.send(JSON.stringify({start:newStart, end:newEnd,
            room:newRoom, modifier:newMod, note:newNote}));
}