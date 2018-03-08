const colorRange = [
    [0.0, "#AA0000" ],
    [0.25, "#FF5500"],
    [0.5, "#FF9A00"],
    [0.75, "#FAD201"],
    [1.0, "#00AA00"]
];

function load() {
    let req = new XMLHttpRequest();
    req.addEventListener("load", function(evt) {
	    let data = JSON.parse(req.response);
	    console.log(req.response);
	    input(data);
    });
    req.open("GET", "http://localhost:8080/api/v1/dashboard");
    req.send();
}

function input(data) {
    let needed = 0;
    if (data.children === 1) {
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
}

// This function will send a booking removal request to the server
function requestRemoval(event) {
    let promptStr = "Are you sure you want to remove yourself from:\n";

    if (!confirm(promptStr + event.start.toString() + ", in the " + event.room + " room")) {
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
            event.bookingCount--;
            $('#calendar').fullCalendar('updateEvent', event);
        },
        error: function(xhr, ajaxOptions, thrownError) {
            alert("Request failed: " + thrownError);
        }
    });
}

/*
 * Spawns our chart (hours vs. time)
 */
function chartInit(elementId, data) {
    let ctx = document.getElementById(elementId).getContext("2d");
    var hoursChart = new Chart(ctx, {
        type: "line",
        data:
            {
                labels:
                    [
                        "January", "February", "March", "April", "May","June","July"
                    ],
            datasets:
                [
                    {
                        label:"Hours/Week",
                        data:[65,59,80,81,56,55,40],
                        fill:false,
                        borderColor:"rgb(75, 192, 192)",
                        lineTension:0.15
                    }
                ]
            },
        options:{
            spanGaps: true
        }
    });
}

// Where elementID is the div to use and value is the number of hours in our case
// colorRange is an array format [[0.00, "color"], ...[1.0, "#color"]]
function gaugeInit(elementId, value) {
    const opts = {
        angle: -0.35, // The span of the gauge arc
        lineWidth: 0.11, // The line thickness
        radiusScale: 1, // Relative radius
        pointer: {
            length: 0.0, // // Relative to gauge radius
            strokeWidth: 0.00, // The thickness
            color: '#000000' // Fill color
        },
        fontSize: 24,
        renderTicks: {
            divisions: 1,
            divWidth: 5,
            divLength: 1,
            divColor: "#000000",
            subDivisions: 0,
            subLength: 0,
            subWidth: 0,
        },
        limitMax: true,     // If false, max value increases automatically if value > maxValue
        limitMin: true,     // If true, the min value of the gauge will be fixed
        percentColors: colorRange,
        strokeColor: '#E0E0E0',  // to see which ones work best for you
        generateGradient: true,
        highDpiSupport: true,     // High resolution support
    };

    var element = document.getElementById(elementId)
    element.style.zIndex = 1;
    var gauge = new Gauge(element).setOptions(opts);
    gauge.maxValue = 5;
    gauge.setMinValue(0);
    gauge.animationSpeed = 60;
    gauge.set(value);
    // setup the text
    element = document.getElementById(elementId + "-text");
    element.style.color = gauge.getColorForValue(value);
    element.innerHTML = value;
    return gauge; // return for use
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
        theme: "bootstrap",
        aspectRatio: 0.33,
        defaultView: "list",
        duration: {days: 14},        // two week intervals shown for upcoming events
        events: "/api/v1/events/dash",    // EventsFeed with dash as its target
        allDayDefault: false,        // blocks are not all-day unless specified
        themeSystem: "bootstrap3",
        editable: false,                 // Need to use templating engine to change bool based on user's rolego ,
        eventRender: function(event, element, view) {
            element.find('.fc-list-item-title').append("  " + event.bookingCount + "/3");
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
    });
});

/* Wait for DOM to load, then get the gauges */
document.addEventListener("DOMContentLoaded", () => {
        // get values from server
        gaugeInit("hoursDone", 2); // replace literal value (2, 5) with value from server/calculated values
        gaugeInit("hoursBooked", 5);
        chartInit("hoursChart", "");
        load();
    });

