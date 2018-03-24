const colorRange = [
    [0.0, "#AA0000" ],
    [0.25, "#FF5500"],
    [0.5, "#FF9A00"],
    [0.75, "#FAD201"],
    [1.0, "#00AA00"]
];

var elems = {
    gDone: "",
    gBooked: "",
    chart: "",
};

function makeToast(type, msg) {
    toastr.options = {
        "closeButton": true,
        "debug": false,
        "newestOnTop": true,
        "progressBar": false,
        "positionClass": "toast-top-right",
        "preventDuplicates": true,
        "onclick": null,
        "showDuration": "300",
        "hideDuration": "1000",
        "timeOut": "5000",
        "extendedTimeOut": "1000",
        "showEasing": "swing",
        "hideEasing": "linear",
        "showMethod": "fadeIn",
        "hideMethod": "fadeOut"
    };
    Command: toastr[type](msg);
}

// This function will send a booking removal request to the server
function requestRemoval(event) {
    let promptStr = "Are you sure you want to remove yourself from:\n";
    if (!confirm(promptStr + event.start.toString() + ", in the " + event.room + " room")) {
        return;
    }
    // Block info for booking
    let booking_json = JSON.stringify({
        id:         event.id,
        start:      event.start.toString(),
        end:        event.end.toString(),
    });

    // Make ajax POST request with booking request or request bookign delete if already booked
    $.ajax({
        url: '/api/v1/events/book',
        type: 'POST',
        contentType:'json',
        data: booking_json,
        dataType:'json',
        success: function(data) {  // We expect the server to return json data with a msg field
            makeToast("success", data.msg);
            $("#calendar").fullCalendar("removeEvents", event.id);
            refreshWidgets();
        },
        error: function(xhr, ajaxOptions, thrownError) {
            makeToast("error", "Request failed: " + xhr.responseText);
        }
    });
}

/*
 * Spawns our chart (hours vs. time)
 */
function chartInit(elementId, data) {
    /* We get datasets in a map of {id => dataset}, break them into array */
    console.log(data.history)
    let datasets = [];
    for (var key in data.history) {
        datasets.push(data.history[key]);
    }

    let ctx = document.getElementById(elementId).getContext("2d");
    // noinspection ES6ConvertVarToLetConst
    var hoursChart = new Chart(ctx, {
        type: "line",
        data:
            {
                datasets: datasets,
            },
        options: {
            scales: {
                xAxes: [{
                    type: 'time',
                    distribution: 'linear',
                    time: {
                        stepSize: 1,
                        unit: 'week',
                        isoWeekday: true,
                        minUnit: "week",
                        parser: "YYYY-MM-DD",
                        max: data.endOfPeriod,
                        min: data.startOfPeriod,
                    },
                    source: "auto",
                }],
                yAxes: [{
                    min: 0,
                    max: 12.5,
                    stepSize: 1,
                }]
            }
        }
    });
    return hoursChart;
}

// Where elementID is the div to use and value is the number of hours in our case
// colorRange is an array format [[0.00, "color"], ...[1.0, "#color"]]
function gaugeInit(elementId, value, goal) {
    const opts = {
        angle: -0.35, // The span of the gauge arc
        lineWidth: 0.11, // The line thickness
        radiusScale: 1, // Relative radius
        pointer: {
            length: 0.0, // // Relative to gauge radius
            strokeWidth: 0.00, // The thickness
            color: '#000000' // Fill color
        },
        limitMax: true,     // If false, max value increases automatically if value > maxValue
        limitMin: true,     // If true, the min value of the gauge will be fixed
        percentColors: colorRange,
        strokeColor: '#E0E0E0',  // to see which ones work best for you
        generateGradient: true,
        highDpiSupport: true,     // High resolution support
    };

    // noinspection ES6ConvertVarToLetConst
    var element = document.getElementById(elementId);
    element.style.zIndex = 1;
    // noinspection ES6ConvertVarToLetConst
    var gauge = new Gauge(element).setOptions(opts);
    gauge.maxValue = goal;
    gauge.setMinValue(0);
    gauge.animationSpeed = 60;
    gauge.set(value);
    // setup the text
    element = document.getElementById(elementId + "-text");
    element.style.color = gauge.getColorForValue(value);
    if (value < goal) {
        if ((value / goal) < 0.5) {
            element.innerHTML = "<h3 class='text-danger'>" + value + "h</h3>";
        } else {
            element.innerHTML = "<h3 class='text-warning'>" + value + "h</h3>";
        }

    } else {
        element.innerHTML = "<h3 class='text-success'>" + value + "h</h3>";
    }
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
        contentHeight: 400,
        defaultView: "list",
        duration: {days: 14},        // two week intervals shown for upcoming events
        events: "/api/v1/events/dash",    // EventsFeed with dash as its target
        allDayDefault: false,        // blocks are not all-day unless specified
        themeSystem: "bootstrap4",
        editable: false,                 // Need to use templating engine to change bool based on user's rolego ,
        eventRender: function(event, element, view) {
            element.find('.fc-list-item-title').append("  " + event.bookingCount + "/3")
            element.find('.fc-list-item-title').append('<i class="fas fa-thumbtack"></i><br/>');

        },
        // DOM-Event handling for Calendar Eventblocks (why do js people suck at naming)
        eventOverlap: false,
        eventClick: function(event, jsEvent, view) {
            requestRemoval(event);
            refreshWidgets();
        },
        businessHours: {
            // days of week. an array of zero-based day of week integers (0=Sunday)
            dow: [1, 2, 3, 4, 5], // Monday - Thursday
            start: '8:00',
            end: '18:00'
        }
    });
});

function refreshGauge(gauge, parentID, newValue) {
    gauge.set(newValue);
    // setup the text
    element = document.getElementById(parentID + "-text");
    element.style.color = gauge.getColorForValue(newValue);
    element.innerHTML = "<h3>" + newValue + "h</h3>";
}


// Perform ajax get-request to fiddle with our non-calendar bits
function refreshWidgets() {
    var data = {};
    $.ajax({
        url: "/api/v1/dashboard",
        type: 'GET',
        contentType:'json',
        data: data,
        success: function(data) {
            refreshGauge(elems.gDone, "hoursDone", data.hoursDone);
            refreshGauge(elems.gBooked, "hoursBooked", data.hoursBooked);
            elems.chart.data.datasets[0] = data.history1;
            elems.chart.data.datasets[1] = data.history2;
            elems.chart.update();
        },
        error: function(xhr, ajaxOptions, thrownError) {
            makeToast("error", "Refresh failed: " + xhr.responseText);
        },
        dataType: 'json'
    });
}

/* Wait for DOM to load, then get the gauges */
document.addEventListener("DOMContentLoaded", () => {
    // get values from server
    // noinspection ES6ConvertVarToLetConst
    var data = {};

    $.ajax({
        url: "/api/v1/dashboard",
        type: 'GET',
        contentType:'json',
        data: data,
        success: function(data) {
            elems.gDone   = gaugeInit("hoursDone", data.hoursDone, data.hoursGoal);
            elems.gBooked = gaugeInit("hoursBooked", data.hoursBooked, data.hoursGoal);
            elems.chart   = chartInit("hoursChart", data);
        },
        dataType: 'json'
    });
});
