// Callback function for drag/drops and resizes of existing events
// Note: We dont want this to be populated if we aren't admin.
// post-demo will refactor this out into templates populated differently based on the role of the user
function storeChangesToEvent(event, delta, revertFunc, jsEvent, ui, view) {
    // Extract block data required for updating on server
    let event_json = JSON.stringify({
        id: event.id,
        start: event.start.format() + "Z",
        end:   event.end.format()+ "Z",
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
            alert("Request failed: " + xhr.responseText);
        }
    });
}

// REmoves an event from the calendar and the associated TB/Bookings from the database
function removeEvent(event, jsEvent, view) {
    yn = confirm("Are you sure you want to delete this event?")
    if (!yn) {
        return false; // event should not be deleted
    }
    return true; // event was deleted
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
        agendaEventMinHeight: 100,
        defaultView: "agendaWeek",
        contentHeight: 'auto',
        events: "/api/v1/events/scheduler",    // link to events (bookings + blocks feed)
        allDayDefault: false,        // blocks are not all-day unless specified
        themeSystem: "bootstrap4",
        editable: true,                 // Need to use templating engine to change bool based on user's rolego ,
        eventRender: function(event, element, view) {
            element.find('.fc-time').css("font-size", "1.2em");
            element.find('.fc-title').css("font-size", "1.2em");
            if (event.booked) {
                element.find('.fc-list-item-title').append('<i class="fas fa-thumbtack"></i><br/>');
            } else {
                element.find('.fc-title').prepend("<br/>");
            }
            element.find('.fc-title').append("<br/>" + event.bookingCount + " / 3<br/>");
        },
        // DOM-Event handling for Calendar Eventblocks (why do js people suck at naming)
        eventOverlap: function(stillEvent, movingEvent) {
            return stillEvent.color !== movingEvent.color;
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
                if (!removeEvent(event, jsEvent, view)) {
                    return;
                }
                console.log("NOT IMPLEMENTED ON THE BACKEND YET");
                $('#calendar').fullCalendar('removeEvents', event.id);
        },
        businessHours: {
            // days of week. an array of zero-based day of week integers (0=Sunday)
            dow: [1, 2, 3, 4, 5], // Monday - Thursday
            start: '8:00', // a start time (10am in this example)
            end: '18:00', // an end time (6pm in this example)
        },
        // Controls view of agendaWeek
        minTime: '07:00:00',
        maxTime: '19:00:00',
        allDaySlot: false,       // shows slot @ top for allday events
        slotDuration: '00:30:00' // hourly divisions
    });
    $('.fc-today').css("background-color", "#FEFEFE");
    loadAddEvent();
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
        || document.querySelector("#modifier").value == "") {
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
    let event = {}
    event.start = moment(document.querySelector("#start").value).format();
    event.end = moment(document.querySelector("#end").value).format();
    event.room = document.querySelector("#room").value;
    event.color = event.room
    event.modifier = parseInt(document.querySelector("#modifier").value);
    event.note = document.querySelector("#note").value;
    eventJson = JSON.stringify(event);
    // Make ajax POST request with booking request or request bookign delete if already booked
    $.ajax({
        url: '/api/v1/events/add',
        type: 'POST',
        contentType: 'json',
        data: eventJson,
        dataType: 'json',
        success: function (data) {
            event.id = data.id;
            event.color = data.color;
            event.title = "<br>Facilitation 0/3";
            $('#calendar').fullCalendar('renderEvent', event); // render event on calendar
        },
        error: function (xhr, ajaxOptions, thrownError) {
            alert("Request failed: " + xhr.responseText);
        }
    });
}