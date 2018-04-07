// This function will send a booking request to the server
// Like the updateEvent function, post-demo I will refactor this out
// and have the templates populate based on role. Additional auth checks
// server side to ensure correct user/role and such should still take place
function requestBooking() {
    let event = JSON.parse($('.modal-footer').find('#modalEventData').text());
    if (event.booked && event.bookingCount >= 3 && !prompt("This block is pretty crowded, are you sure you want to proceed?")) {
        return;
    }
    // Block info for booking
    let booking_json = JSON.stringify({
        id:         event.id,
        start:      event.start,
        end:        event.end,
    });

    // Make ajax POST request with booking request or request bookign delete if already booked
    $.ajax({
        url: '/api/v1/events/book',
        type: 'POST',
        contentType:'json',
        data: booking_json,
        dataType:'json',
        success: (data => { bookingRequestSuccess(data, event.id); }),
        error: (xhr => { makeToast("error", "Booking request failed: " + xhr.responseText); })
    });
}

function bookingRequestSuccess(data, id) {  // We expect the server to return json data with a msg field
    // noinspection Annotator
    makeToast("success", data.msg);
    let event = $('#calendar').fullCalendar('clientEvents', id)[0];
    console.log(event);
    $('#modalConfirm').removeClass((event.booked) ? "btn-warning" : "btn-success");
    event.booked = !event.booked;
    if (event.booked === true) {
        if (!event.bookings) {
            event.bookings = [];
        }
        // noinspection Annotator
        event.bookings.push({userName: data.userName, userId: data.userId});
        event.bookingCount++;
    } else {
        // noinspection Annotator
        event.bookingCount--;
        let len = event.bookings.length;
        for (let i = 0; i < len; i++) {
            if (event.bookings[i].userId == data.userId) {
                event.bookings.splice(i, 1);
                break;
            }
        }
    }
    $('#calendar').fullCalendar('updateEvent', event);
    $('#eventDetailsModal').modal('hide');
}

$(document).ready(function() {
    // page is now ready, initialize the calendar...
    $('#calendar').fullCalendar({
        weekends: false,
        header: {
            left: 'today',
            center: 'prev, title, next',
            right: 'agendaWeek, month'
        },
        contentHeight: 'auto',
        defaultView: "agendaWeek",
        events: "/api/v1/events/scheduler",    // link to events (bookings + blocks feed)
        allDayDefault: false,        // blocks are not all-day unless specified
        themeSystem: "bootstrap4",
        editable: false,                 // Need to use templating engine to change bool based on user's rolego ,
        eventRender: function(event, element, view) {
            let fctime = element.find('.fc-time');
            let fctitle = element.find('.fc-title');
            fctime.css("font-size", "1em");
            fctitle.prepend("<br/>");
            fctitle.css("font-size", "1.0em");

            if (event.booked) {
                fctime.append('<br><i class="fas fa-thumbtack"></i><br>');
            } else {
                fctime.append('<br>' + event.bookingCount + "/" +
                    ((!!event.capacity) ? event.capacity.toString() : "3"));
            }

            return renderFiltered(event);
        },
        loading: ((isLoading, view) => {
            if (isLoading) { clearFilterButtons(); }
        }),
        eventOverlap: function(stillEvent, movingEvent) {
            return stillEvent.color === movingEvent.color;
        },
        eventClick: function(event, jsEvent, view) {
            $('#eventModalTitle').html("Book " + event.title);
            $('#modalEventRoom').html(event.room + " Room").css("color", event.color);
            $('#modalEventTime').html(event.start.format("ddd, hA") + " - " + event.end.format("hA"))
            $('#modalEventValue').html(moment.duration(event.end.diff(event.start)).asHours() * event.modifier);
            $('#eventNote').html(event.note);
            $('#modalEventData').html(JSON.stringify({
                id: event.id,
                start: event.start,
                end:   event.end,
                booked: event.booked,
                room: event.room,
                bookingCount: event.bookingCount,
            }));
            let len = (!event.bookings) ? 0 : event.bookings.length;
            let bookingsHTML = "";
            if (len === 0) {
                bookingsHTML = "No bookings yet <br> You could be the first!"
            }
            for (let i = 0; i < len; i++) {
                bookingsHTML += "<p>" + event.bookings[i].userName + "</p><br>";
            }
            $('#eventBookings').html(bookingsHTML);
            // set text and colour of submit button:
            // if booking and potential overbooking --> danger,
            // else if booking open slot --> success
            // else if unbook --> warning
            $('#modalConfirm').html((!event.booked) ? "Book" : "Unbook")
                              .addClass((!event.booked) ? ((event.bookingCount < 3) ? "btn-success" : "btn-danger") : "btn-warning");
            $('#eventDetailsModal').modal('show');
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
});