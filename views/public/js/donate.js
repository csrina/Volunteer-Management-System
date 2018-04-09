document.addEventListener("DOMContentLoaded", () => {
    // Make request for the donation page data
    $.ajax({
        url: "/api/v1/donate",
        type: 'GET',
        contentType:'json',
        success: (data => { createWidgets(data); }),
        error: (xhr => makeToast("error", "Failed to retrieve page data: " + xhr.responseText)),
        dataType: 'json'
    });
    document.querySelector('select[name="DoneeSelect"]').onchange=doneeChangeHandler;
    document.getElementById('AmountField').onchange=amtFieldChangeHandler;
});

function amtFieldChangeHandler(event) {
    $(this).attr("value", event.target.value);
    console.log($(this).attr("value") + "   " + event.target.value);
}

function doneeChangeHandler(event) {
    event.target.setAttribute("value", event.target.value);
}

function sendDonation() {
    if (!validateForm()) {
        makeToast("warning", "Donation amount cannot exceed hours available!");
        return
    }

    let jsonData = JSON.stringify({
        amount: parseFloat(document.getElementById("AmountField").getAttribute("value")),
        donee: parseInt(document.getElementById("DoneeSelect").getAttribute("value")),
        donor: 0,
        id: 0,
        date: moment(),
    });

    console.log(jsonData);
    $.ajax({
        url: "/api/v1/donate",
        type: 'POST',
        contentType:'json',
        data: jsonData,
        success: function(data) {
            console.log(data);
            makeToast("success", "Donation sent! Donation ID: " + data.id);
            refreshWidgets(data);
        },
        error: (xhr => makeToast("error", "Failed to send donation: " + xhr.responseText)),
        dataType: 'json'
    });
}

function createWidgets(data) {
    data.families.forEach(fam => {
        $("#DoneeSelect").append("<option value='" + fam.familyID + "' >" + fam.familyName + "</option>");
    });
    document.getElementById("hoursAvail-text").innerHTML = (data.hoursAvail > 0) ? data.hoursAvail : 0;
}

function refreshWidgets(data) {
    let e = document.getElementById("hoursAvail-text")
    e.innerHTML = parseFloat(e.innerHTML) - parseFloat(data.amount);
}

// Changes appearance and returns a bool (isValid? true if valid)
function validateForm() {
    return true;
}
