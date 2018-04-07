// Callback function for drag/drops and resizes of existing events
// Note: We dont want this to be populated if we aren't admin.
// post-demo will refactor this out into templates populated differently based on the role of the user

// For teachers, this is classes other than their own
function showStaticModal(btn) {
    let event = $('#calendar').fullCalendar('clientEvents', btn.getAttribute("data-id"))[0]; // get event from returned array
    $('#eventModalTitle').html("Book " + event.title);
    $('#modalEventRoom').html(event.room + " Room").css("color", event.color);
    $('#modalEventTime').html(event.start.format("ddd, hA") + " - " + event.end.format("hA"))
    $('#modalEventValue').html(moment.duration(event.end.diff(event.start)).asHours() * event.modifier);
    $('#eventNote').html(event.note);
    let len = (!event.bookings) ? 0 : event.bookings.length;
    let bookingsHTML = "";
    if (len === 0) {
        bookingsHTML = "No bookings yet :(";
    }
    for (let i = 0; i < len; i++) {
        bookingsHTML += "<p>" + event.bookings[i].userName + "</p><br>";
    }
    $('#eventBookings').html(bookingsHTML);
    $('#eventDetailsModal').modal('show');
}

// This is the modal with edit power, shown to teachers for the classroom they teach in
function showEditModal(btn) {
    let event = $('#calendar').fullCalendar('clientEvents', btn.getAttribute("data-id"))[0]; // get event from returned array
    // Make open/closing button tags which allows us to insert  the current data as the btton text
    let openEditNoteButton = "<button type='button' class='btn btn-outline-secondary border-0 mpb-1' "
        + "data-fieldName='note' onclick='editEventDetails(this)' data-id='"
        + event.id + "'>";
    // Need another prefix tag for title
    let openEditTitleButton = "<button type='button' class='btn btn-outline-secondary border-0 mpb-1' "
        + "data-fieldName='title' onclick='editEventDetails(this)' data-id='"
        + event.id + "'>";

    let closeEditButton = "    <span class='far fa-edit fa-lg'></span></button>"; // close the edit button
    // use the open/close button strings to create edit buttons containing the data to be altered
    $('#eventModalTitle').html(openEditTitleButton + "<h5>" + event.title + closeEditButton + "</h5>");
    $('#modalEventRoom').html(event.room + " Room").css("color", event.color);
    $('#modalEventTime').html(event.start.format("ddd, hhA") + " - " + event.end.format("hA"))

    let hourlyValue = moment.duration(event.end.diff(event.start)).asHours() * event.modifier;
    $('#modalValueLabel').append("<h5 id='modalEventValue' class='text-primary'>" + hourlyValue + "</h5>");
    $('#eventNote').html(openEditNoteButton + "<p class='text-muted'>" + event.note + closeEditButton + "</p>");

    let len = (!event.bookings) ? 0 : event.bookings.length;
    let bookingsHTML = "";
    if (len === 0) {
        bookingsHTML = "No bookings yet <br> You could be the first!";
    }

    // make a button for each user (so they can be unbooked easily and clearly... unlike this code hehe)
    for (let i = 0; i < len; i++) {
        bookingsHTML += "<p>" + event.bookings[i].userName + "</p><br>";
    }
    $('#eventBookings').html(bookingsHTML); // set event bookings with the html built in the loop

    $('#eventDetailsModal').modal('show'); // spawn our modal
}

function storeChangesToEvent(event, delta, revertFunc, jsEvent, ui, view) {
    // Extract block data required for updating on server
    let temp = {
        id: event.id,
        start: event.start.format(),
        end:   event.end.format(),
        title: event.title,
        note:  event.note,
        modifier: event.modifier,
    };

    if (!temp.start.endsWith("Z")) { temp.start = temp.start + "Z"; }
    if (!temp.end.endsWith("Z")) { temp.end = temp.end + "Z"; }

    let event_json = JSON.stringify(temp);
    // Make ajax post request with updated event data
    $.ajax({
        url: '/api/v1/events/update',
        type: 'POST',
        contentType:'json',
        data: event_json,
        dataType:'json',
        success: function(data) {
            makeToast("success", data.msg);
        },
        error: function(xhr, ajaxOptions, thrownError) {
            if (!!revertFunc) {
                revertFunc();
            }
            makeToast("error", "Request failed: " + xhr.responseText);
        }
    });
}




function updateEventRefreshModal(event, btn) {
    $('#calendar').fullCalendar('updateEvent', event);
    $('#eventDetailsModal').one('hidden.bs.modal', function(e) {
        showEditModal(btn);
    }) .modal('hide');
}

function editEventDetails(btn) {
    let event = $('#calendar').fullCalendar('clientEvents', btn.getAttribute("data-id"))[0];
    let field = btn.getAttribute("data-fieldName");

    /*
     * Gross nested if-elses... just a bunch of error checking with slightly different edge case checking
     * Basically we can have empty notes, but we warn; cannot have empty title so throw error
     * If numeric notes/title, we make change but show warning.
     * If cancel is clicked (only matter for notes because can be empty) --> we leave silently [for title we throw an error toast b/c can't have empty title so cancel is covered]
     * Lastly, we client-side validate the modifier value
     */
    if (field === "title") {
        temp = prompt("Enter the new title: ");
        if (temp === null || temp === "" || temp === undefined) {
            makeToast("error", "Cannot have empty title");
            return;
        } else if (!isNaN(temp)) {
            makeToast("warning", "The new title is a number, is that a typo?");
        }
        event.title = temp;
    } else if (field === "note") {
        temp = prompt("Enter the new description: ");
        if (temp === null) {
            return;
        } else if (temp === "") {
            makeToast("warning", "You left the description field empty, was that on purpose?")
        } else if (!isNaN(event.note)) {
            makeToast("warning", "The new description is numeric, was that on purpose?");
        }
        event.note = temp;
    } else if (field == "modifier") {
        makeToast("warning", "Only admins may modify this, sorry.");
        return;
    } else {
        makeToast("error", "An unpredicted error occurred");
        return
    }
    storeChangesToEvent(event);
    updateEventRefreshModal(event, btn)
}

// REmoves an event from the calendar and the associated TB/Bookings from the database
function removeEvent(btn) {
    let event = $('#calendar').fullCalendar('clientEvents', btn.getAttribute("data-id"))[0];
    yn = confirm("Are you sure you want to delete this event?")
    if (!yn) {
        return false; // event should not be deleted
    }
    let event_json = JSON.stringify({
        id:    event.id,
    });
    // Make ajax post request with updated event data
    $.ajax({
        url: '/api/v1/events/delete',
        type: 'POST',
        contentType:'json',
        data: event_json,
        dataType:'json',
        success: function(data) {
            makeToast("success", data.msg);
        },
        error: function(xhr, ajaxOptions, thrownError) {
            makeToast("error", "Request failed: " + xhr.responseText);
        }
    });
    // remove event from calendar
    $('#calendar').fullCalendar('removeEvents', event.id);
}

$(document).ready(function() {
    // ensure created button is deleted on modal close
    $('#eventDetailsModal').on('hide.bs.modal', function (e) {
        $('#modalNoteLabel').html("Description:");
        $('#modalBookedLabel').html("Attending:");
        $('#modalValueLabel').html("Value:");
    });

    // page is now ready, initialize the calendar...
    $('#calendar').fullCalendar({
        weekends: false,
        header: {
            left: 'today',
            center: 'prev, title, next',
            right: 'agendaWeek, month'
        },
        agendaEventMinHeight: 75,
        defaultView: "agendaWeek",
        contentHeight: 'auto',
        events: "/api/v1/events/scheduler",    // link to events (bookings + blocks feed)
        allDayDefault: false,        // blocks are not all-day unless specified
        themeSystem: "bootstrap4",
        editable: false,                 // Need to use templating engine to change bool based on user's rolego ,
        eventRender: function(event, element, view) {
            element.find('.fc-time').css("font-size", "1em")
                .append('    ' + event.bookingCount + "/" +
                            ((!!event.capacity) ? event.capacity.toString() : "/3"));
            let fcTitle = element.find('.fc-title')
                                 .css("font-size", "1.2em")
                                 .append("<br>"); // gets the fcTitle jQuery elem
            /* give edit event button to the events in a teacher's classroom */
            let teachesRoom = event.booked; // WE USE THIS FOR TEACHERS TO INDICATE IF CLASSROOM IS ONE THEY TEACH IN ROOM
            if (teachesRoom) {
                fcTitle.append("<button type='button' class='btn btn-outline-primary border-0 btn-sm' data-id='" +
                    event.id + "' onclick='showEditModal(this)'><i class='far fa-edit fa-lg'></i></button>");
            } else {
                fcTitle.append("<button type='button' class='btn btn-outline-primary border-0 btn-sm' data-id='" +
                    event.id + "' onclick='showStaticModal(this)'><i class=\"fas fa-external-link-alt fa-lg\"></i></button>");
            }
            return renderFiltered(event);
        },
        // DOM-Event handling for Calendar Eventblocks (why do js people suck at naming)
        eventOverlap: function(stillEvent, movingEvent) {
            if (stillEvent.color === movingEvent.color) {
                makeToast("warning", "Events of same color may not overlap");
            }
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
});
