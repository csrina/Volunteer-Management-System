// JAVASCRIPT FOR ALL CALENDARS REGARDLESS OF ROLE -- MORE TO COME ONCE MODAL IS STANDARDIZED TO DATA-ID (SEE ADMIN)

// Determines if an event satisfies the filters applied
function isVisibleEvent(event) {
    return $("#" + event.room + "btn").attr("data-value") == "on";
}

function clearFilterButtons() {
    $("#filterButtons").html();
}

// Add a filter with the given text value
function addFilterButton(buttonText) {
    if ($("#" + buttonText + "btn").length != 0) {
        return; // Already have dis
    }
    let btn = "<button type='button' class='btn btn-sm btn-primary active' aria-pressed='true' " +
        "data-value='on' id='" + buttonText + "btn' onclick='changeFilter(this)' style='font-size: 0.85rem;'>";
    btn +=  buttonText + "</button>";
    $('#filterButtons').append(btn);
}

// Update filters & refresh displayed events
function changeFilter(btn) {
    if (btn.getAttribute("data-value") == "on") { // button WAS active (not filtered) --> set to inactive (filtered)
        btn.setAttribute("data-value", "off"); // change data attr
        $(btn).removeClass("btn-primary").addClass("btn-light"); // set light to clearly indicate not shown
    } else {
        btn.setAttribute("data-value", "on"); // change data attr to on
        $(btn).removeClass("btn-light").addClass("btn-primary"); // set primary to indicate more clearly
    }

    let events = $("#calendar").fullCalendar("clientEvents");
    $("#calendar").fullCalendar("updateEvents", events);
}

// called during eventRender to apply filtering (+dynamically add filter buttons by room)
function renderFiltered(event) {
    addFilterButton(event.room);
    return isVisibleEvent(event);
}