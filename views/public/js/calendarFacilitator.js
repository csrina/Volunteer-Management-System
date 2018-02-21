// This function will send a booking request to the server
// Like the updateEvent function, post-demo I will refactor this out
// and have the templates populate based on role. Additional auth checks
// server side to ensure correct user/role and such should still take place
function requestBooking(event, jsEvent, view) {
    var promptStr = "Confirm booking ";
    console.log(event.booked);
    if (event.booked == true) {
        promptStr += "Cancellation (";
    } else {
        promptStr += "Booking (";
    }
    if (!confirm(promptStr + event.start.toString() + ", in the " + event.room + " room)")) {
        return;
    }
    // Block info for booking
    booking_json = JSON.stringify({
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
            if (event.booked == true) {
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
    // page is now ready, initialize the calendar...
    $('#calendar').fullCalendar({
        // Education use (both now and if deployed!)
        weekends: false,
        header: {
            left: 'today',
            center: 'prev, title, next',
            right: 'agendaWeek, month'
        },
        agendaEventMinHeight: 5,
        defaultView: "agendaWeek",
        events: "/api/v1/events",    // link to events (bookings + blocks feed)
        allDayDefault: false,        // blocks are not all-day unless specified
        themeSystem: "bootstrap3",
        editable: false,                 // Need to use templating engine to change bool based on user's rolego ,
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
        eventOverlap: function(stillEvent, movingEvent) {      // Event blocks in different rooms may overlap, events in same room may not
            // Note: events may overlap on import; moving events will not be allowed to over lap
            //       That is, we must constrain the overlap when we make event creation possible
            if (stillEvent.color === movingEvent.color) {
                return false;
            }
            return true;
        },
        eventClick: function(event, jsEvent, view) {
                if (event.bookingCount > 3) {
                    alert("Sorry, only administrators can over-book time blocks.")
                    return;
                }
                requestBooking(event, jsEvent, view);
        },
        eventMouseover: function (event, jsEvent, view) {
            $(this).addClass("expand");
        },
        eventMouseout: function (event, jsEvent, view) {
            $(this).removeClass("expand");
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
